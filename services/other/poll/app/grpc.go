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

var errServer = errors.New("server error")
var errDoesNotExist = errors.New("does not exist")
var errAlreadyVoted = errors.New("already voted")
var errPollExpired = errors.New("poll expired")
var errMaxChoicesExceeded = errors.New("max choices exceeded")
var errMaxExpireExceeded = errors.New("max expire minutes exceeded")
var errPollDoesNotExist = errors.New("poll does not exist")
var errPollChoiceDoesNotExist = errors.New("poll choice does not exist")

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

	if len(in.GetChoices()) > int(s.cfg.MaxPollChoices) {
		return nil, status.Error(codes.FailedPrecondition, errMaxChoicesExceeded.Error())
	}

	kvs := map[string]interface{}{}

	if in.GetExpireMinutes() > 0 {
		if in.GetExpireMinutes() > int64(s.cfg.MaxExpireMinutes) {
			return nil, status.Error(codes.FailedPrecondition, errMaxExpireExceeded.Error())
		}
		expireTime := time.Now().Add(time.Minute * time.Duration(in.GetExpireMinutes())).Unix()
		kvs["expire"] = strconv.FormatInt(expireTime, 10)
	}

	jsonChoices, err := json.Marshal(in.GetChoices())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	kvs["choices"] = string(jsonChoices)

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := s.db.HSet(id.String(), kvs).Err(); err != nil {
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

	if len(choices) > int(s.cfg.MaxPollChoices) {
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
	var hasExpiration bool
	expire, ok := res["expire"]
	if ok {
		expireTime, err := strconv.ParseInt(expire, 10, 64)
		if err != nil {
			return nil, status.Error(codes.Internal, errServer.Error())
		}
		dur := time.Until(time.Unix(expireTime, 0))
		expireMinutes = int64(dur.Minutes())
		hasExpiration = true
	}

	return &pb.GetPollResponse{
		PollId:        in.GetPollId(),
		Results:       pollResults,
		ExpireMinutes: expireMinutes,
		HasExpiration: hasExpiration,
	}, nil
}

//contains possible sync issues but should be fine
func (s *pollServer) Vote(ctx context.Context, in *pb.VotePollRequest) (*pb.VotePollResponse, error) {
	redisChoices, err := s.db.HGet(in.GetPollId(), "choices").Result()
	if err != nil {
		if err == redis.Nil {
			return nil, status.Error(codes.FailedPrecondition, errPollDoesNotExist.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	var choices []string
	if err := json.Unmarshal([]byte(redisChoices), &choices); err != nil {
		return nil, status.Error(codes.Internal, errServer.Error())
	}

	var choiceExists bool
	for _, v := range choices {
		if v == in.GetChoice() {
			choiceExists = true
			break
		}
	}

	if !choiceExists {
		return nil, status.Error(codes.FailedPrecondition, errPollChoiceDoesNotExist.Error())
	}

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
