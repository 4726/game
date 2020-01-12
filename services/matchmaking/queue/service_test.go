package main

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/4726/game/services/matchmaking/queue/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type test struct {
	s  *grpc.Server
	c  pb.QueueClient
	l  net.Listener
	qs *QueueService
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

	return &test{server, c, lis, service}
}

func (te *test) teardown() {
	te.s.Stop()
	te.l.Close()
}

func TestServiceJoinAlreadyInQueue(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	in := &pb.JoinQueueRequest{
		UserId: 1,
		Rating: 1000,
	}
	outStream, err := te.c.Join(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	assert.NoError(t, err)

	outStream, err = te.c.Join(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	time.Sleep(time.Second * 2)
	assert.True(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceJoinMatchFound(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	for i := 1; i < 10; i++ {
		in := &pb.JoinQueueRequest{
			UserId: uint64(i),
			Rating: 1000,
		}
		_, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 2)

	in := &pb.JoinQueueRequest{
		UserId: 10,
		Rating: 1000,
	}
	outStream, err := te.c.Join(context.Background(), in)
	assert.NoError(t, err)
	resp, err := outStream.Recv()
	assert.NoError(t, err)
	resp2, err := outStream.Recv()
	assert.NoError(t, err)
	expectedResp := &pb.JoinQueueResponse{
		UserId:          in.GetUserId(),
		MatchId:         uint64(1),
		Found:           true,
		SecondsToAccept: 20,
	}
	expectedResp2 := &pb.JoinQueueResponse{
		UserId:          in.GetUserId(),
		MatchId:         uint64(0),
		Found:           false,
		SecondsToAccept: 20,
	}
	assert.ElementsMatch(t, []*pb.JoinQueueResponse{resp, resp2}, []*pb.JoinQueueResponse{expectedResp, expectedResp2})
	assertEmptyRecv(t, outStream)
	time.Sleep(time.Second * 2)
	assert.False(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceJoinMatchFoundLater(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	in := &pb.JoinQueueRequest{
		UserId: 10,
		Rating: 1000,
	}
	outStream, err := te.c.Join(context.Background(), in)
	assert.NoError(t, err)

	for i := 1; i < 10; i++ {
		in := &pb.JoinQueueRequest{
			UserId: uint64(i),
			Rating: 1000,
		}
		_, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
	}

	resp, err := outStream.Recv()
	assert.NoError(t, err)
	resp2, err := outStream.Recv()
	assert.NoError(t, err)
	expectedResp := &pb.JoinQueueResponse{
		UserId:          in.GetUserId(),
		MatchId:         uint64(1),
		Found:           true,
		SecondsToAccept: 20,
	}
	expectedResp2 := &pb.JoinQueueResponse{
		UserId:          in.GetUserId(),
		MatchId:         uint64(0),
		Found:           false,
		SecondsToAccept: 20,
	}
	assert.ElementsMatch(t, []*pb.JoinQueueResponse{resp, resp2}, []*pb.JoinQueueResponse{expectedResp, expectedResp2})
	assertEmptyRecv(t, outStream)
	time.Sleep(time.Second * 2)
	assert.False(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceJoin(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	in := &pb.JoinQueueRequest{
		UserId: 1,
		Rating: 1000,
	}
	outStream, err := te.c.Join(context.Background(), in)
	assert.NoError(t, err)
	resp, err := outStream.Recv()
	assert.NoError(t, err)
	expectedResp := &pb.JoinQueueResponse{
		UserId:          in.GetUserId(),
		MatchId:         uint64(0),
		Found:           false,
		SecondsToAccept: 20,
	}
	assert.Equal(t, expectedResp, resp)
	assertEmptyRecv(t, outStream)
	time.Sleep(time.Second * 2)
	assert.True(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceLeaveNotInQueue(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	in := &pb.LeaveQueueRequest{
		UserId: 1,
	}
	resp, err := te.c.Leave(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.LeaveQueueResponse{
		UserId: in.GetUserId(),
	}
	assert.Equal(t, expectedResp, resp)
	time.Sleep(time.Second * 2)
	assert.False(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceLeave(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	joinIn := &pb.JoinQueueRequest{
		UserId: 1,
		Rating: 1000,
	}
	te.c.Join(context.Background(), joinIn)
	time.Sleep(time.Second * 2)

	in := &pb.LeaveQueueRequest{
		UserId: 1,
	}
	resp, err := te.c.Leave(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.LeaveQueueResponse{
		UserId: in.GetUserId(),
	}
	assert.Equal(t, expectedResp, resp)
	time.Sleep(time.Second * 2)
	assert.False(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceAcceptMatchDoesNotExist(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.AcceptQueueRequest{
		UserId:  1,
		MatchId: 1,
	}
	outStream, err := te.c.Accept(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	expectedErr := status.Error(codes.FailedPrecondition, ErrUserNotInMatch.Error())
	assert.Equal(t, expectedErr, err)
	time.Sleep(time.Second * 2)
	assert.False(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceAcceptNotInMatch(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	for i := 1; i < 11; i++ {
		in := &pb.JoinQueueRequest{
			UserId: uint64(i),
			Rating: 1000,
		}
		_, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 2)

	in := &pb.AcceptQueueRequest{
		UserId:  11,
		MatchId: 1,
	}
	outStream, err := te.c.Accept(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	expectedErr := status.Error(codes.FailedPrecondition, ErrUserNotInMatch.Error())
	assert.Equal(t, expectedErr, err)
	time.Sleep(time.Second * 2)
	assert.False(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceAcceptAllAccepted(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	for i := 1; i < 11; i++ {
		in := &pb.JoinQueueRequest{
			UserId: uint64(i),
			Rating: 1000,
		}
		_, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 2)
	for i := 1; i < 10; i++ {
		in := &pb.AcceptQueueRequest{
			UserId:  uint64(i),
			MatchId: 1,
		}
		_, err := te.c.Accept(context.Background(), in)
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 2)

	in := &pb.AcceptQueueRequest{
		UserId:  10,
		MatchId: 1,
	}
	outStream, err := te.c.Accept(context.Background(), in)
	assert.NoError(t, err)
	resp, err := outStream.Recv()
	assert.NoError(t, err)
	expectedResp := &pb.AcceptQueueResponse{
		TotalAccepted: 10,
		TotalNeeded:   10,
		Cancelled:     false,
		UserIds:       []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}
	assert.ElementsMatch(t, expectedResp.GetUserIds(), resp.GetUserIds())
	expectedResp.UserIds = nil
	resp.UserIds = nil
	assert.Equal(t, expectedResp, resp)
	time.Sleep(time.Second * 2)
	assert.False(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceAcceptAllAcceptedLater(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	for i := 1; i < 11; i++ {
		in := &pb.JoinQueueRequest{
			UserId: uint64(i),
			Rating: 1000,
		}
		outStream, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
		_, err = outStream.Recv()
		assert.NoError(t, err)
	}

	in := &pb.AcceptQueueRequest{
		UserId:  10,
		MatchId: 1,
	}
	outStream, err := te.c.Accept(context.Background(), in)
	assert.NoError(t, err)
	time.Sleep(time.Second * 2) //sleep so userID(10) can subscribe before anyone else accepts

	for i := 1; i < 10; i++ {
		in := &pb.AcceptQueueRequest{
			UserId:  uint64(i),
			MatchId: 1,
		}
		_, err := te.c.Accept(context.Background(), in)
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 2)

	resps := []*pb.AcceptQueueResponse{}
	expectedResps := []*pb.AcceptQueueResponse{}
	for i := 1; i < 11; i++ {
		expectedResp := &pb.AcceptQueueResponse{
			TotalAccepted: uint32(i),
			TotalNeeded:   10,
			Cancelled:     false,
			UserIds:       []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		}
		resp, err := outStream.Recv()
		assert.NoError(t, err)
		assert.ElementsMatch(t, expectedResp.GetUserIds(), resp.GetUserIds())
		resp.UserIds = nil
		expectedResp.UserIds = nil
		expectedResps = append(expectedResps, resp)
		resps = append(resps, resp)
	}

	assert.ElementsMatch(t, expectedResps, resps)
	time.Sleep(time.Second * 2)
	assert.False(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceAcceptOneDeniedBefore(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	for i := 1; i < 11; i++ {
		in := &pb.JoinQueueRequest{
			UserId: uint64(i),
			Rating: 1000,
		}
		_, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 2)

	deleteIn := &pb.DeclineQueueRequest{
		UserId:  1,
		MatchId: 1,
	}
	_, err := te.c.Decline(context.Background(), deleteIn)
	assert.NoError(t, err)

	in := &pb.AcceptQueueRequest{
		UserId:  10,
		MatchId: 1,
	}
	outStream, err := te.c.Accept(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	expectedErr := status.Error(codes.FailedPrecondition, ErrUserNotInMatch.Error())
	assert.Equal(t, expectedErr, err)
	time.Sleep(time.Second * 2)
	assert.True(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceAcceptOneDeniedAfter(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	for i := 1; i < 11; i++ {
		in := &pb.JoinQueueRequest{
			UserId: uint64(i),
			Rating: 1000,
		}
		_, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 2)

	in := &pb.AcceptQueueRequest{
		UserId:  10,
		MatchId: 1,
	}
	outStream, err := te.c.Accept(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	assert.NoError(t, err)

	deleteIn := &pb.DeclineQueueRequest{
		UserId:  1,
		MatchId: 1,
	}
	_, err = te.c.Decline(context.Background(), deleteIn)
	assert.NoError(t, err)

	_, err = outStream.Recv()
	assert.NoError(t, err)
	_, err = outStream.Recv()
	assert.Equal(t, io.EOF, err)
	time.Sleep(time.Second * 2)
	assert.True(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceAcceptTimeout(t *testing.T) {
	oldDefaultMatchAcceptTimeout := defaultMatchAcceptTimeout
	oldDefaultCheckTimeout := defaultCheckTimeout
	defaultMatchAcceptTimeout = time.Second * 3
	defaultCheckTimeout = time.Second
	defer func() {
		defaultMatchAcceptTimeout = oldDefaultMatchAcceptTimeout
		defaultCheckTimeout = oldDefaultCheckTimeout
	}()

	te := newTest(t)
	defer te.teardown()

	for i := 1; i < 11; i++ {
		in := &pb.JoinQueueRequest{
			UserId: uint64(i),
			Rating: 1000,
		}
		outStream, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
		_, err = outStream.Recv()
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 5)

	in := &pb.AcceptQueueRequest{
		UserId:  10,
		MatchId: 1,
	}
	outStream, err := te.c.Accept(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	expectedErr := status.Error(codes.FailedPrecondition, ErrUserNotInMatch.Error())
	assert.Equal(t, expectedErr, err)
	time.Sleep(time.Second * 2)
	assert.False(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceDeclineMatchDoesNotExist(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	in := &pb.DeclineQueueRequest{
		UserId:  1,
		MatchId: 1,
	}
	_, err := te.c.Decline(context.Background(), in)
	expectedErr := status.Error(codes.FailedPrecondition, ErrUserNotInMatch.Error())
	assert.Equal(t, expectedErr, err)
	time.Sleep(time.Second * 2)
	assert.False(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceDeclineNotInMatch(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	for i := 1; i < 11; i++ {
		in := &pb.JoinQueueRequest{
			UserId: uint64(i),
			Rating: 1000,
		}
		_, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 2)

	in := &pb.DeclineQueueRequest{
		UserId:  11,
		MatchId: 1,
	}
	_, err := te.c.Decline(context.Background(), in)
	expectedErr := status.Error(codes.FailedPrecondition, ErrUserNotInMatch.Error())
	assert.Equal(t, expectedErr, err)
	time.Sleep(time.Second * 2)
	assert.False(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceDecline(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	for i := 1; i < 11; i++ {
		in := &pb.JoinQueueRequest{
			UserId: uint64(i),
			Rating: 1000,
		}
		_, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 2)

	in := &pb.DeclineQueueRequest{
		UserId:  1,
		MatchId: 1,
	}
	resp, err := te.c.Decline(context.Background(), in)
	expectedResp := &pb.DeclineQueueResponse{
		UserId: in.GetUserId(),
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
	time.Sleep(time.Second * 2)
	assert.False(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceInfo(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	for i := 1; i < 6; i++ {
		in := &pb.JoinQueueRequest{
			UserId: uint64(i),
			Rating: 1000,
		}
		_, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 2)

	in := &pb.QueueInfoRequest{
		UserId: 6,
		Rating: 1000,
	}
	resp, err := te.c.Info(context.Background(), in)
	expectedResp := &pb.QueueInfoResponse{
		SecondsEstimated: uint32(0),
		UserCount:        5,
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
}

func testIsInQueue(te *test, userID uint64) bool {
	var inQueue bool
	te.qs.queue.ForEach(func(qd QueueData) bool {
		if qd.UserID == userID {
			inQueue = !qd.MatchFound
			return false
		}
		return true
	})
	return inQueue
}

func assertEmptyRecv(t testing.TB, stream grpc.ClientStream) {
	ch := make(chan struct{}, 1)
	go func() {
		m := &unimplementedProtoMessage{}
		stream.RecvMsg(m)
		ch <- struct{}{}
	}()
	select {
	case <-time.After(time.Second * 2):
	case <-ch:
		assert.Fail(t, "did not expect another message to stream")
	}
}

type unimplementedProtoMessage struct{}

func (u *unimplementedProtoMessage) Reset()         {}
func (u *unimplementedProtoMessage) String() string { return "" }
func (u *unimplementedProtoMessage) ProtoMessage()  {}
