package app

import (
	"context"
	"testing"
	"time"

	"github.com/4726/game/services/constants/constants/config"
	"github.com/4726/game/services/constants/constants/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type test struct {
	c       pb.ConstantsClient
	service *Service
}

func newTest(t testing.TB) *test {
	cfg := config.Config{
		Redis:   config.RedisConfig{"localhost:6379", "", 0},
		Port:    14000,
		Metrics: config.MetricsConfig{14001, "/metrics"},
	}
	service, err := NewService(cfg)
	assert.NoError(t, err)

	go service.Run()
	time.Sleep(time.Second * 2)

	conn, err := grpc.Dial("127.0.0.1:14000", grpc.WithInsecure())
	assert.NoError(t, err)
	c := pb.NewConstantsClient(conn)

	assert.NoError(t, service.cs.db.FlushAll().Err())

	return &test{c, service}
}

func (te *test) addData(t testing.TB, k, v string) {
	err := te.service.cs.db.Set(k, v, 0).Err()
	assert.NoError(t, err)
}

func (te *test) teardown() {
	te.service.Close()
}

func TestServiceGetOneDoesNotExist(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.GetConstantsRequest{
		Keys: []string{"key1"},
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	item1 := &pb.ConstantItem{
		Key:   "key1",
		Value: "",
	}
	assert.Equal(t, []*pb.ConstantItem{item1}, resp.GetItems())
}

func TestServiceGetOne(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	te.addData(t, "key1", "value1")

	in := &pb.GetConstantsRequest{
		Keys: []string{"key1"},
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	item1 := &pb.ConstantItem{
		Key:   "key1",
		Value: "value1",
	}
	assert.Equal(t, []*pb.ConstantItem{item1}, resp.GetItems())
}

func TestServiceGetManyAllDoesNotExist(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.GetConstantsRequest{
		Keys: []string{"key1", "key2"},
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	item1 := &pb.ConstantItem{
		Key:   "key1",
		Value: "",
	}
	item2 := &pb.ConstantItem{
		Key:   "key2",
		Value: "",
	}
	assert.Equal(t, []*pb.ConstantItem{item1, item2}, resp.GetItems())
}

func TestServiceGetManySomeDoNotExist(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	te.addData(t, "key1", "value1")

	in := &pb.GetConstantsRequest{
		Keys: []string{"key1", "key2"},
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	item1 := &pb.ConstantItem{
		Key:   "key1",
		Value: "value1",
	}
	item2 := &pb.ConstantItem{
		Key:   "key2",
		Value: "",
	}
	assert.Equal(t, []*pb.ConstantItem{item1, item2}, resp.GetItems())
}

func TestServiceGetMany(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	te.addData(t, "key1", "value1")
	te.addData(t, "key2", "value2")

	in := &pb.GetConstantsRequest{
		Keys: []string{"key1", "key2"},
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	item1 := &pb.ConstantItem{
		Key:   "key1",
		Value: "value1",
	}
	item2 := &pb.ConstantItem{
		Key:   "key2",
		Value: "value2",
	}
	assert.ElementsMatch(t, []*pb.ConstantItem{item1, item2}, resp.GetItems())
}

func TestServiceGetAllNone(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	resp, err := te.c.GetAll(context.Background(), &pb.GetAllConstantsRequest{})
	assert.NoError(t, err)
	expectedResp := &pb.GetAllConstantsResponse{}
	assert.Equal(t, expectedResp, resp)
}

func TestServiceGetAll(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	te.addData(t, "key1", "value1")
	te.addData(t, "key2", "value2")

	resp, err := te.c.GetAll(context.Background(), &pb.GetAllConstantsRequest{})
	assert.NoError(t, err)
	item1 := &pb.ConstantItem{
		Key:   "key1",
		Value: "value1",
	}
	item2 := &pb.ConstantItem{
		Key:   "key2",
		Value: "value2",
	}
	assert.ElementsMatch(t, []*pb.ConstantItem{item1, item2}, resp.GetItems())
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
