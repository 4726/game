package app

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"

	"github.com/4726/game/services/matchmaking/history/pb"
	"github.com/golang/protobuf/jsonpb"
	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type test struct {
	s       *grpc.Server
	c       pb.HistoryClient
	l       net.Listener
	service *Service
}

func newTest(t testing.TB) *test {
	cfg := Config{
		DBConfig{"history_test", "collection_test"},
		NSQConfig{"127.0.0.1:4150", "matches_test", "db_test"},
	}
	service, err := NewService(cfg)
	assert.NoError(t, err)

	lis, err := net.Listen("tcp", "127.0.0.1:14000")
	assert.NoError(t, err)
	server := grpc.NewServer()
	pb.RegisterHistoryServer(server, service)
	go server.Serve(lis)
	time.Sleep(time.Second * 2)

	conn, err := grpc.Dial("127.0.0.1:14000", grpc.WithInsecure())
	assert.NoError(t, err)
	c := pb.NewHistoryClient(conn)

	collection := service.db.Database("history_test").Collection("collection_test")
	assert.NoError(t, collection.Drop(context.Background()))

	return &test{server, c, lis, service}
}

func (te *test) addData(t testing.TB) []*pb.MatchHistoryInfo {
	producer, err := nsq.NewProducer("127.0.0.1:4150", nsq.NewConfig())
	assert.NoError(t, err)
	defer producer.Stop()

	testMatches := matchHistoryInfoTestData()
	for _, v := range testMatches {
		marshaler := &jsonpb.Marshaler{}
		buffer := bytes.NewBuffer([]byte{})
		assert.NoError(t, marshaler.Marshal(buffer, v))
		assert.NoError(t, producer.Publish("matches_test", buffer.Bytes()))
	}
	time.Sleep(time.Second * 5) //wait for consumer to process message
	return testMatches
}

func (te *test) teardown() {
	te.s.Stop()
	te.l.Close()
	te.service.Close()
}

func TestServiceGetOne(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	matchesAdded := te.addData(t)

	in := &pb.GetHistoryRequest{
		Total: 1,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.GetHistoryResponse{
		Match: []*pb.MatchHistoryInfo{matchesAdded[0]},
	}
	assert.Equal(t, expectedResp, resp)
}

func TestServiceGet(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	matchesAdded := te.addData(t)

	in := &pb.GetHistoryRequest{
		Total: 10,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.GetHistoryResponse{
		Match: matchesAdded,
	}
	assert.Equal(t, expectedResp, resp)
}

func TestServiceGetUserNone(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	te.addData(t)

	in := &pb.GetUserHistoryRequest{
		UserId: 21,
		Total:  10,
	}
	resp, err := te.c.GetUser(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.GetUserHistoryResponse{
		UserId: in.GetUserId(),
	}
	assert.Equal(t, expectedResp, resp)
}

func TestServiceGetUser(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	matchesAdded := te.addData(t)

	in := &pb.GetUserHistoryRequest{
		UserId: 1,
		Total:  10,
	}
	resp, err := te.c.GetUser(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.GetUserHistoryResponse{
		Match:  []*pb.MatchHistoryInfo{matchesAdded[0], matchesAdded[1]},
		UserId: in.GetUserId(),
	}
	assert.Equal(t, expectedResp, resp)
}
