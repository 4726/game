package app

import (
	"context"
	"testing"
	"time"

	"github.com/4726/game/services/matchmaking/live/config"
	"github.com/4726/game/services/matchmaking/live/pb"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type test struct {
	c       pb.LiveClient
	service *Service
}

func newTest(t testing.TB) *test {
	cfg := config.Config{
		DB:      config.DBConfig{"live_test", "collection_test"},
		Port:    14000,
		Metrics: config.MetricsConfig{14001, "/metrics"},
	}
	service := NewService(cfg)

	go service.Run()
	time.Sleep(time.Second * 2)

	conn, err := grpc.Dial("127.0.0.1:14000", grpc.WithInsecure())
	assert.NoError(t, err)
	c := pb.NewLiveClient(conn)

	collection := service.ls.db.Database("live_test").Collection("collection_test")
	assert.NoError(t, collection.Drop(context.Background()))

	return &test{c, service}
}

func (te *test) teardown() {
	te.service.Close()
}

func TestServiceGetNone(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.GetLiveRequest{
		MatchId: 1,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.GetLiveResponse{}
	assert.Equal(t, expectedResp, resp)
}

func TestServiceGet(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	t1 := &pb.TeamLiveInfo{
		Users:         []uint64{1, 2, 3, 4, 5},
		Score:         10,
		AverageRating: 2000,
	}
	t2 := &pb.TeamLiveInfo{
		Users:         []uint64{6, 7, 8, 9, 10},
		Score:         10,
		AverageRating: 2020,
	}
	addIn := &pb.AddLiveRequest{
		MatchId:   1,
		Team1:     t1,
		Team2:     t2,
		StartTime: ptypes.TimestampNow(),
	}
	_, err := te.c.Add(context.Background(), addIn)
	assert.NoError(t, err)

	in := &pb.GetLiveRequest{
		MatchId: 1,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.GetLiveResponse{
		MatchId:   addIn.GetMatchId(),
		Team1:     addIn.GetTeam1(),
		Team2:     addIn.GetTeam2(),
		StartTime: addIn.GetStartTime(),
	}
	expectedResp.Team1.XXX_sizecache = 0
	expectedResp.Team2.XXX_sizecache = 0
	expectedResp.StartTime.XXX_sizecache = 0
	assert.Equal(t, expectedResp, resp)
}

func TestServiceGetTotalNone(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.GetTotalLiveRequest{}
	resp, err := te.c.GetTotal(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.GetTotalLiveResponse{Total: 0}
	assert.Equal(t, expectedResp, resp)
}

func TestServiceGetTotal(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddLiveRequest{
		MatchId:   1,
		Team1:     &pb.TeamLiveInfo{},
		Team2:     &pb.TeamLiveInfo{},
		StartTime: ptypes.TimestampNow(),
	}
	_, err := te.c.Add(context.Background(), addIn)
	assert.NoError(t, err)

	in := &pb.GetTotalLiveRequest{}
	resp, err := te.c.GetTotal(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.GetTotalLiveResponse{Total: 1}
	assert.Equal(t, expectedResp, resp)
}

func TestServiceFindMultipleNone(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddLiveRequest{
		MatchId:   1,
		Team1:     &pb.TeamLiveInfo{AverageRating: 3000},
		Team2:     &pb.TeamLiveInfo{AverageRating: 3050},
		StartTime: ptypes.TimestampNow(),
	}
	_, err := te.c.Add(context.Background(), addIn)
	assert.NoError(t, err)

	in := &pb.FindMultipleLiveRequest{
		Total:       10,
		RatingOver:  0,
		RatingUnder: 2000,
	}
	resp, err := te.c.FindMultiple(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.FindMultipleLiveResponse{Matches: nil}
	assert.Equal(t, expectedResp, resp)
}

func TestServiceFindMultiple(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddLiveRequest{
		MatchId:   1,
		Team1:     &pb.TeamLiveInfo{AverageRating: 3000},
		Team2:     &pb.TeamLiveInfo{AverageRating: 3050},
		StartTime: ptypes.TimestampNow(),
	}
	_, err := te.c.Add(context.Background(), addIn)
	assert.NoError(t, err)

	in := &pb.FindMultipleLiveRequest{
		Total:       10,
		RatingOver:  0,
		RatingUnder: 3001,
	}
	resp, err := te.c.FindMultiple(context.Background(), in)
	assert.NoError(t, err)
	expectedMatch := &pb.GetLiveResponse{
		MatchId:   addIn.GetMatchId(),
		Team1:     addIn.GetTeam1(),
		Team2:     addIn.GetTeam2(),
		StartTime: addIn.GetStartTime(),
	}
	expectedResp := &pb.FindMultipleLiveResponse{
		Matches: []*pb.GetLiveResponse{expectedMatch},
	}
	for i, v := range expectedResp.Matches {
		actual := resp.GetMatches()[i]
		v.Team1.XXX_sizecache = 0
		v.Team2.XXX_sizecache = 0
		v.StartTime.XXX_sizecache = 0
		assert.Equal(t, v, actual)
	}
}

func TestServiceFindUserNone(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	t1 := &pb.TeamLiveInfo{
		Users:         []uint64{1, 2, 3, 4, 5},
		Score:         10,
		AverageRating: 2000,
	}
	t2 := &pb.TeamLiveInfo{
		Users:         []uint64{6, 7, 8, 9, 10},
		Score:         10,
		AverageRating: 2020,
	}
	addIn := &pb.AddLiveRequest{
		MatchId:   1,
		Team1:     t1,
		Team2:     t2,
		StartTime: ptypes.TimestampNow(),
	}
	_, err := te.c.Add(context.Background(), addIn)
	assert.NoError(t, err)

	in := &pb.FindUserLiveRequest{
		UserId: 11,
	}
	resp, err := te.c.FindUser(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.FindUserLiveResponse{UserId: in.GetUserId()}
	assert.Equal(t, expectedResp, resp)
}

func TestServiceFindUser(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	t1 := &pb.TeamLiveInfo{
		Users:         []uint64{1, 2, 3, 4, 5},
		Score:         10,
		AverageRating: 2000,
	}
	t2 := &pb.TeamLiveInfo{
		Users:         []uint64{6, 7, 8, 9, 10},
		Score:         10,
		AverageRating: 2020,
	}
	addIn := &pb.AddLiveRequest{
		MatchId:   1,
		Team1:     t1,
		Team2:     t2,
		StartTime: ptypes.TimestampNow(),
	}
	_, err := te.c.Add(context.Background(), addIn)
	assert.NoError(t, err)

	in := &pb.FindUserLiveRequest{
		UserId: 1,
	}
	resp, err := te.c.FindUser(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.FindUserLiveResponse{
		UserId:  in.GetUserId(),
		MatchId: 1,
	}
	assert.Equal(t, expectedResp, resp)
}

func TestServiceAddDuplicate(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	t1 := &pb.TeamLiveInfo{
		Users:         []uint64{1, 2, 3, 4, 5},
		Score:         10,
		AverageRating: 2000,
	}
	t2 := &pb.TeamLiveInfo{
		Users:         []uint64{6, 7, 8, 9, 10},
		Score:         10,
		AverageRating: 2020,
	}
	addIn := &pb.AddLiveRequest{
		MatchId:   1,
		Team1:     t1,
		Team2:     t2,
		StartTime: ptypes.TimestampNow(),
	}
	_, err := te.c.Add(context.Background(), addIn)
	assert.NoError(t, err)

	resp, err := te.c.Add(context.Background(), addIn)
	assert.NoError(t, err)
	expectedResp := &pb.AddLiveResponse{}
	assert.Equal(t, expectedResp, resp)

	totalResp, err := te.c.GetTotal(context.Background(), &pb.GetTotalLiveRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), totalResp.GetTotal())
}

func TestServiceAdd(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	t1 := &pb.TeamLiveInfo{
		Users:         []uint64{1, 2, 3, 4, 5},
		Score:         10,
		AverageRating: 2000,
	}
	t2 := &pb.TeamLiveInfo{
		Users:         []uint64{6, 7, 8, 9, 10},
		Score:         10,
		AverageRating: 2020,
	}
	addIn := &pb.AddLiveRequest{
		MatchId:   1,
		Team1:     t1,
		Team2:     t2,
		StartTime: ptypes.TimestampNow(),
	}
	resp, err := te.c.Add(context.Background(), addIn)
	assert.NoError(t, err)
	expectedResp := &pb.AddLiveResponse{}
	assert.Equal(t, expectedResp, resp)
}

func TestServiceRemoveDoesNotExist(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddLiveRequest{
		MatchId:   1,
		Team1:     &pb.TeamLiveInfo{},
		Team2:     &pb.TeamLiveInfo{},
		StartTime: ptypes.TimestampNow(),
	}
	_, err := te.c.Add(context.Background(), addIn)
	assert.NoError(t, err)

	in := &pb.RemoveLiveRequest{
		MatchId: 2,
	}
	resp, err := te.c.Remove(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.RemoveLiveResponse{}
	assert.Equal(t, expectedResp, resp)

	totalResp, err := te.c.GetTotal(context.Background(), &pb.GetTotalLiveRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), totalResp.GetTotal())
}

func TestServiceRemove(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	t1 := &pb.TeamLiveInfo{
		Users:         []uint64{1, 2, 3, 4, 5},
		Score:         10,
		AverageRating: 2000,
	}
	t2 := &pb.TeamLiveInfo{
		Users:         []uint64{6, 7, 8, 9, 10},
		Score:         10,
		AverageRating: 2020,
	}
	addIn := &pb.AddLiveRequest{
		MatchId:   1,
		Team1:     t1,
		Team2:     t2,
		StartTime: ptypes.TimestampNow(),
	}
	_, err := te.c.Add(context.Background(), addIn)
	assert.NoError(t, err)

	in := &pb.RemoveLiveRequest{
		MatchId: 1,
	}
	resp, err := te.c.Remove(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.RemoveLiveResponse{}
	assert.Equal(t, expectedResp, resp)

	totalResp, err := te.c.GetTotal(context.Background(), &pb.GetTotalLiveRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), totalResp.GetTotal())
}
