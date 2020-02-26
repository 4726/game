package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/4726/game/services/other/poll/config"
	"github.com/4726/game/services/other/poll/pb"
	"github.com/cenkalti/backoff"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var maxPollChoices = 10
var errServer = errors.New("server error")
var errDoesNotExist = errors.New("does not exist")
var errAlreadyVoted = errors.New("already voted")
var errPollExpired = errors.New("poll expired")

type pollServer struct {
	db  *redis.Client
	cfg config.Config
}

func newPollServer(c config.Config) (*pollServer, error) {
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

	return &pollServer{db, c}, nil
}

func (s *pollServer) Add(ctx context.Context, in *pb.AddPollRequest) (*pb.AddPollResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}

	expireTime := time.Now().Add(time.Minute * time.Duration(in.GetExpireMinutes())).Unix()
	jsonChoices, err := json.Marshal(in.GetChoices())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := s.db.HSet(id.String(), "choices", string(jsonChoices), "expire", strconv.FormatInt(expireTime, 10)).Err(); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.AddPollResponse{
		Id: id.String(),
	}, nil
}

func (s *pollServer) Get(ctx context.Context, in *pb.GetPollRequest) (*pb.GetPollResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}

	res, err := s.db.HGetAll(in.GetPollId()).Result()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	redisChoices, ok := res["choices"]
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, errDoesNotExist.Error())
	}
	var choices []string
	if err := json.Unmarshal([]byte(redisChoices), &choices); err != nil {
		return nil, status.Error(codes.Internal, errServer.Error())
	}

	if len(choices) > maxPollChoices {
		return nil, status.Error(codes.Internal, errServer.Error())
	}

	userChoices := map[string][]uint64{}
	for _, v := range choices {
		userChoices[v] = []uint64{}
	}

	for k, v := range res {
		userID, err := strconv.ParseUint(k, 10, 64)
		if err != nil {
			continue
		}
		users, ok := userChoices[v]
		if !ok {
			continue
		}
		users = append(users, userID)
		userChoices[v] = users
	}

	var totalVotes int
	for _, v := range userChoices {
		totalVotes += len(v)
	}
	percentages := map[string]int{}
	for k, v := range userChoices {
		if totalVotes == 0 {
			percentages[k] = 0
			continue
		}
		percentages[k] = int((float64(len(v)) / float64(totalVotes)) * 100)
	}

	var pollResults []*pb.PollChoice

	for k, v := range userChoices {
		pollResult := &pb.PollChoice{
			Choice:     k,
			Percentage: uint32(percentages[k]),
			Users:      v,
		}
		pollResults = append(pollResults, pollResult)
	}

	var expireMinutes int64
	expire, ok := res["expire"]
	if ok {
		expireTime, err := strconv.ParseInt(expire, 10, 64)
		if err != nil {
			return nil, status.Error(codes.Internal, errServer.Error())
		}
		dur := time.Until(time.Unix(expireTime, 0))
		expireMinutes = int64(dur.Minutes())
	}

	return &pb.GetPollResponse{
		PollId:        in.GetPollId(),
		Results:       pollResults,
		ExpireMinutes: uint64(expireMinutes),
	}, nil
}

func (s *pollServer) Vote(ctx context.Context, in *pb.VotePollRequest) (*pb.VotePollResponse, error) {
	expireRes, err := s.db.HGet(in.GetPollId(), "expire").Result()
	if err != nil {
		if err == redis.Nil {
			res, err := s.db.HSetNX(in.GetPollId(), strconv.FormatUint(in.GetUserId(), 10), in.GetChoice()).Result()
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			if !res {
				return nil, status.Error(codes.FailedPrecondition, errAlreadyVoted.Error())
			}
			return &pb.VotePollResponse{}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	expireTime, err := strconv.ParseInt(expireRes, 10, 64)
	if err != nil {
		return nil, status.Error(codes.Internal, errServer.Error())
	}

	if time.Now().After(time.Unix(expireTime, 0)) {
		return nil, status.Error(codes.FailedPrecondition, errPollExpired.Error())
	}

	//sync issues here should be fine

	res, err := s.db.HSetNX(in.GetPollId(), strconv.FormatUint(in.GetUserId(), 10), in.GetChoice()).Result()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !res {
		return nil, status.Error(codes.FailedPrecondition, errAlreadyVoted.Error())
	}

	return &pb.VotePollResponse{}, nil
}

//Close gracefully stops the server
func (s *pollServer) Close() {
	s.db.Close()
}
