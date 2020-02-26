package app

import (
	"context"
	"testing"
	"time"

	"github.com/4726/game/services/other/poll/config"
	"github.com/4726/game/services/other/poll/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type test struct {
	c       pb.PollClient
	service *Service
}

func newTest(t testing.TB) *test {
	cfg := config.Config{
		Redis:            config.RedisConfig{"localhost:6379", "", 0},
		Port:             14000,
		Metrics:          config.MetricsConfig{14001, "/metrics"},
		MaxPollChoices:   10,
		MaxExpireMinutes: 144000,
	}
	service, err := NewService(cfg)
	assert.NoError(t, err)

	go service.Run()
	time.Sleep(time.Second * 2)

	conn, err := grpc.Dial("127.0.0.1:14000", grpc.WithInsecure())
	assert.NoError(t, err)
	c := pb.NewPollClient(conn)

	assert.NoError(t, service.s.db.FlushAll().Err())

	return &test{c, service}
}

func (te *test) vote(t testing.TB, userID uint64, pollID, choice string) {
	in := &pb.VotePollRequest{
		UserId: userID,
		PollId: pollID,
		Choice: choice,
	}
	_, err := te.c.Vote(context.Background(), in)
	assert.NoError(t, err)
}

func (te *test) teardown() {
	te.service.Close()
}

func TestServiceAddMaxChoices(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.AddPollRequest{
		Creator:       1,
		Choices:       []string{"yes", "no", "maybe", "hi", "bye", "ok", "what", "why", "when", "one", "two"},
		ExpireMinutes: 60,
	}
	_, err := te.c.Add(context.Background(), in)
	assert.Equal(t, status.Error(codes.FailedPrecondition, errMaxChoicesExceeded.Error()), err)
}

func TestServiceAddMaxExpire(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.AddPollRequest{
		Creator:       1,
		Choices:       []string{"yes", "no", "maybe"},
		ExpireMinutes: 144001,
	}
	_, err := te.c.Add(context.Background(), in)
	assert.Equal(t, status.Error(codes.FailedPrecondition, errMaxExpireExceeded.Error()), err)
}

func TestServiceAddNoExpiration(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.AddPollRequest{
		Creator:       1,
		Choices:       []string{"yes", "no", "maybe"},
		ExpireMinutes: 0,
	}
	res, err := te.c.Add(context.Background(), in)
	assert.NoError(t, err)
	assert.NotEqual(t, "", res.GetId())
}

func TestServiceAddNegativeExpiration(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.AddPollRequest{
		Creator:       1,
		Choices:       []string{"yes", "no", "maybe"},
		ExpireMinutes: -1,
	}
	res, err := te.c.Add(context.Background(), in)
	assert.NoError(t, err)
	assert.NotEqual(t, "", res.GetId())
}

func TestServiceAdd(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.AddPollRequest{
		Creator:       1,
		Choices:       []string{"yes", "no", "maybe"},
		ExpireMinutes: 60,
	}
	res, err := te.c.Add(context.Background(), in)
	assert.NoError(t, err)
	assert.NotEqual(t, "", res.GetId())
}

func TestServiceGetDoesNotExist(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.GetPollRequest{
		PollId: "qwqwqwqwq",
	}
	_, err := te.c.Get(context.Background(), in)
	assert.Equal(t, status.Error(codes.FailedPrecondition, errDoesNotExist.Error()), err)
}

