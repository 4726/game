package app

import (
	"context"
	"fmt"
	"strconv"

	"github.com/4726/game/services/matchmaking/ranking/config"
	"github.com/4726/game/services/matchmaking/ranking/pb"
	"github.com/go-redis/redis/v7"
	"github.com/nsqio/go-nsq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//historyServer implements pb.HistoryServer
type rankingServer struct {
	consumer *nsq.Consumer
	db       *redis.Client
	cfg      config.Config
}

func newRankingServer(c config.Config) (*rankingServer, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})

	if err := db.Ping().Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %v", err)
	}

	consumer, err := nsq.NewConsumer(c.NSQ.Topic, c.NSQ.Channel, nsq.NewConfig())
	if err != nil {
		return nil, fmt.Errorf("could not create nsq consumer: %v", err)
	}
	consumer.AddHandler(&nsqMessageHandler{db, c.Redis.SetName})
	if err := consumer.ConnectToNSQD(c.NSQ.Addr); err != nil {
		return nil, fmt.Errorf("could not connect to nsqd: %v", err)
	}

	return &rankingServer{consumer, db, c}, nil
}

func (s *rankingServer) Get(ctx context.Context, in *pb.GetRankingRequest) (*pb.GetRankingResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}

	rank, err := s.db.ZRevRank(s.cfg.Redis.SetName, strconv.FormatUint(in.GetUserId(), 10)).Result()
	if err != nil {
		if err == redis.Nil {
			return &pb.GetRankingResponse{
				UserId: in.GetUserId(),
				Rating: 0,
				Rank:   0,
			}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	rank++ //position is 0-based so need to add 1

	rating, err := s.db.ZScore(s.cfg.Redis.SetName, strconv.FormatUint(in.GetUserId(), 10)).Result()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetRankingResponse{
		UserId: in.GetUserId(),
		Rating: uint64(rating),
		Rank:   uint64(rank),
	}, nil
}

func (s *rankingServer) GetTop(ctx context.Context, in *pb.GetTopRankingRequest) (*pb.GetTopRankingResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}

	start := in.GetSkip()
	stop := start + in.GetLimit() - 1
	if stop < 0 {
		stop = 0
	}

	var ratings []*pb.GetRankingResponse

	res, err := s.db.ZRevRange(s.cfg.Redis.SetName, int64(start), int64(stop)).Result()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	for _, v := range res {
		userID, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		rank, err := s.db.ZRevRank(s.cfg.Redis.SetName, v).Result()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		rank++

		rating, err := s.db.ZScore(s.cfg.Redis.SetName, v).Result()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		ratings = append(ratings, &pb.GetRankingResponse{
			UserId: userID,
			Rating: uint64(rating),
			Rank:   uint64(rank),
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