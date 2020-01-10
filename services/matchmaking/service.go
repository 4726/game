package main

import (
	"context"
	"sync"
	"time"

	"github.com/4726/game/services/matchmaking/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var defaultMatchAcceptTimeout = time.Second * 20
var defaultCheckTimeout = time.Minute

type QueueService struct {
	queue       *Queue
	matches     map[uint64]*Match
	matchesLock sync.Mutex
	queueTimes  *QueueTimes
	opts        QueueServiceOptions
	inQueue     map[uint64]QueueChannels
	inQueueLock sync.Mutex
}

type QueueChannels struct {
	LeaveCh        chan struct{} //no longer in queue
	MatchFoundCh   chan uint64   //match found waiting for accept
	MatchStartFail chan struct{} //not all users accepted the match
}

type QueueServiceOptions struct {
	//rating range where players can get matched with each other
	//ex) RatingRange of 100 allows a player with 1000 rating to match with a player with 1100 rating
	RatingRange uint64
	//number of players in a single match
	PlayerCount int
}

func NewQueueService(opts QueueServiceOptions) *QueueService {
	qs := &QueueService{
		queue:      NewQueue(10000),
		matches:    map[uint64]*Match{},
		queueTimes: NewQueueTimes(1000),
		opts:       opts,
		inQueue:    map[uint64]QueueChannels{},
	}

	go func() {
		for _ = range time.Tick(defaultCheckTimeout) {
			qs.matchesLock.Lock()
			for k, v := range qs.matches {
				if v.TimeSince() < defaultMatchAcceptTimeout {
					state := v.State()
					for _, userID := range state.Unknown {
						//remove all players that didnt reply from the queue
						qs.queue.DeleteOne(userID)
					}
					delete(qs.matches, k)
				}
			}
			qs.matchesLock.Unlock()
		}
	}()

	go qs.notifyQueueStateChanges()
	return qs
}

func (s *QueueService) Join(in *pb.JoinQueueRequest, outStream pb.Queue_JoinServer) error {
	s.inQueueLock.Lock()
	_, ok := s.inQueue[in.GetUserId()]
	if ok {
		s.inQueueLock.Unlock()
		return status.Error(codes.FailedPrecondition, ErrAlreadyInQueue.Error())
	}
	leaveQueueCh := make(chan struct{}, 1)
	matchFoundCh := make(chan uint64, 1)
	matchStartFailCh := make(chan struct{}, 1)
	s.inQueue[in.GetUserId()] = QueueChannels{leaveQueueCh, matchFoundCh, matchStartFailCh}
	s.inQueueLock.Unlock()

	defer func() {
		s.inQueueLock.Lock()
		delete(s.inQueue, in.GetUserId())
		s.inQueueLock.Unlock()
	}()

	found, matchID, users, err := s.queue.EnqueueAndFindMatch(in.GetUserId(), in.GetRating(), s.opts.RatingRange, s.opts.PlayerCount)
	if err != nil {
		if err == ErrAlreadyInQueue {
			return status.Error(codes.FailedPrecondition, err.Error())
		}
		return status.Error(codes.Internal, err.Error())
	}

	resp := &pb.JoinQueueResponse{
		UserId:          in.GetUserId(),
		MatchId:         uint64(0),
		Found:           false,
		SecondsToAccept: 20,
	}
	if err := outStream.Send(resp); err != nil {
		s.queue.DeleteOne(in.GetUserId())
		return err
	}

	if found {
		userIDs := []uint64{}
		for _, v := range users {
			userIDs = append(userIDs, v.UserID)
		}
		s.matches[matchID] = NewMatch(userIDs)
	}

	for {
		select {
		case <-leaveQueueCh:
			return nil
		case matchID := <-matchFoundCh:
			resp := &pb.JoinQueueResponse{
				UserId:          in.GetUserId(),
				MatchId:         matchID,
				Found:           true,
				SecondsToAccept: 20,
			}
			if err := outStream.Send(resp); err != nil {
				s.queue.DeleteOne(in.GetUserId())
				return err
			}
		case <-matchStartFailCh:
			resp := &pb.JoinQueueResponse{
				UserId:          in.GetUserId(),
				MatchId:         uint64(0),
				Found:           false,
				SecondsToAccept: 20,
			}
			if err := outStream.Send(resp); err != nil {
				s.queue.DeleteOne(in.GetUserId())
				return err
			}
		}
	}
}

func (s *QueueService) Leave(ctx context.Context, in *pb.LeaveQueueRequest) (*pb.LeaveQueueResponse, error) {
	s.queue.DeleteOne(in.GetUserId())

	return &pb.LeaveQueueResponse{
		UserId: in.GetUserId(),
	}, nil
}

func (s *QueueService) Accept(in *pb.AcceptQueueRequest, outStream pb.Queue_AcceptServer) error {
	s.matchesLock.Lock()
	match, ok := s.matches[in.GetMatchId()]
	s.matchesLock.Unlock()

	readdToQueue := true
	defer func() {
		if readdToQueue {
			s.queue.MarkMatchFound(in.GetUserId(), false)
		}
	}()

	if !ok {
		//happens when someone denied or match times out
		return status.Error(codes.FailedPrecondition, ErrUserNotInMatch.Error())
	}
	ch := make(chan MatchStatus, 1)
	if err := match.Accept(in.GetUserId(), ch); err != nil {
		if err == ErrUserNotInMatch {
			//client error
			return status.Error(codes.FailedPrecondition, err.Error())
		}
		if err == ErrMatchCancelled {
			//happens when nobody gave match response and  match times out
			s.matchesLock.Lock()
			delete(s.matches, in.GetMatchId())
			s.matchesLock.Unlock()
			readdToQueue = false
			return status.Error(codes.FailedPrecondition, err.Error())
		}
		return status.Error(codes.Internal, err.Error())
	}

	for {
		status, ok := <-ch
		if !ok {
			//should not go here, keep just in case
			return nil
		}

		resp := &pb.AcceptQueueResponse{
			TotalAccepted: uint32(status.TotalAccepted),
			TotalNeeded:   uint32(status.TotalNeeded),
			Cancelled:     status.Cancelled,
			UserIds:       status.Players,
		}
		if err := outStream.Send(resp); err != nil {
			return err
		}
		if resp.Cancelled {
			//someone declined or timeout
			s.queue.MarkMatchFound(in.GetUserId(), false)
			s.matchesLock.Lock()
			delete(s.matches, in.GetMatchId())
			s.matchesLock.Unlock()
			return nil
		}
		if resp.GetTotalAccepted() == resp.GetTotalNeeded() {
			//everyone accepted, remove user from queue
			//also removes match
			s.queue.DeleteOne(in.GetUserId())
			delete(s.matches, in.GetMatchId())
			readdToQueue = false
			return nil
		}
	}
}

func (s *QueueService) Decline(ctx context.Context, in *pb.DeclineQueueRequest) (*pb.DeclineQueueResponse, error) {
	s.matchesLock.Lock()
	match, ok := s.matches[in.GetMatchId()]
	s.matchesLock.Unlock()
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, ErrUserNotInMatch.Error())
	}
	if err := match.Decline(in.GetUserId()); err != nil {
		if err == ErrUserNotInMatch {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	s.queue.DeleteOne(in.GetUserId())
	s.matchesLock.Lock()
	delete(s.matches, in.GetMatchId())
	s.matchesLock.Unlock()

	return &pb.DeclineQueueResponse{
		UserId: in.GetUserId(),
	}, nil
}

func (s *QueueService) Info(ctx context.Context, in *pb.QueueInfoRequest) (*pb.QueueInfoResponse, error) {
	return &pb.QueueInfoResponse{
		SecondsEstimated: uint32(s.queueTimes.EstimatedWaitTime(in.GetRating(), 100).Seconds()),
		UserCount:        uint32(s.queue.Len()),
	}, nil
}

func (s *QueueService) notifyQueueStateChanges() {
	ch := make(chan PubSubMessage, 1)
	s.queue.Subscribe(ch)

	for msg := range ch {
		userID := msg.Data.UserID

		s.inQueueLock.Lock()
		qChs, ok := s.inQueue[userID]
		s.inQueueLock.Unlock()
		if !ok {
			continue
		}

		switch msg.Topic {
		case PubSubTopicDelete:
			qChs.LeaveCh <- struct{}{}
		case PubSubTopicMatchFound:
			qChs.MatchFoundCh <- msg.Data.MatchID
		case PubSubTopicMatchNotFound:
			qChs.MatchStartFail <- struct{}{}
		}
	}
}
