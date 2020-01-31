package app

import (
	"context"
	"time"

	"github.com/4726/game/services/matchmaking/queue/config"
	"github.com/4726/game/services/matchmaking/queue/engine"
	"github.com/4726/game/services/matchmaking/queue/engine/memory"
	"github.com/4726/game/services/matchmaking/queue/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//queueServer implements the grpc server
type queueServer struct {
	q          engine.Queue
	queueTimes *queueTimes
}

func newQueueServer(cfg config.Config) *queueServer {
	q := memory.New(cfg.Limit, cfg.PerMatch, int(cfg.RatingRange), time.Second*time.Duration(cfg.AcceptTimeoutSeconds))
	qs := &queueServer{
		q:          q,
		queueTimes: newQueueTimes(1000),
	}

	inQueueTicker := time.NewTicker(time.Minute)
	go func(q engine.Queue) {
		for {
			<-inQueueTicker.C
			total, _ := q.Len()
			inQueueTotal.Set(float64(total))
		}
	}(qs.q)

	return qs
}

func (s *queueServer) Join(in *pb.JoinQueueRequest, outStream pb.Queue_JoinServer) error {
	ch, err := s.q.Join(in.GetUserId(), in.GetRating())
	if err != nil {
		if err == memory.ErrQueueFull || err == memory.ErrAlreadyInQueue {
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
		case engine.JoinStateEntered:
			resp := &pb.JoinQueueResponse{
				UserId:          in.GetUserId(),
				MatchId:         0,
				Found:           false,
				SecondsToAccept: 20,
			}
			if err := outStream.Send(resp); err != nil {
				return err
			}
		case engine.JoinStateLeft:
			return nil
		case engine.JoinStateGroupFound:
			data := msg.Data.(engine.JoinStateGroupFoundData)
			resp := &pb.JoinQueueResponse{
				UserId:          in.GetUserId(),
				MatchId:         data.MatchID,
				Found:           true,
				SecondsToAccept: 20,
			}
			if err := outStream.Send(resp); err != nil {
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
		if err == memory.ErrUserNotInMatch || err == memory.ErrUserAlreadyAccepted {
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
		case engine.AcceptStateUpdate:
			data := msg.Data.(engine.AcceptStatusUpdateData)
			resp := &pb.AcceptQueueResponse{
				TotalAccepted: uint32(data.UsersAccepted),
				TotalNeeded:   uint32(data.UsersNeeded),
				Cancelled:     false,
			}
			if err := outStream.Send(resp); err != nil {
				return err
			}
		case engine.AcceptStateFailed:
			resp := &pb.AcceptQueueResponse{
				Cancelled: true,
			}
			return outStream.Send(resp)
		case engine.AcceptStateExpired:
			resp := &pb.AcceptQueueResponse{
				Cancelled: true,
			}
			return outStream.Send(resp)
		case engine.AcceptStateSuccess:
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
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}

	usersInQueue, err := s.q.All()
	if err != nil {
		return nil, err
	}
	return &pb.QueueInfoResponse{
		SecondsEstimated: uint32(s.queueTimes.EstimatedWaitTime(in.GetRating(), 100).Seconds()),
		UserCount:        uint32(len(usersInQueue)),
	}, nil
}

func (s *queueServer) Listen(in *pb.ListenQueueRequest, outStream pb.Queue_ListenServer) error {
	ch := s.q.Channel()

	for {
		msg, ok := <-ch
		if !ok {
			return nil
		}
		var users []*pb.QueueUser
		for k, v := range msg.Users {
			user := &pb.QueueUser{
				UserId: k,
				Rating: v,
			}
			users = append(users, user)
		}
		resp := &pb.ListenQueueResponse{
			MatchId: msg.MatchID,
			User:    users,
		}
		if err := outStream.Send(resp); err != nil {
			return err
		}
	}
}
