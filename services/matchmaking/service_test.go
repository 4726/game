package main

import (
	"context"
	"io"
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

func TestServiceJoinMatchFound(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	for i := 1; i < 10; i++ {
		in := &pb.JoinQueueRequest{
			UserId:    uint64(i),
			QueueType: pb.QueueType_UNRANKED,
			Rating:    1000,
		}
		_, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 2)

	in := &pb.JoinQueueRequest{
		UserId:    10,
		QueueType: pb.QueueType_UNRANKED,
		Rating:    1000,
	}
	outStream, err := te.c.Join(context.Background(), in)
	assert.NoError(t, err)
	resp, err := outStream.Recv()
	assert.NoError(t, err)
	resp2, err := outStream.Recv()
	assert.NoError(t, err)
	expectedResp := &pb.JoinQueueResponse{
		UserId:          in.GetUserId(),
		QueueType:       in.GetQueueType(),
		MatchId:         uint64(1),
		Found:           true,
		SecondsToAccept: 20,
	}
	expectedResp2 := &pb.JoinQueueResponse{
		UserId:          in.GetUserId(),
		QueueType:       in.GetQueueType(),
		MatchId:         uint64(0),
		Found:           false,
		SecondsToAccept: 20,
	}
	assert.ElementsMatch(t, []*pb.JoinQueueResponse{resp, resp2}, []*pb.JoinQueueResponse{expectedResp, expectedResp2})
	assertEmptyRecv(t, outStream)
}

func TestServiceJoinMatchFoundLater(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	in := &pb.JoinQueueRequest{
		UserId:    10,
		QueueType: pb.QueueType_UNRANKED,
		Rating:    1000,
	}
	outStream, err := te.c.Join(context.Background(), in)
	assert.NoError(t, err)

	for i := 1; i < 10; i++ {
		in := &pb.JoinQueueRequest{
			UserId:    uint64(i),
			QueueType: pb.QueueType_UNRANKED,
			Rating:    1000,
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
		QueueType:       in.GetQueueType(),
		MatchId:         uint64(1),
		Found:           true,
		SecondsToAccept: 20,
	}
	expectedResp2 := &pb.JoinQueueResponse{
		UserId:          in.GetUserId(),
		QueueType:       in.GetQueueType(),
		MatchId:         uint64(0),
		Found:           false,
		SecondsToAccept: 20,
	}
	assert.ElementsMatch(t, []*pb.JoinQueueResponse{resp, resp2}, []*pb.JoinQueueResponse{expectedResp, expectedResp2})
	assertEmptyRecv(t, outStream)
}

func TestServiceJoinMultipleQueues(t *testing.T) {
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

	in2 := &pb.JoinQueueRequest{
		UserId:    1,
		QueueType: pb.QueueType_RANKED,
		Rating:    1000,
	}
	outStream2, err := te.c.Join(context.Background(), in2)
	assert.NoError(t, err)
	resp2, err := outStream2.Recv()
	assert.NoError(t, err)
	expectedResp2 := &pb.JoinQueueResponse{
		UserId:          in2.GetUserId(),
		QueueType:       in2.GetQueueType(),
		MatchId:         uint64(0),
		Found:           false,
		SecondsToAccept: 20,
	}
	assert.Equal(t, expectedResp2, resp2)
	assertEmptyRecv(t, outStream2)
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
	te := newTest(t)
	defer te.teardown()
	in := &pb.LeaveQueueRequest{
		UserId:    1,
		QueueType: pb.QueueType_UNRANKED,
	}
	resp, err := te.c.Leave(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.LeaveQueueResponse{
		UserId:    in.GetUserId(),
		QueueType: in.GetQueueType(),
	}
	assert.Equal(t, expectedResp, resp)
}

func TestServiceLeave(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	joinIn := &pb.JoinQueueRequest{
		UserId:    1,
		QueueType: pb.QueueType_UNRANKED,
		Rating:    1000,
	}
	te.c.Join(context.Background(), joinIn)

	in := &pb.LeaveQueueRequest{
		UserId:    1,
		QueueType: pb.QueueType_UNRANKED,
	}
	resp, err := te.c.Leave(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.LeaveQueueResponse{
		UserId:    in.GetUserId(),
		QueueType: in.GetQueueType(),
	}
	assert.Equal(t, expectedResp, resp)
}

func TestServiceAcceptMatchDoesNotExist(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.AcceptQueueRequest{
		UserId:    1,
		QueueType: pb.QueueType_UNRANKED,
		MatchId:   1,
	}
	outStream, err := te.c.Accept(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	expectedErr := status.Error(codes.FailedPrecondition, ErrUserNotInMatch.Error())
	assert.Equal(t, expectedErr, err)
}

func TestServiceAcceptNotInMatch(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	for i := 1; i < 11; i++ {
		in := &pb.JoinQueueRequest{
			UserId:    uint64(i),
			QueueType: pb.QueueType_UNRANKED,
			Rating:    1000,
		}
		_, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 2)

	in := &pb.AcceptQueueRequest{
		UserId:    11,
		QueueType: pb.QueueType_UNRANKED,
		MatchId:   1,
	}
	outStream, err := te.c.Accept(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	expectedErr := status.Error(codes.FailedPrecondition, ErrUserNotInMatch.Error())
	assert.Equal(t, expectedErr, err)
}

func TestServiceAcceptAllAccepted(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	for i := 1; i < 11; i++ {
		in := &pb.JoinQueueRequest{
			UserId:    uint64(i),
			QueueType: pb.QueueType_UNRANKED,
			Rating:    1000,
		}
		_, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 2)
	for i := 1; i < 10; i++ {
		in := &pb.AcceptQueueRequest{
			UserId:    uint64(i),
			QueueType: pb.QueueType_UNRANKED,
			MatchId:   1,
		}
		_, err := te.c.Accept(context.Background(), in)
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 2)

	in := &pb.AcceptQueueRequest{
		UserId:    10,
		QueueType: pb.QueueType_UNRANKED,
		MatchId:   1,
	}
	outStream, err := te.c.Accept(context.Background(), in)
	assert.NoError(t, err)
	resp, err := outStream.Recv()
	assert.NoError(t, err)
	expectedResp := &pb.AcceptQueueResponse{
		TotalAccepted: 10,
		TotalNeeded:   10,
		QueueType:     in.GetQueueType(),
		Cancelled:     false,
		UserIds:       []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}
	assert.ElementsMatch(t, expectedResp.GetUserIds(), resp.GetUserIds())
	expectedResp.UserIds = nil
	resp.UserIds = nil
	assert.Equal(t, expectedResp, resp)
}

func TestServiceAcceptAllAcceptedLater(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	for i := 1; i < 11; i++ {
		in := &pb.JoinQueueRequest{
			UserId:    uint64(i),
			QueueType: pb.QueueType_UNRANKED,
			Rating:    1000,
		}
		outStream, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
		_, err = outStream.Recv()
		assert.NoError(t, err)
	}

	in := &pb.AcceptQueueRequest{
		UserId:    10,
		QueueType: pb.QueueType_UNRANKED,
		MatchId:   1,
	}
	outStream, err := te.c.Accept(context.Background(), in)
	assert.NoError(t, err)
	time.Sleep(time.Second * 2) //sleep so userID(10) can subscribe before anyone else accepts

	for i := 1; i < 10; i++ {
		in := &pb.AcceptQueueRequest{
			UserId:    uint64(i),
			QueueType: pb.QueueType_UNRANKED,
			MatchId:   1,
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
			QueueType:     in.GetQueueType(),
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
}

func TestServiceAcceptOneDeniedBefore(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	for i := 1; i < 11; i++ {
		in := &pb.JoinQueueRequest{
			UserId:    uint64(i),
			QueueType: pb.QueueType_UNRANKED,
			Rating:    1000,
		}
		_, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 2)

	deleteIn := &pb.DeclineQueueRequest{
		UserId:    1,
		QueueType: pb.QueueType_UNRANKED,
		MatchId:   1,
	}
	_, err := te.c.Decline(context.Background(), deleteIn)
	assert.NoError(t, err)

	in := &pb.AcceptQueueRequest{
		UserId:    10,
		QueueType: pb.QueueType_UNRANKED,
		MatchId:   1,
	}
	outStream, err := te.c.Accept(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	expectedErr := status.Error(codes.FailedPrecondition, ErrUserNotInMatch.Error())
	assert.Equal(t, expectedErr, err)
}

func TestServiceAcceptOneDeniedAfter(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	for i := 1; i < 11; i++ {
		in := &pb.JoinQueueRequest{
			UserId:    uint64(i),
			QueueType: pb.QueueType_UNRANKED,
			Rating:    1000,
		}
		_, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 2)

	in := &pb.AcceptQueueRequest{
		UserId:    10,
		QueueType: pb.QueueType_UNRANKED,
		MatchId:   1,
	}
	outStream, err := te.c.Accept(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	assert.NoError(t, err)

	deleteIn := &pb.DeclineQueueRequest{
		UserId:    1,
		QueueType: pb.QueueType_UNRANKED,
		MatchId:   1,
	}
	_, err = te.c.Decline(context.Background(), deleteIn)
	assert.NoError(t, err)

	_, err = outStream.Recv()
	assert.NoError(t, err)
	_, err = outStream.Recv()
	assert.Equal(t, io.EOF, err)
}

func TestServiceAcceptTimeout(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	defaultMatchAcceptTimeout = time.Second * 3
	defer func() {
		defaultMatchAcceptTimeout = time.Second * 20
	}()

	for i := 1; i < 11; i++ {
		in := &pb.JoinQueueRequest{
			UserId:    uint64(i),
			QueueType: pb.QueueType_UNRANKED,
			Rating:    1000,
		}
		outStream, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
		_, err = outStream.Recv()
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 5)

	in := &pb.AcceptQueueRequest{
		UserId:    10,
		QueueType: pb.QueueType_UNRANKED,
		MatchId:   1,
	}
	outStream, err := te.c.Accept(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	expectedErr := status.Error(codes.FailedPrecondition, ErrMatchCancelled.Error())
	assert.Equal(t, expectedErr, err)
}

func TestServiceDeclineMatchDoesNotExist(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	in := &pb.DeclineQueueRequest{
		UserId:    1,
		QueueType: pb.QueueType_UNRANKED,
		MatchId:   1,
	}
	_, err := te.c.Decline(context.Background(), in)
	expectedErr := status.Error(codes.FailedPrecondition, ErrUserNotInMatch.Error())
	assert.Equal(t, expectedErr, err)
}

func TestServiceDeclineNotInMatch(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	for i := 1; i < 11; i++ {
		in := &pb.JoinQueueRequest{
			UserId:    uint64(i),
			QueueType: pb.QueueType_UNRANKED,
			Rating:    1000,
		}
		_, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 2)

	in := &pb.DeclineQueueRequest{
		UserId:    11,
		QueueType: pb.QueueType_UNRANKED,
		MatchId:   1,
	}
	_, err := te.c.Decline(context.Background(), in)
	expectedErr := status.Error(codes.FailedPrecondition, ErrUserNotInMatch.Error())
	assert.Equal(t, expectedErr, err)
}

func TestServiceDecline(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	for i := 1; i < 11; i++ {
		in := &pb.JoinQueueRequest{
			UserId:    uint64(i),
			QueueType: pb.QueueType_UNRANKED,
			Rating:    1000,
		}
		_, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 2)

	in := &pb.DeclineQueueRequest{
		UserId:    1,
		QueueType: pb.QueueType_UNRANKED,
		MatchId:   1,
	}
	resp, err := te.c.Decline(context.Background(), in)
	expectedResp := &pb.DeclineQueueResponse{
		UserId:    in.GetUserId(),
		QueueType: in.GetQueueType(),
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
}

func TestServiceInfo(t *testing.T) {
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
