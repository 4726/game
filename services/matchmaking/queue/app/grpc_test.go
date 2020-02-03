package app

import (
	"context"
	"testing"
	"time"

	"github.com/4726/game/services/matchmaking/queue/config"
	"github.com/4726/game/services/matchmaking/queue/engine"
	"github.com/4726/game/services/matchmaking/queue/engine/memory"
	"github.com/4726/game/services/matchmaking/queue/pb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type test struct {
	c       pb.QueueClient
	service *Service
}

func newTest(t testing.TB, conf ...config.Config) *test {
	var cfg config.Config
	if len(conf) < 1 {
		cfg = config.Config{
			Port:                 14000,
			Limit:                10000,
			PerMatch:             10,
			RatingRange:          100,
			AcceptTimeoutSeconds: 20,
			Metrics:              config.MetricsConfig{14001, "/metrics"},
		}
	} else {
		cfg = conf[0]
		cfg.Port = 14000
	}

	service := NewService(cfg)

	go service.Run()
	time.Sleep(time.Second * 2)

	conn, err := grpc.Dial("127.0.0.1:14000", grpc.WithInsecure())
	assert.NoError(t, err)
	c := pb.NewQueueClient(conn)

	return &test{c, service}
}

func (te *test) insert10(t testing.TB) {
	for i := 1; i < 11; i++ {
		in := &pb.JoinQueueRequest{
			UserId: uint64(i),
			Rating: 1000,
		}
		_, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
	}
	time.Sleep(time.Second * 2)
}

func (te *test) teardown() {
	te.service.Close()
}

