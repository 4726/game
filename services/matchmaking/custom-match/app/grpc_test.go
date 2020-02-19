package app

import (
	"context"
	"testing"
	"time"

	"github.com/4726/game/services/matchmaking/custom-match/config"
	"github.com/4726/game/services/matchmaking/custom-match/pb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type test struct {
	c       pb.CustomMatchClient
	service *Service
}

func newTest(t testing.TB, conf ...config.Config) *test {
	var cfg config.Config
	if len(conf) < 1 {
		cfg = config.Config{
			Port:    14000,
			Metrics: config.MetricsConfig{14001, "/metrics"},
		}
	} else {
		cfg = conf[0]
		cfg.Port = 14000
	}

	service, err := NewService(cfg)
	assert.NoError(t, err)

	go service.Run()
	time.Sleep(time.Second * 2)

	conn, err := grpc.Dial("127.0.0.1:14000", grpc.WithInsecure())
	assert.NoError(t, err)
	c := pb.NewCustomMatchClient(conn)

	service.s.db.Exec("TRUNCATE users, groups")
	service.s.db.Exec("ALTER SEQUENCE groups_id_seq RESTART WITH 1")

	return &test{c, service}
}

func (te *test) teardown() {
	te.service.Close()
}

func TestServiceAddLeaderAlreadyHasGroup(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.AddCustomMatchRequest{
		Leader:   1,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}
	outStream, _ := te.c.Add(context.Background(), in)
	go func() {
		for {
			outStream.Recv()
		}
	}()
	time.Sleep(time.Second * 5)

	outStream, err := te.c.Add(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	assert.Error(t, err)
}

func TestServiceGetAll(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddCustomMatchRequest{
		Leader:   1,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}
	outStream, _ := te.c.Add(context.Background(), addIn)
	go func() {
		for {
			outStream.Recv()
		}
	}()
	time.Sleep(time.Second * 5)

	in := &pb.GetAllCustomMatchRequest{Total: 100, Skip: 0}
	resp, err := te.c.GetAll(context.Background(), in)
	assert.NoError(t, err)
	g1 := &pb.CustomMatchGroup{
		GroupId:          1,
		Name:             "game",
		Leader:           1,
		PasswordRequired: false,
		MaxUsers:         10,
		TotalUsers:       1,
	}
	expectedResp := &pb.GetAllCustomMatchResponse{
		Groups: []*pb.CustomMatchGroup{g1},
	}
	assert.Equal(t, expectedResp, resp)
}
