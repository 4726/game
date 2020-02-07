package app

import (
	"context"
	"fmt"
	"time"

	"github.com/4726/game/services/matchmaking/ranking/config"
	"github.com/4726/game/services/matchmaking/ranking/pb"
	"github.com/nsqio/go-nsq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/go-redis/redis/v7"
)

//historyServer implements pb.HistoryServer
type rankingServer struct {
	consumer *nsq.Consumer
	db       *redis.Client
	cfg      config.Config
}

func newRankingServer(c config.Config) (*rankingServer, error) {
	db := redis.NewClient(&redis.Options{
		Addr: c.Redis.Addr,
		Password: c.Redis.Password,
		DB: c.Redis.DB,
	})

	if err := db.Ping().Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %v", err)
	}

	consumer, err := nsq.NewConsumer(c.NSQ.Topic, c.NSQ.Channel, nsq.NewConfig())
	if err != nil {
		return nil, fmt.Errorf("could not create nsq consumer: %v", err)
	}
	consumer.AddHandler(&nsqMessageHandler{db, c.DB.Name, c.DB.Collection})
	if err := consumer.ConnectToNSQD(c.NSQ.Addr); err != nil {
		return nil, fmt.Errorf("could not connect to nsqd: %v", err)
	}

	return &rankingServer{consumer, db, c}, nil
}

func (s *rankingServer) Get(ctx context.Context, in *pb.GetRankingRequest) (*pb.GetRankingResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}

	rank, err := s.db.ZRevRank(s.cfg.Redis.SetName, strconv.Itoa(in.GetUserId)).Result()
	if err != nil {
		return nil, status.Error(codes.Internal, err)
	}
	rank++ //position is 0-based so need to add 1

	rating, err := s.db.ZScore((s.cfg.Redis.SetName, strconv.Itoa(in.GetUserId)).Result()
	if err != nil {
		return nil, status.Error(codes.Internal, err)
	}

	return &pb.GetRankingResponse{
		UserId: in.GetUserId(),
		Rating: uint64(rating),
		Rank: uint64(rank),
	}, nil
}

func (s *rankingServer) GetTop(ctx context.Context, in *pb.GetTopRankingRequest) (*pb.GetTopRankingResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}

	start := in.GetSkip()
	stop := start + in.GetTotal()

	var ratings []*pb.GetRankingResponse

	res, err := s.db.ZRange((s.cfg.Redis.SetName, int64(start), int64(stop)).Result()
	if err != nil {
		return nil, status.Error(codes.Internal, err)
	}

	for _, v := range res {
		userID, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, status.Error(codes.Internal, err)
		}

		rank, err := s.db.ZRevRank((s.cfg.Redis.SetName, strconv.Itoa(in.GetUserId)).Result()
		if err != nil {
			return nil, status.Error(codes.Internal, err)
		}
		rank++
	
		rating, err := s.db.ZScore((s.cfg.Redis.SetName, strconv.Itoa(in.GetUserId)).Result()
		if err != nil {
			return nil, status.Error(codes.Internal, err)
		}

		ratings = append(ratings, &pb.GetRankingResponse{
			UserId: userID,
			Rating: uint64(rating),
			Rank: uint64(rank),
		})
	}

	return &pb.GetTopRankingResponse{
		Ratings: ratings,
	}, nil
}

//Close gracefully stops the server
func (s *rankingServer) Close() {
	s.db.Close()
	s.consumer.Stop()
	<-s.consumer.StopChan
}
