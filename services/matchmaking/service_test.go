package main

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/4726/game/services/matchmaking/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type test struct {
	s *grpc.Server
	c pb.QueueClient
	l net.Listener
}

func newTest(t testing.TB) *test {
	lis, err := net.Listen("tcp", "127.0.0.1:14000")
	assert.NoError(t, err)
	server := grpc.NewServer()
	service := NewQueueService(QueueServiceOptions{100, 10})
	pb.RegisterQueueServer(server, service)
	go server.Serve(lis)
	time.Sleep(time.Second * 2)

	conn, err := grpc.Dial("127.0.0.1:14000", grpc.WithInsecure())
	assert.NoError(t, err)
	c := pb.NewQueueClient(conn)

	return &test{server, c, lis}
}

func (te *test) teardown() {
	te.s.Stop()
	te.l.Close()
}

func TestServiceJoinAlreadyInQueue(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	in := &pb.JoinQueueRequest{
		UserId:    1,
		QueueType: pb.QueueType_UNRANKED,
		Rating:    1000,
	}
	outStream, err := te.c.Join(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	assert.NoError(t, err)

	outStream, err = te.c.Join(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
}

func TestServiceJoinFoundMatch(t *testing.T) {
}

func TestServiceJoinFoundMatchLater(t *testing.T) {
}

func TestServiceJoin(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	in := &pb.JoinQueueRequest{
		UserId:    1,
		QueueType: pb.QueueType_UNRANKED,
		Rating:    1000,
	}
	outStream, err := te.c.Join(context.Background(), in)
	assert.NoError(t, err)
	resp, err := outStream.Recv()
	assert.NoError(t, err)
	expectedResp := &pb.JoinQueueResponse{
		UserId:          in.GetUserId(),
		QueueType:       in.GetQueueType(),
		MatchId:         uint64(0),
		Found:           false,
		SecondsToAccept: 20,
	}
	assert.Equal(t, expectedResp, resp)
	assertEmptyRecv(t, outStream)
}

func TestServiceLeaveNotInQueue(t *testing.T) {
}

func TestServiceLeave(t *testing.T) {
}

func TestServiceAcceptInvalidMatchID(t *testing.T) {
}

func TestServiceAcceptAllAccepted(t *testing.T) {
}

func TestServiceAcceptOneDenied(t *testing.T) {
}

func TestServiceAcceptTimeout(t *testing.T) {
}

func TestServiceDeclineInvalidMatchID(t *testing.T) {
}

func TestServiceInfo(t *testing.T) {
}

func assertEmptyRecv(t testing.TB, stream grpc.ClientStream) {
	ch := make(chan struct{}, 1)
	go func() {
		m := map[string]interface{}{}
		stream.RecvMsg(&m)
		ch <- struct{}{}
	}()
	select {
	case <-time.After(time.Second * 2):
	case <-ch:
		assert.Fail(t, "did not expect another message to stream")
	}
}
