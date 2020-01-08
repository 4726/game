package main

import (
	"context"
	"sync"
	"time"

	"github.com/4726/game/services/matchmaking/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type QueueService struct {
	queues      map[pb.QueueType]*Queue
	matches     map[uint64]*Match
	queueTimes  map[pb.QueueType]*QueueTimes
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
	queues := map[pb.QueueType]*Queue{}
	queues[pb.QueueType_UNRANKED] = NewQueue(10000)
	queues[pb.QueueType_RANKED] = NewQueue(10000)
	queueTimes := map[pb.QueueType]*QueueTimes{}
	queueTimes[pb.QueueType_UNRANKED] = NewQueueTimes(1000)
	queueTimes[pb.QueueType_RANKED] = NewQueueTimes(1000)
	qs := &QueueService{
		queues,
		map[uint64]*Match{},
		queueTimes,
		opts,
		map[uint64]QueueChannels{},
		sync.Mutex{},
	}
	go qs.notifyQueueStateChanges()
	return qs
}

func (s *QueueService) Join(in *pb.JoinQueueRequest, outStream pb.Queue_JoinServer) error {
	queue := s.queues[in.GetQueueType()]
	found, matchID, users, err := queue.EnqueueAndFindMatch(in.GetUserId(), in.GetRating(), s.opts.RatingRange, s.opts.PlayerCount)
	if err != nil {
		if err == ErrAlreadyInQueue {
			return status.Error(codes.FailedPrecondition, err.Error())
		}
		return status.Error(codes.Internal, err.Error())
	}

	if !found {
		resp := &pb.JoinQueueResponse{
			UserId:          in.GetUserId(),
			QueueType:       in.GetQueueType(),
			MatchId:         uint64(0),
			Found:           false,
			SecondsToAccept: 20,
		}
		if err := outStream.Send(resp); err != nil {
			queue.DeleteOne(in.GetUserId())
			return err
		}
	} else {
		userIDs := []uint64{}
		for _, v := range users {
			userIDs = append(userIDs, v.UserID)
		}
		s.matches[matchID] = NewMatch(userIDs, time.Second*20)
	}

	leaveQueueCh := make(chan struct{}, 1)
	matchFoundCh := make(chan uint64, 1)
	matchStartFailCh := make(chan struct{}, 1)
	s.inQueueLock.Lock()
	s.inQueue[in.GetUserId()] = QueueChannels{leaveQueueCh, matchFoundCh, matchStartFailCh}
	s.inQueueLock.Unlock()

	for {
		select {
		case <-leaveQueueCh:
			return nil
		case matchID := <-matchFoundCh:
			resp := &pb.JoinQueueResponse{
				UserId:          in.GetUserId(),
				QueueType:       in.GetQueueType(),
				MatchId:         matchID,
				Found:           true,
				SecondsToAccept: 20,
			}
			if err := outStream.Send(resp); err != nil {
				queue.DeleteOne(in.GetUserId())
				return err
			}
		case <-matchStartFailCh:
			resp := &pb.JoinQueueResponse{
				UserId:          in.GetUserId(),
				QueueType:       in.GetQueueType(),
				MatchId:         uint64(0),
				Found:           false,
				SecondsToAccept: 20,
			}
			if err := outStream.Send(resp); err != nil {
				queue.DeleteOne(in.GetUserId())
				return err
			}
		}
	}
}

func (s *QueueService) Leave(ctx context.Context, in *pb.LeaveQueueRequest) (*pb.LeaveQueueResponse, error) {
	queue := s.queues[in.GetQueueType()]

	queue.DeleteOne(in.GetUserId())

	return &pb.LeaveQueueResponse{
		UserId:    in.GetUserId(),
		QueueType: in.GetQueueType(),
	}, nil
}

func (s *QueueService) Accept(in *pb.AcceptQueueRequest, outStream pb.Queue_AcceptServer) error {
	queue := s.queues[in.GetQueueType()]

	match, ok := s.matches[in.GetMatchId()]
	if !ok {
		return status.Error(codes.FailedPrecondition, "invalid match id")
	}
	ch := make(chan MatchStatus, 1)
	if err := match.Accept(in.GetUserId(), ch); err != nil {
		if err == ErrUserNotInMatch {
			return status.Error(codes.FailedPrecondition, err.Error())
		}
		return status.Error(codes.Internal, err.Error())
	}

	for {
		status := <-ch
		resp := &pb.AcceptQueueResponse{
			TotalAccepted: uint32(status.TotalAccepted),
			TotalNeeded:   uint32(status.TotalNeeded),
			QueueType:     in.GetQueueType(),
			Cancelled:     status.Cancelled,
			UserIds:       status.Players,
		}
		if err := outStream.Send(resp); err != nil {
			return err
		}
		if resp.Cancelled {
			//someone declined, adds user back into queue
			queue.MarkMatchFound(in.GetUserId(), false)
			return nil
		}
		if resp.GetTotalAccepted() == resp.GetTotalNeeded() {
			//everyone accepted, remove user from queue
			//also removes match
			queue.DeleteOne(in.GetUserId())
			delete(s.matches, in.GetMatchId())
			return nil
		}
	}
}

func (s *QueueService) Decline(ctx context.Context, in *pb.DeclineQueueRequest) (*pb.DeclineQueueResponse, error) {
	queue := s.queues[in.GetQueueType()]

	match, ok := s.matches[in.GetMatchId()]
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "invalid match id")
	}
	if err := match.Decline(in.GetUserId()); err != nil {
		if err == ErrUserNotInMatch {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	queue.DeleteOne(in.GetUserId())

	return &pb.DeclineQueueResponse{
		UserId:    in.GetUserId(),
		QueueType: in.GetQueueType(),
	}, nil
}

func (s *QueueService) Info(ctx context.Context, in *pb.QueueInfoRequest) (*pb.QueueInfoResponse, error) {
	queue := s.queues[in.GetQueueType()]
	queueTimes := s.queueTimes[in.GetQueueType()]

	return &pb.QueueInfoResponse{
		SecondsEstimated: uint32(queueTimes.EstimatedWaitTime(in.GetRating(), 100).Seconds()),
		UserCount:        uint32(queue.Len()),
	}, nil
}

func (s *QueueService) notifyQueueStateChanges() {
	ch := make(chan PubSubMessage, 1)
	s.queues[0].Subscribe(ch)
	for {
		msg := <-ch
		userID := msg.Data.UserID
		s.inQueueLock.Lock()
		qChs, ok := s.inQueue[userID]
		if !ok {
			s.inQueueLock.Unlock()
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
		s.inQueueLock.Unlock()
	}
}
