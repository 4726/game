package app

import (
	"context"
	"fmt"

	"github.com/4726/game/services/other/poll/config"
	"github.com/4726/game/services/other/poll/pb"
	"github.com/cenkalti/backoff"
	"github.com/go-redis/redis/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/google/uuid"
)

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

func (s *pollServer) Add(context.Context, in *pb.AddPollRequest) (*pb.AddPollResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}

	expireTime := time.Now().Add(time.Minute * in.GetExpireMinutes()).Unix()
	jsonChoices, err := json.Marshal(in.GetChoices())
	if err != nil {
		return nil, stats.Error(codes.Internal, "server error")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, status.Error(codes.Internal, err)
	}

	if err := s.db.HSet(id.String(), "choices", string(jsonChoices), "expire", strconv.FormatInt(expireTime, 10)); err != nil {
		return nil, stats.Error(codes.Internal, err)
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
		return nil, status.Error(codes.Internal, err)
	}

	redisChoices, ok := res["choices"]
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "does not exist")
	}
	var choices []string
	if err := json.Unmarshal([]byte(redisChoices), &choices); err != nil {
		return nil, status.Error(codes.Internal, "server error")
	}

	if len(choices) > maxPollChoices {
		return nil, status.Error(codes.Internal, "server error")
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
		percentages[k] = len(v) / totalVotes
	}

	var pollResults []*pb.PollChoiceResult

	for k, v := range userChoices {
		pollResult := &pb.PollChoiceResult{
			Choice: k,
			Percentage: percentages[k],
			Users: v,
		}
		pollResults = append(pollResults, pollResult)
	}

	return &pb.GetPollResponse{
		PollId: in.GetPollId()
		Results: pollResults.
	}, nil
}

func (s *pollServer) Vote(context.Context, in *pb.VotePollRequest) (*pb.VotePollResponse, error) {
	res, err := s.db.HGet(in.GetPollId(), "expire")
	if err != nil {
		return nil, status.Error(codes.Internal, err)
	}

	expireTime, err := strconv.ParseInt(res, 10, 64)
	if err != nil {
		return nil, status.Error(codes.Internal, "server error")
	}

	if time.Now().After(time.Unix(expireTime, 0)) {
		return nil, status.Error(codes.FailedPrecondition, "poll expired")
	}

	//sync issues here should be fine

	res, err := s.db.HSetNX(in.GetPollId(), strconv.FormatUint(in.GetUserId(), 10), in.GetChoice()).Result()
	if err != nil {
		return nil, status.Error(codes.Internal, err)
	}
	if !res {
		return nil, status.Error(codes.FailedPrecondition, "already voted")
	}

	return &pb.VotePollResponse{}, nil
}

//Close gracefully stops the server
func (s *pollServer) Close() {
	s.db.Close()
}
