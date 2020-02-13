package app

import (
	"context"
	"fmt"

	"github.com/4726/game/services/constants/constants/config"
	"github.com/4726/game/services/constants/constants/pb"
	"github.com/cenkalti/backoff"
	"github.com/go-redis/redis/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//constantsServer implements pb.ConstantsServer
type constantsServer struct {
	db  *redis.Client
	cfg config.Config
}

func newConstantsServer(c config.Config) (*constantsServer, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})

	redisOp := func() error {
		logEntry.Info("connecting to redis: ", c.Redis.Addr)
		err := db.Ping().Err()
		if err != nil {
			logEntry.Warn("could not connect to redis, retrying")
		}
		return err
	}

	if err := backoff.Retry(redisOp, backoff.NewExponentialBackOff()); err != nil {
		logEntry.Error("could not connect to redis, max retries reached")
		return nil, fmt.Errorf("could not connect to redis: %v", err)
	}

	return &constantsServer{db, c}, nil
}

func (s *constantsServer) Get(ctx context.Context, in *pb.GetConstantsRequest) (*pb.GetConstantsResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}

	items := []*pb.ConstantItem{}
	for _, v := range in.GetKeys() {
		value, err := s.db.Get(v).Result()
		if err != nil {
			if err == redis.Nil {
				items = append(items, &pb.ConstantItem{
					Key:   v,
					Value: "",
				})
				continue
			}
			return nil, status.Error(codes.Internal, err.Error())
		}
		items = append(items, &pb.ConstantItem{
			Key:   v,
			Value: value,
		})
	}

	return &pb.GetConstantsResponse{
		Items: items,
	}, nil
}

func (s *constantsServer) GetAll(ctx context.Context, in *pb.GetAllConstantsRequest) (*pb.GetAllConstantsResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}

	keys, err := s.db.Keys("*").Result()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	items := []*pb.ConstantItem{}
	for _, v := range keys {
		value, err := s.db.Get(v).Result()
		if err != nil {
			if err == redis.Nil {
				items = append(items, &pb.ConstantItem{
					Key:   v,
					Value: "",
				})
				continue
			}
			return nil, status.Error(codes.Internal, err.Error())
		}
		items = append(items, &pb.ConstantItem{
			Key:   v,
			Value: value,
		})
	}

	return &pb.GetAllConstantsResponse{
		Items: items,
	}, nil
}

//Close gracefully stops the server
func (s *constantsServer) Close() {
	s.db.Close()
}
