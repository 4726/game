package main

import (
	"context"
	"time"
	"sync"

	"github.com/4726/game/services/matchmaking/pb"
)

type QueueService struct {
	queues     map[pb.QueueType]*Queue
	matches    map[uint64]Match
	queueTimes *QueueTimes
	opts       QueueServiceOptions
	inQueue map[uint64]QueueChannels
	inQueueLock sync.Mutex
}

type QueueChannels struct {
	LeaveCh chan struct{} //no longer in queue
	MatchFoundCh chan uint64 //match found waiting for accept
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
	queues[pb.QueueType_UNRANKED] = NewQueue()
	queues[pb.QueueType_RANKED] = NewQueue()
	qs := &QueueService{
		queues,
		map[uint64]Match{},
		NewQueueTimes(1000),
		opts,
		map[uint64]QueueChannels{},
		sync.Mutex{},
	}
	qs.notifyQueueStateChanges()
	return qs
}

func (s *QueueService) Join(in *pb.JoinQueueRequest, outStream pb.Queue_JoinServer) error {
	queue := s.queues[in.GetQueueType()]

	found, matchID, users, err := queue.EnqueueAndFindMatch(in.GetUserId(), in.GetRating(), s.opts.RatingRange, s.opts.PlayerCount)
	if err != nil {
		return err
	}

	if !found {
		resp := &JoinQueueResponse {
			UserID: in.GetUserId(),
			QueueType: in.GetQueueType(),
			MatchId: uint64(0),
			Found: false,
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
		s.matches[matchID] = NewMatch(userIDs, time.Second * 20)
	}

	leaveQueueCh := make(chan struct{}, 1)
	matchFoundCh := make(chan uint64, 1)
	matchStartFailCh := make(chan struct{}, 1)
	s.inQueueLock.Lock()
	s.inQueue[in.GetQueueType()] = QueueChannels{leaveQueueCh, matchFoundCh, matchStartFailCh}
	s.inQueueLock.Unlock()
	
	for {
		select {
		case <- leaveQueueCh:
			return nil
		case matchID := <- matchFoundCh:
			resp := &JoinQueueResponse {
				UserID: in.GetUserId(),
				QueueType: in.GetQueueType(),
				MatchId: matchID,
				Found: true,
				SecondsToAccept: 20,
			}
			if err := outStream.Send(resp); err != nil {
				queue.DeleteOne(in.GetUserId())
				return err
			}
		case <- matchStartFailCh:
			resp := &JoinQueueResponse {
				UserID: in.GetUserId(),
				QueueType: in.GetQueueType(),
				MatchId: uint64(0),
				Found: false,
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

	return &pb.LeaveQueueResponse{in.GetUserId(), in.GetQueueType()}, nil
}

func (s *QueueService) Accept(in *pb.AcceptQueueRequest, outStream pb.Queue_AcceptServer) error {
	queue := s.queues[in.GetQueueType()]

	match := s.matches[in.GetMatchId()]
	ch := make(chan MatchStatus, 1)
	if err := match.Accept(in.GetUserId(), ch); err != nil {
		return err
	}

	for {
		status <- ch
		resp := &AcceptQueueResponse{
			TotalAccepted: status.TotalAccepted,
			TotalNeeded:   stauts.TotalNeeded,
			QueueType:     in.GetQueueType(),
			Cancelled:     status.Cancelled,
		}
		if err := outStream.Send(resp); err != nil {
			return err
		}
		if resp.GetCancelled() {
			//someone declined, adds user back into queue
			queue.MarkMatchFound(in.GetUserId(), false)
			return nil
		}
		if resp.GetTotalAccepted() == resp.GetTotalNeeded() {
			//everyone accepted, remove user from queue
			//also removes match
			queue.SetMatchStartedAndDelete(in.GetUserId(), in.GetMatchId(), true)
			delete(s.matches, in.GetMatchId())
			return nil
		}
	}
}

func (s *QueueService) Decline(ctx context.Context, in *pb.DeclineQueueRequest) (*pb.DeclineQueueResponse, error) {
	match := s.matches[in.GetMatchId()]
	ch := make(chan MatchResponse, 1)
	if err := match.Decline(in.GetUserId()); err != nil {
		return nil, err
	}

	return &DeclineQueueResponse{in.GetUserId(), in.GetQueueType()}, nil
}

func (s *QueueService) Info(ctx context.Context, in *pb.QueueInfoRequest) (*pb.QueueInfoResponse, error) {
	estimatedWaitTime := s.queueTimes.EstimatedWaitTime(in.GetRating(), 100)

	return &pb.QueueInfoResponse{
		uint32(estimatedWaitTime.Seconds()),
	}, nil
}

func (s *QueueStatus) notifyQueueStateChanges() {
	ch := make(chan PubSubMessage, 1)
	s.queues[0].Subscribe(ch)
	for {
		msg := <- ch
		userID := msg.QueueData.UserID
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
			qChs.MatchFoundCh <- msg.QueueData.MatchID
		case PubSubTopicMatchNotFound:
			qChs. MatchStartFail <- struct{}{}
		}
		s.inQueueLock.Unlock()
	}
}