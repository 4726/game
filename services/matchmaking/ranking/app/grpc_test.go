package app

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/4726/game/services/matchmaking/ranking/config"
	"github.com/4726/game/services/matchmaking/ranking/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type test struct {
	c            pb.RankingClient
	service      *Service
	redisSetName string
}

func newTest(t testing.TB) *test {
	cfg := config.Config{
		Redis:   config.RedisConfig{"localhost:6379", "", 0, "ranking_test"},
		NSQ:     config.NSQConfig{"127.0.0.1:4150", "ranking_test", "test"},
		Port:    14000,
		Metrics: config.MetricsConfig{14001, "/metrics"},
	}
	service, err := NewService(cfg)
	fmt.Println(err)
	assert.NoError(t, err)

	go service.Run()
	time.Sleep(time.Second * 2)

	conn, err := grpc.Dial("127.0.0.1:14000", grpc.WithInsecure())
	assert.NoError(t, err)
	c := pb.NewRankingClient(conn)

	assert.NoError(t, service.rs.db.Del(cfg.Redis.SetName).Err())

	return &test{c, service, cfg.Redis.SetName}
}

func (te *test) addData(t testing.TB, r Rating) {
	err := te.service.rs.db.ZAdd(te.redisSetName, r.ToRedisZ()).Err()
	assert.NoError(t, err)
}

func (te *test) teardown() {
	te.service.Close()
}

func TestServiceGetDoesNotExist(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.GetRankingRequest{
		UserId: 1,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.GetRankingResponse{
		UserId: in.GetUserId(),
		Rating: 0,
		Rank:   0,
	}
	assert.Equal(t, expectedResp, resp)
}

func TestServiceGetSameRating(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	te.addData(t, Rating{1, 5000})
	te.addData(t, Rating{2, 2000})
	te.addData(t, Rating{3, 2000})

	in := &pb.GetRankingRequest{
		UserId: 2,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.GetRankingResponse{
		UserId: 2,
		Rating: 2000,
		Rank:   3,
	}
	assert.Equal(t, expectedResp, resp)
}

func TestServiceGet(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	te.addData(t, Rating{1, 5000})
	te.addData(t, Rating{2, 2000})
	te.addData(t, Rating{3, 2000})

	in := &pb.GetRankingRequest{
		UserId: 1,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.GetRankingResponse{
		UserId: 1,
		Rating: 5000,
		Rank:   1,
	}
	assert.Equal(t, expectedResp, resp)
}

func TestServiceGetTopNone(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.GetTopRankingRequest{
		Limit: 10,
		Skip:  0,
	}
	resp, err := te.c.GetTop(context.Background(), in)
	assert.NoError(t, err)
	assert.Len(t, resp.GetRatings(), 0)
}

func TestServiceGetTopSkip(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	te.addData(t, Rating{1, 5000})
	te.addData(t, Rating{2, 2000})
	te.addData(t, Rating{3, 2005})
	te.addData(t, Rating{4, 6000})

	in := &pb.GetTopRankingRequest{
		Limit: 10,
		Skip:  2,
	}
	resp, err := te.c.GetTop(context.Background(), in)
	assert.NoError(t, err)

	r1 := &pb.GetRankingResponse{
		UserId: 3,
		Rating: 2005,
		Rank:   3,
	}
	r2 := &pb.GetRankingResponse{
		UserId: 2,
		Rating: 2000,
		Rank:   4,
	}
	expectedResp := &pb.GetTopRankingResponse{
		Ratings: []*pb.GetRankingResponse{r1, r2},
	}
	assert.Equal(t, expectedResp, resp)
}

func TestServiceGetTopLimit(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	te.addData(t, Rating{1, 5000})
	te.addData(t, Rating{2, 2000})
	te.addData(t, Rating{3, 2005})
	te.addData(t, Rating{4, 6000})

	in := &pb.GetTopRankingRequest{
		Limit: 3,
		Skip:  0,
	}
	resp, err := te.c.GetTop(context.Background(), in)
	assert.NoError(t, err)

	r1 := &pb.GetRankingResponse{
		UserId: 4,
		Rating: 6000,
		Rank:   1,
	}
	r2 := &pb.GetRankingResponse{
		UserId: 1,
		Rating: 5000,
		Rank:   2,
	}
	r3 := &pb.GetRankingResponse{
		UserId: 3,
		Rating: 2005,
		Rank:   3,
	}
	expectedResp := &pb.GetTopRankingResponse{
		Ratings: []*pb.GetRankingResponse{r1, r2, r3},
	}
	assert.Equal(t, expectedResp, resp)
}

func TestServiceGetTop(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	te.addData(t, Rating{1, 5000})
	te.addData(t, Rating{2, 2000})
	te.addData(t, Rating{3, 2005})
	te.addData(t, Rating{4, 6000})

	in := &pb.GetTopRankingRequest{
		Limit: 10,
		Skip:  0,
	}
	resp, err := te.c.GetTop(context.Background(), in)
	assert.NoError(t, err)

	r1 := &pb.GetRankingResponse{
		UserId: 4,
		Rating: 6000,
		Rank:   1,
	}
	r2 := &pb.GetRankingResponse{
		UserId: 1,
		Rating: 5000,
		Rank:   2,
	}
	r3 := &pb.GetRankingResponse{
		UserId: 3,
		Rating: 2005,
		Rank:   3,
	}
	r4 := &pb.GetRankingResponse{
		UserId: 2,
		Rating: 2000,
		Rank:   4,
	}
	expectedResp := &pb.GetTopRankingResponse{
		Ratings: []*pb.GetRankingResponse{r1, r2, r3, r4},
	}
	assert.Equal(t, expectedResp, resp)
}