func TestServiceGetNoVotes(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddPollRequest{
		Creator:       1,
		Choices:       []string{"yes", "no", "maybe"},
		ExpireMinutes: 60,
	}
	addRes, err := te.c.Add(context.Background(), addIn)
	pollID := addRes.GetId()

	in := &pb.GetPollRequest{
		PollId: pollID,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	c1 := &pb.PollChoice{
		Choice:     "yes",
		Percentage: 0,
	}
	c2 := &pb.PollChoice{
		Choice:     "no",
		Percentage: 0,
	}
	c3 := &pb.PollChoice{
		Choice:     "maybe",
		Percentage: 0,
	}
	assertPollChoices(t, []*pb.PollChoice{c1, c2, c3}, resp.GetResults())
	assert.Contains(t, []int64{60, 59, 58}, resp.GetExpireMinutes())
	assert.Equal(t, true, resp.GetHasExpiration())
}

func TestServiceGetOneVote(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddPollRequest{
		Creator:       1,
		Choices:       []string{"yes", "no", "maybe"},
		ExpireMinutes: 60,
	}
	addRes, err := te.c.Add(context.Background(), addIn)
	pollID := addRes.GetId()

	te.vote(t, 2, pollID, "yes")

	in := &pb.GetPollRequest{
		PollId: pollID,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	c1 := &pb.PollChoice{
		Choice:     "yes",
		Percentage: 100,
		Users:      []uint64{2},
	}
	c2 := &pb.PollChoice{
		Choice:     "no",
		Percentage: 0,
		Users:      []uint64{},
	}
	c3 := &pb.PollChoice{
		Choice:     "maybe",
		Percentage: 0,
		Users:      []uint64{},
	}
	assert.Equal(t, pollID, resp.GetPollId())
	assertPollChoices(t, []*pb.PollChoice{c1, c2, c3}, resp.GetResults())
	assert.Contains(t, []int64{60, 59, 58}, resp.GetExpireMinutes())
	assert.Equal(t, true, resp.GetHasExpiration())
}

func TestServiceGetEqualVotes(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddPollRequest{
		Creator:       1,
		Choices:       []string{"yes", "no", "maybe"},
		ExpireMinutes: 60,
	}
	addRes, err := te.c.Add(context.Background(), addIn)
	pollID := addRes.GetId()

	te.vote(t, 2, pollID, "yes")
	te.vote(t, 3, pollID, "no")
	te.vote(t, 4, pollID, "maybe")
	te.vote(t, 5, pollID, "yes")
	te.vote(t, 6, pollID, "no")
	te.vote(t, 7, pollID, "maybe")

	in := &pb.GetPollRequest{
		PollId: pollID,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	c1 := &pb.PollChoice{
		Choice:     "yes",
		Percentage: 33,
		Users:      []uint64{2, 5},
	}
	c2 := &pb.PollChoice{
		Choice:     "no",
		Percentage: 33,
		Users:      []uint64{3, 6},
	}
	c3 := &pb.PollChoice{
		Choice:     "maybe",
		Percentage: 33,
		Users:      []uint64{4, 7},
	}
	assert.Equal(t, pollID, resp.GetPollId())
	assertPollChoices(t, []*pb.PollChoice{c1, c2, c3}, resp.GetResults())
	assert.Contains(t, []int64{60, 59, 58}, resp.GetExpireMinutes())
	assert.Equal(t, true, resp.GetHasExpiration())
}

func TestServiceGetNoExpiration(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddPollRequest{
		Creator:       1,
		Choices:       []string{"yes", "no", "maybe"},
		ExpireMinutes: 0,
	}
	addRes, err := te.c.Add(context.Background(), addIn)
	pollID := addRes.GetId()

	te.vote(t, 2, pollID, "yes")
	te.vote(t, 3, pollID, "no")
	te.vote(t, 4, pollID, "maybe")
	te.vote(t, 5, pollID, "yes")
	te.vote(t, 6, pollID, "yes")

	in := &pb.GetPollRequest{
		PollId: pollID,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	c1 := &pb.PollChoice{
		Choice:     "yes",
		Percentage: 60,
		Users:      []uint64{2, 5, 6},
	}
	c2 := &pb.PollChoice{
		Choice:     "no",
		Percentage: 20,
		Users:      []uint64{3},
	}
	c3 := &pb.PollChoice{
		Choice:     "maybe",
		Percentage: 20,
		Users:      []uint64{4},
	}
	assert.Equal(t, pollID, resp.GetPollId())
	assertPollChoices(t, []*pb.PollChoice{c1, c2, c3}, resp.GetResults())
	assert.Equal(t, false, resp.GetHasExpiration())
}

func TestServiceGetNegativeExpiration(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddPollRequest{
		Creator:       1,
		Choices:       []string{"yes", "no", "maybe"},
		ExpireMinutes: -1,
	}
	addRes, err := te.c.Add(context.Background(), addIn)
	pollID := addRes.GetId()

	te.vote(t, 2, pollID, "yes")
	te.vote(t, 3, pollID, "no")
	te.vote(t, 4, pollID, "maybe")
	te.vote(t, 5, pollID, "yes")
	te.vote(t, 6, pollID, "yes")

	in := &pb.GetPollRequest{
		PollId: pollID,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	c1 := &pb.PollChoice{
		Choice:     "yes",
		Percentage: 60,
		Users:      []uint64{2, 5, 6},
	}
	c2 := &pb.PollChoice{
		Choice:     "no",
		Percentage: 20,
		Users:      []uint64{3},
	}
	c3 := &pb.PollChoice{
		Choice:     "maybe",
		Percentage: 20,
		Users:      []uint64{4},
	}
	assert.Equal(t, pollID, resp.GetPollId())
	assertPollChoices(t, []*pb.PollChoice{c1, c2, c3}, resp.GetResults())
	assert.Equal(t, false, resp.GetHasExpiration())
}

func TestServiceGet(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddPollRequest{
		Creator:       1,
		Choices:       []string{"yes", "no", "maybe"},
		ExpireMinutes: 60,
	}
	addRes, err := te.c.Add(context.Background(), addIn)
	pollID := addRes.GetId()

	te.vote(t, 2, pollID, "yes")
	te.vote(t, 3, pollID, "no")
	te.vote(t, 4, pollID, "maybe")
	te.vote(t, 5, pollID, "yes")
	te.vote(t, 6, pollID, "yes")

	in := &pb.GetPollRequest{
		PollId: pollID,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	c1 := &pb.PollChoice{
		Choice:     "yes",
		Percentage: 60,
		Users:      []uint64{2, 5, 6},
	}
	c2 := &pb.PollChoice{
		Choice:     "no",
		Percentage: 20,
		Users:      []uint64{3},
	}
	c3 := &pb.PollChoice{
		Choice:     "maybe",
		Percentage: 20,
		Users:      []uint64{4},
	}
	assert.Equal(t, pollID, resp.GetPollId())
	assertPollChoices(t, []*pb.PollChoice{c1, c2, c3}, resp.GetResults())
	assert.Contains(t, []int64{60, 59, 58}, resp.GetExpireMinutes())
	assert.Equal(t, true, resp.GetHasExpiration())
}

func TestServiceVotePollDoesNotExist(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.VotePollRequest{
		UserId: 2,
		PollId: "qwqwqwq",
		Choice: "yes",
	}
	_, err := te.c.Vote(context.Background(), in)
	assert.Equal(t, status.Error(codes.FailedPrecondition, errPollDoesNotExist.Error()), err)
}

func TestServiceVoteExpired(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddPollRequest{
		Creator:       1,
		Choices:       []string{"yes", "no", "maybe"},
		ExpireMinutes: 1,
	}
	addRes, err := te.c.Add(context.Background(), addIn)
	pollID := addRes.GetId()

	time.Sleep(time.Second * 65)
	in := &pb.VotePollRequest{
		UserId: 2,
		PollId: pollID,
		Choice: "yes",
	}
	_, err = te.c.Vote(context.Background(), in)
	assert.Equal(t, status.Error(codes.FailedPrecondition, errPollExpired.Error()), err)
}

func TestServiceVoteChoiceDoesNotExist(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddPollRequest{
		Creator:       1,
		Choices:       []string{"yes", "no", "maybe"},
		ExpireMinutes: 1,
	}
	addRes, err := te.c.Add(context.Background(), addIn)
	pollID := addRes.GetId()

	in := &pb.VotePollRequest{
		UserId: 2,
		PollId: pollID,
		Choice: "other",
	}
	_, err = te.c.Vote(context.Background(), in)
	assert.Equal(t, status.Error(codes.FailedPrecondition, errPollChoiceDoesNotExist.Error()), err)
}

func TestServiceVoteAlreadyVotedSame(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddPollRequest{
		Creator:       1,
		Choices:       []string{"yes", "no", "maybe"},
		ExpireMinutes: 1,
	}
	addRes, err := te.c.Add(context.Background(), addIn)
	pollID := addRes.GetId()

	te.vote(t, 2, pollID, "yes")

	in := &pb.VotePollRequest{
		UserId: 2,
		PollId: pollID,
		Choice: "yes",
	}
	_, err = te.c.Vote(context.Background(), in)
	assert.Equal(t, status.Error(codes.FailedPrecondition, errAlreadyVoted.Error()), err)
}

func TestServiceVoteAlreadyVotedDifferent(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddPollRequest{
		Creator:       1,
		Choices:       []string{"yes", "no", "maybe"},
		ExpireMinutes: 1,
	}
	addRes, err := te.c.Add(context.Background(), addIn)
	pollID := addRes.GetId()

	te.vote(t, 2, pollID, "no")

	in := &pb.VotePollRequest{
		UserId: 2,
		PollId: pollID,
		Choice: "yes",
	}
	_, err = te.c.Vote(context.Background(), in)
	assert.Equal(t, status.Error(codes.FailedPrecondition, errAlreadyVoted.Error()), err)
}

func TestServiceVote(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddPollRequest{
		Creator:       1,
		Choices:       []string{"yes", "no", "maybe"},
		ExpireMinutes: 1,
	}
	addRes, err := te.c.Add(context.Background(), addIn)
	pollID := addRes.GetId()

	in := &pb.VotePollRequest{
		UserId: 2,
		PollId: pollID,
		Choice: "yes",
	}
	_, err = te.c.Vote(context.Background(), in)
	assert.NoError(t, err)
}

func TestServiceTLSInvalidPath(t *testing.T) {
	cfg := config.Config{
		Redis:   config.RedisConfig{"localhost:6379", "", 0},
		Port:    14000,
		Metrics: config.MetricsConfig{14001, "/metrics"},
		TLS:     config.TLSConfig{"crt.pem", "key.pem"},
	}

	_, err := NewService(cfg)
	assert.Error(t, err)
}

func TestServiceTLS(t *testing.T) {
	cfg := config.Config{
		Redis:   config.RedisConfig{"localhost:6379", "", 0},
		Port:    14000,
		Metrics: config.MetricsConfig{14001, "/metrics"},
		TLS:     config.TLSConfig{"../../../../tests/tls/localhost.crt", "../../../../tests/tls/localhost.key"},
	}

	service, err := NewService(cfg)
	assert.NoError(t, err)
	defer service.Close()
}

func assertPollChoices(t testing.TB, expected, actual []*pb.PollChoice) {
	for _, v := range expected {
		var found bool
		for _, v2 := range actual {
			if v.Choice == v2.Choice {
				assert.ElementsMatch(t, v.Users, v2.Users)
				assert.Equal(t, v.Percentage, v2.Percentage)
				found = true
			}
		}

		if !found {
			assert.Fail(t, "not found")
		}
	}
}
