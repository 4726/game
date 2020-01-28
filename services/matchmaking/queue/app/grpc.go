package app

import (
	"context"
	"fmt"
	"time"

	"github.com/4726/game/services/matchmaking/queue/config"
	"github.com/4726/game/services/matchmaking/queue/pb"
	"github.com/4726/game/services/matchmaking/queue/queue"
	"github.com/4726/game/services/matchmaking/queue/queue/inmemory"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type queueServer struct {
	q          queue.Queue
	queueTimes *QueueTimes
}

func newQueueServer(cfg config.Config) *queueServer {
	q := inmemory.New(cfg.Limit, cfg.PerMatch, int(cfg.RatingRange), time.Second*time.Duration(cfg.AcceptTimeoutSeconds))
	qs := &queueServer{
		q:          q,
		queueTimes: NewQueueTimes(1000),
	}

	go func(matchFoundCh <-chan queue.Match) {
		for {
			msg, ok := <-matchFoundCh
			if !ok {
				return
			}
			fmt.Println("new match started: ", msg.MatchID)
		}
	}(q.Channel())

	return qs
}

func (s *queueServer) Join(in *pb.JoinQueueRequest, outStream pb.Queue_JoinServer) error {
	ch, err := s.q.Join(in.GetUserId(), in.GetRating())
	if err != nil {
		if err == inmemory.ErrQueueFull || err == inmemory.ErrAlreadyInQueue {
			return status.Error(codes.FailedPrecondition, err.Error())
		}
		return status.Error(codes.Internal, err.Error())
	}

	for {
		msg, ok := <-ch
		if !ok {
			return nil
		}
		switch msg.State {
		case queue.JoinStateEntered:
			resp := &pb.JoinQueueResponse{
				UserId:          in.GetUserId(),
				MatchId:         0,
				Found:           false,
				SecondsToAccept: 20,
			}
			if err := outStream.Send(resp); err != nil {
				s.q.Leave(in.GetUserId())
				return err
			}
		case queue.JoinStateLeft:
			return nil
		case queue.JoinStateGroupFound:
			data := msg.Data.(queue.JoinStateGroupFoundData)
			resp := &pb.JoinQueueResponse{
				UserId:          in.GetUserId(),
				MatchId:         data.MatchID,
				Found:           true,
				SecondsToAccept: 20,
			}
			if err := outStream.Send(resp); err != nil {
				s.q.Leave(in.GetUserId())
				return err
			}
		default:
			return status.Error(codes.Internal, "unknown join state")
		}
	}
}

func (s *queueServer) Leave(ctx context.Context, in *pb.LeaveQueueRequest) (*pb.LeaveQueueResponse, error) {
	err := s.q.Leave(in.GetUserId())
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &pb.LeaveQueueResponse{
		UserId: in.GetUserId(),
	}, nil
}

func (s *queueServer) Accept(in *pb.AcceptQueueRequest, outStream pb.Queue_AcceptServer) error {
	ch, err := s.q.Accept(in.GetUserId(), in.GetMatchId())
	if err != nil {
		if err == inmemory.ErrUserNotInMatch || err == inmemory.ErrUserAlreadyAccepted {
			return status.Error(codes.FailedPrecondition, err.Error())
		}
		return status.Error(codes.Internal, err.Error())
	}

	for {
		msg, ok := <-ch
		if !ok {
			return nil
		}
		switch msg.State {
		case queue.AcceptStateUpdate:
			data := msg.Data.(queue.AcceptStatusUpdateData)
			resp := &pb.AcceptQueueResponse{
				TotalAccepted: uint32(data.UsersAccepted),
				TotalNeeded:   uint32(data.UsersNeeded),
				Cancelled:     false,
			}
			if err := outStream.Send(resp); err != nil {
				return err
			}
		case queue.AcceptStateFailed:
			resp := &pb.AcceptQueueResponse{
				Cancelled: true,
			}
			return outStream.Send(resp)
		case queue.AcceptStateExpired:
			resp := &pb.AcceptQueueResponse{
				Cancelled: true,
			}
			return outStream.Send(resp)
		case queue.AcceptStateSuccess:
			resp := &pb.AcceptQueueResponse{
				Success: true,
			}
			return outStream.Send(resp)
		default:
			return status.Error(codes.Internal, "unknown join state")
		}
	}
}

func (s *queueServer) Decline(ctx context.Context, in *pb.DeclineQueueRequest) (*pb.DeclineQueueResponse, error) {
	if err := s.q.Decline(in.GetUserId(), in.GetMatchId()); err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &pb.DeclineQueueResponse{
		UserId: in.GetUserId(),
	}, nil
}

func (s *queueServer) Info(ctx context.Context, in *pb.QueueInfoRequest) (*pb.QueueInfoResponse, error) {
	usersInQueue, err := s.q.All()
	if err != nil {
		return nil, err
	}
	return &pb.QueueInfoResponse{
		SecondsEstimated: uint32(s.queueTimes.EstimatedWaitTime(in.GetRating(), 100).Seconds()),
		UserCount:        uint32(len(usersInQueue)),
	}, nil
}