func TestServiceJoinQueueFull(t *testing.T) {
	te := newTest(t, config.Config{
		Limit:                5,
		PerMatch:             10,
		RatingRange:          100,
		AcceptTimeoutSeconds: 20,
	})
	defer te.teardown()

	for i := 1; i < 6; i++ {
		in := &pb.JoinQueueRequest{
			UserId: uint64(i),
			Rating: 1000,
		}
		te.c.Join(context.Background(), in)
	}
	time.Sleep(time.Second * 2)
	in := &pb.JoinQueueRequest{
		UserId: uint64(6),
		Rating: 1000,
	}
	outStream, err := te.c.Join(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	expectedErr := status.Error(codes.FailedPrecondition, memory.ErrQueueFull.Error())
	assert.Equal(t, expectedErr, err)
	assert.False(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceJoinAlreadyInQueue(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	in := &pb.JoinQueueRequest{
		UserId: 1,
		Rating: 1000,
	}
	te.c.Join(context.Background(), in)
	time.Sleep(time.Second)

	outStream, err := te.c.Join(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	expectedErr := status.Error(codes.FailedPrecondition, memory.ErrAlreadyInQueue.Error())
	assert.Equal(t, expectedErr, err)
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
	resp, err := outStream.Recv()
	assert.NoError(t, err)
	expectedResp := &pb.JoinQueueResponse{
		UserId:          in.GetUserId(),
		MatchId:         uint64(0),
		Found:           false,
		SecondsToAccept: 20,
	}
	assert.Equal(t, expectedResp, resp)

	for i := 1; i < 10; i++ {
		in := &pb.JoinQueueRequest{
			UserId: uint64(i),
			Rating: 1000,
		}
		_, err := te.c.Join(context.Background(), in)
		assert.NoError(t, err)
	}

	resp, err = outStream.Recv()
	assert.NoError(t, err)
	expectedResp = &pb.JoinQueueResponse{
		UserId:          in.GetUserId(),
		MatchId:         uint64(1),
		Found:           true,
		SecondsToAccept: 20,
	}
	assert.Equal(t, expectedResp, resp)
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
	_, err := te.c.Leave(context.Background(), in)
	expectedErr := status.Error(codes.FailedPrecondition, memory.ErrDoesNotExist.Error())
	assert.Equal(t, expectedErr, err)
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
	expectedErr := status.Error(codes.FailedPrecondition, memory.ErrUserNotInMatch.Error())
	assert.Equal(t, expectedErr, err)
	time.Sleep(time.Second * 2)
	assert.False(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceAcceptNotInMatch(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	te.insert10(t)

	in := &pb.AcceptQueueRequest{
		UserId:  11,
		MatchId: 1,
	}
	outStream, err := te.c.Accept(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	expectedErr := status.Error(codes.FailedPrecondition, memory.ErrUserNotInMatch.Error())
	assert.Equal(t, expectedErr, err)
	time.Sleep(time.Second * 2)
	assert.False(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceAcceptAllAccepted(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	te.insert10(t)
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
		Success: true,
	}
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
	time.Sleep(time.Second * 2)

	for i := 1; i < 10; i++ {
		in := &pb.AcceptQueueRequest{
			UserId:  uint64(i),
			MatchId: 1,
		}
		_, err := te.c.Accept(context.Background(), in)
		assert.NoError(t, err)
		time.Sleep(time.Second * 2)
	}

	resps := []*pb.AcceptQueueResponse{}
	expectedResps := []*pb.AcceptQueueResponse{}
	for i := 1; i < 10; i++ {
		expectedResp := &pb.AcceptQueueResponse{
			TotalAccepted: uint32(i),
			TotalNeeded:   10,
			Cancelled:     false,
		}
		resp, err := outStream.Recv()
		assert.NoError(t, err)
		expectedResps = append(expectedResps, expectedResp)
		resps = append(resps, resp)
	}
	assert.ElementsMatch(t, expectedResps, resps)
	time.Sleep(time.Second * 2)

	resp, err := outStream.Recv()
	assert.NoError(t, err)
	expectedResp := &pb.AcceptQueueResponse{
		Success: true,
	}
	assert.Equal(t, expectedResp, resp)
	time.Sleep(time.Second * 2)
	assert.False(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceAcceptOneDeniedBefore(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	te.insert10(t)

	deleteIn := &pb.DeclineQueueRequest{
		UserId:  1,
		MatchId: 1,
	}
	_, err := te.c.Decline(context.Background(), deleteIn)
	assert.NoError(t, err)
	time.Sleep(time.Second * 2)

	in := &pb.AcceptQueueRequest{
		UserId:  10,
		MatchId: 1,
	}
	outStream, err := te.c.Accept(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	expectedErr := status.Error(codes.FailedPrecondition, memory.ErrUserNotInMatch.Error())
	assert.Equal(t, expectedErr, err)
	time.Sleep(time.Second * 2)
	assert.True(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceAcceptOneDeniedAfter(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	te.insert10(t)

	in := &pb.AcceptQueueRequest{
		UserId:  10,
		MatchId: 1,
	}
	outStream, err := te.c.Accept(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	assert.NoError(t, err)
	time.Sleep(time.Second * 2)

	deleteIn := &pb.DeclineQueueRequest{
		UserId:  1,
		MatchId: 1,
	}
	_, err = te.c.Decline(context.Background(), deleteIn)
	assert.NoError(t, err)

	resp, err := outStream.Recv()
	assert.NoError(t, err)
	expectedResp := &pb.AcceptQueueResponse{
		Cancelled: true,
	}
	assert.Equal(t, expectedResp, resp)
	time.Sleep(time.Second * 2)
	assert.True(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceAcceptTimeout(t *testing.T) {
	te := newTest(t, config.Config{
		Limit:                10000,
		PerMatch:             10,
		RatingRange:          100,
		AcceptTimeoutSeconds: 5,
	})
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
	time.Sleep(time.Second * 7)

	in := &pb.AcceptQueueRequest{
		UserId:  10,
		MatchId: 1,
	}
	outStream, err := te.c.Accept(context.Background(), in)
	assert.NoError(t, err)
	_, err = outStream.Recv()
	expectedErr := status.Error(codes.FailedPrecondition, memory.ErrUserNotInMatch.Error())
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
	expectedErr := status.Error(codes.FailedPrecondition, memory.ErrUserNotInMatch.Error())
	assert.Equal(t, expectedErr, err)
	time.Sleep(time.Second * 2)
	assert.False(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceDeclineNotInMatch(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	for i := 1; i < 11; i++ {
		if i == 10 {
			time.Sleep(time.Second * 2)
		}
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
	expectedErr := status.Error(codes.FailedPrecondition, memory.ErrUserNotInMatch.Error())
	assert.Equal(t, expectedErr, err)
	time.Sleep(time.Second * 2)
	assert.False(t, testIsInQueue(te, in.GetUserId()))
}

func TestServiceDecline(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	te.insert10(t)

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

func TestServiceListen(t *testing.T) {
	te := newTest(t)
	defer te.teardown()
	outStream, err := te.c.Listen(context.Background(), &pb.ListenQueueRequest{})
	assert.NoError(t, err)
	te.insert10(t)
	for i := 1; i < 11; i++ {
		in := &pb.AcceptQueueRequest{
			UserId:  uint64(i),
			MatchId: 1,
		}
		_, err := te.c.Accept(context.Background(), in)
		assert.NoError(t, err)
	}
	msg, err := outStream.Recv()
	assert.NoError(t, err)
	var expectedUsers []*pb.QueueUser
	for i := 1; i < 11; i++ {
		expectedUser := &pb.QueueUser{
			UserId: uint64(i),
			Rating: 1000,
		}
		expectedUsers = append(expectedUsers, expectedUser)
	}
	assert.ElementsMatch(t, expectedUsers, msg.User)
	assert.Equal(t, uint64(1), msg.MatchId)
	assertEmptyRecv(t, outStream)
}

func testIsInQueue(te *test, userID uint64) bool {
	all, _ := te.service.qs.q.All()
	userData, ok := all[userID]
	if !ok {
		return false
	}
	return userData.State == engine.QueueStateInQueue
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
