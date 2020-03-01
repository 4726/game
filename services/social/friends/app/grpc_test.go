package app

import (
	"context"
	"testing"
	"time"

	"sync"

	"github.com/4726/game/services/social/friends/config"
	"github.com/4726/game/services/social/friends/pb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type test struct {
	c       pb.FriendsClient
	service *Service
}

func newTest(t testing.TB, conf ...config.Config) *test {
	var cfg config.Config

	if len(conf) < 1 {
		cfg = config.Config{
			Port:    14000,
			Metrics: config.MetricsConfig{14001, "/metrics"},
			DB:      config.DBConfig{"postgres", "postgres", "postgres", "localhost", 5432},
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

	service.s.db.Exec("TRUNCATE requests")
	service.s.db.Exec("ALTER SEQUENCE requests_id_seq RESTART WITH 1")

	return &test{c, service}
}

func (te *test) teardown() {
	te.service.Close()
}

func (te *test) add(t testing.TB, from, to uint64) {
	in := &pb.AddFriendRequest{
		UserId: from,
		FriendId: to,
	}
	_, err := te.c.Add(context.Background(), in)
	assert.NoError(t, err)
}

func (te *test) accept(t testing.TB, from, to uint64) {
	in := &pb.AcceptFriendRequest{
		UserId: from,
		FriendId: to,
	}
	_, err := te.c.Accept(context.Background(), in)
	assert.NoError(t, err)
}


func (te *test) fillData(t testing.TB) []Request {
	te.add(t, 1, 2)
	te.accept(t, 2, 1)
	te.add(t, 4, 1)
	te.accept(t, 1, 4)
	te.add(t, 3, 4)
	te.accept(t, 4, 3)
	te.add(t, 2, 3)
	te.add(t, 1, 5)
	return te.queryAll(t)
}

func (te *test) queryAll(t testing.TB) []Request {
	var requests []Request
	assert.NoError(t, te.s.db.Find(&requests).Error)
	return requests
}

func (te *test) assertRequestsEqual(t testing.TB, expected, actual []Request) {
	for i, v := range actual {
		v.ID = 0
		actual[i] = v
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestAddAlreadyFriends(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.AddFriendRequest{
		UserId: 1,
		FriendId: 2,
	}
	_, err := te.c.Add(context.Background(), in)
	assert.Error(t, err)
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestAddAlreadyFriendsReverse(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.AddFriendRequest{
		UserId: 2,
		FriendId: 1,
	}
	_, err := te.c.Add(context.Background(), in)
	assert.Error(t, err)
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestAddHasPendingRequest(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.AddFriendRequest{
		UserId: 3,
		FriendId: 2,
	}
	_, err := te.c.Add(context.Background(), in)
	assert.Error(t, err)
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestAddAlreadyRequested(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.AddFriendRequest{
		UserId: 2,
		FriendId: 3,
	}
	_, err := te.c.Add(context.Background(), in)
	assert.Error(t, err)
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestAdd(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests :+ te.fillData(t)
	in := &pb.AddFriendRequest{
		UserId: 4,
		FriendId: 5,
	}
	_, err := te.c.Add(context.Background(), in)
	assert.NoError(t, err)
	requests = append(requests, Request{
		From: in.GetUserId(),
		To: in.GetFriendId(),
		Accepted: false,
	})
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestDeleteNotFriends(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.AddFriendRequest{
		UserId: 2,
		FriendId: 4,
	}
	_, err := te.c.Add(context.Background(), in)
	assert.Error(t, err)
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestDeleteHasPendingRequest(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.DeleteFriendRequest{
		UserId: 3,
		FriendId: 2,
	}
	_, err := te.c.Delete(context.Background(), in)
	assert.Error(t, err)
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestDeleteAlreadyRequested(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.DeleteFriendRequest{
		UserId: 2,
		FriendId: 3,
	}
	_, err := te.c.Delete(context.Background(), in)
	assert.Error(t, err)
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestDeleteReverse(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.DeleteFriendRequest{
		UserId: 2,
		FriendId: 1,
	}
	_, err := te.c.Delete(context.Background(), in)
	assert.NoError(t, err)
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestDelete(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.DeleteFriendRequest{
		UserId: 1,
		FriendId: 2,
	}
	_, err := te.c.Delete(context.Background(), in)
	assert.NoError(t, err)
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestGetNone(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.GetFriendRequest{
		UserId: 5,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	assert.Len(t, resp.GetFriends(), 0)
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestGet(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.GetFriendRequest{
		UserId: 1,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []uint64{1, 4}, resp.GetFriends())
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestGetRequestsNone(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.GetFriendRequest{
		UserId: 2,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	assert.Len(t, resp.GetFriends(), 0)
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestGetRequests(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.GetFriendRequest{
		UserId: 5,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []uint64{1}, resp.GetFriends())
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestAcceptAlreadyFriends(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.AcceptFriendRequest{
		UserId: 1,
		FriendId: 2,
	}
	_, err := te.c.Accept(context.Background(), in)
	assert.Error(t, err)
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestAcceptNoRequest(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.AcceptFriendRequest{
		UserId: 1,
		FriendId: 3,
	}
	_, err := te.c.Accept(context.Background(), in)
	assert.Error(t, err)
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestAcceptAlreadyRequsted(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.AcceptFriendRequest{
		UserId: 1,
		FriendId: 5,
	}
	_, err := te.c.Accept(context.Background(), in)
	assert.Error(t, err)
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestAccept(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.AcceptFriendRequest{
		UserId: 5,
		FriendId: 1,
	}
	_, err := te.c.Accept(context.Background(), in)
	assert.NoError(t, err)
	for i, v := range requests {
		if v.From == in.GetUserId() && v.To == in.GetFriendId() {
			v.Accepted = true
			requests[i] = v
			break
		}
	}
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestDenyAlreadyFriends(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.DenyFriendRequest{
		UserId: 1,
		FriendId: 2,
	}
	_, err := te.c.Deny(context.Background(), in)
	assert.Error(t, err)
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestDenyNoRequest(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.DenyFriendRequest{
		UserId: 1,
		FriendId: 3,
	}
	_, err := te.c.Deny(context.Background(), in)
	assert.Error(t, err)
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestDenyAlreadyRequsted(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.DenyFriendRequest{
		UserId: 1,
		FriendId: 5,
	}
	_, err := te.c.Deny(context.Background(), in)
	assert.Error(t, err)
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}

func TestDeny(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	requests := te.fillData(t)
	in := &pb.DenyFriendRequest{
		UserId: 5,
		FriendId: 1,
	}
	_, err := te.c.Deny(context.Background(), in)
	assert.NoError(t, err)
	for i, v := range requests {
		if v.From == in.GetUserId() && v.To == in.GetFriendId() {
			requests = requests[:i+copy(requests[i:], requests[i+1:])]
			break
		}
	}
	te.assertRequestsEqual(t, requests, te.queryAll(t))
}