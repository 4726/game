package app

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/4726/game/services/social/chat/config"
	"github.com/4726/game/services/social/chat/pb"
	"github.com/golang/protobuf/ptypes"
	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type test struct {
	c       pb.ChatClient
	service *Service
}

func newTest(t testing.TB) *test {
	cfg := config.Config{
		Cassandra: config.CassandraConfig{"127.0.0.1", 9042, 10},
		Port:      14000,
		Metrics:   config.MetricsConfig{14001, "/metrics"},
		NSQ:       config.NSQConfig{"127.0.0.1:4150", "notification_test"},
	}
	service, err := NewService(cfg)
	assert.NoError(t, err)

	go service.Run()
	time.Sleep(time.Second * 2)

	conn, err := grpc.Dial("127.0.0.1:14000", grpc.WithInsecure())
	assert.NoError(t, err)
	c := pb.NewChatClient(conn)

	assert.NoError(t, service.cs.db.Query("TRUNCATE chat.messages;").Exec())

	return &test{c, service}
}

func (te *test) addData(t testing.TB, from, to uint64, msg string) {
	in := &pb.SendChatRequest{
		From:    from,
		To:      to,
		Message: msg,
	}
	_, err := te.c.Send(context.Background(), in)
	assert.NoError(t, err)
}

func (te *test) teardown() {
	te.service.Close()
}

func TestServiceSend(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.SendChatRequest{
		From:    1,
		To:      2,
		Message: "hello",
	}
	resp, err := te.c.Send(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.SendChatResponse{}
	assert.Equal(t, expectedResp, resp)

	ch := make(chan Message, 1)
	consumer, err := nsq.NewConsumer("notification_test", "testing", nsq.NewConfig())
	assert.NoError(t, err)
	consumer.AddHandler(&testNSQMessageHandler{ch, t})
	err = consumer.ConnectToNSQD("127.0.0.1:4150")
	assert.NoError(t, err)
	var msgs []Message
Loop:
	for {
		select {
		case msg := <-ch:
			msgs = append(msgs, msg)
		case <-time.After(time.Second * 5):
			break Loop
		}
	}

	assert.Equal(t, in.GetFrom(), msgs[0].MessagesFrom)
	assert.Equal(t, in.GetMessage(), msgs[0].MessagesMessage)
	assert.Equal(t, in.GetFrom(), msgs[0].MessagesUser1)
	assert.Equal(t, in.GetTo(), msgs[0].MessagesUser2)
	assert.NotEqual(t, "", msgs[0].MessagesTime.String())
	assert.Len(t, msgs, 1)
}

func TestServiceGetNone(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.GetChatRequest{
		User1: 1,
		User2: 2,
		Total: 100,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	assert.Len(t, resp.GetMessages(), 0)
}

func TestServiceGetTotal(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	te.addData(t, 1, 2, "hello")
	te.addData(t, 2, 1, "hi")
	te.addData(t, 1, 2, "bye bye")
	te.addData(t, 2, 1, "bye")

	in := &pb.GetChatRequest{
		User1: 1,
		User2: 2,
		Total: 2,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	expectedMsg1 := &pb.ChatMessage{
		From:    2,
		To:      1,
		Message: "bye",
		Time:    nil,
	}

	expectedMsg2 := &pb.ChatMessage{
		From:    1,
		To:      2,
		Message: "bye bye",
		Time:    nil,
	}

	var respMsgs []*pb.ChatMessage
	for _, v := range resp.GetMessages() {
		ts, err := ptypes.Timestamp(v.GetTime())
		assert.NoError(t, err)
		assert.WithinDuration(t, time.Now(), ts, time.Second*20)

		v.Time = nil
		respMsgs = append(respMsgs, v)
	}
	assert.Equal(t, []*pb.ChatMessage{expectedMsg1, expectedMsg2}, resp.GetMessages())
}

func TestServiceGetUsersReversed(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	te.addData(t, 1, 2, "hello")
	te.addData(t, 2, 1, "hi")
	te.addData(t, 1, 2, "bye bye")

	in := &pb.GetChatRequest{
		User1: 2,
		User2: 1,
		Total: 100,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	expectedMsg1 := &pb.ChatMessage{
		From:    1,
		To:      2,
		Message: "bye bye",
		Time:    nil,
	}

	expectedMsg2 := &pb.ChatMessage{
		From:    2,
		To:      1,
		Message: "hi",
		Time:    nil,
	}

	expectedMsg3 := &pb.ChatMessage{
		From:    1,
		To:      2,
		Message: "hello",
		Time:    nil,
	}

	var respMsgs []*pb.ChatMessage
	for _, v := range resp.GetMessages() {
		ts, err := ptypes.Timestamp(v.GetTime())
		assert.NoError(t, err)
		assert.WithinDuration(t, time.Now(), ts, time.Second*20)

		v.Time = nil
		respMsgs = append(respMsgs, v)
	}
	assert.Equal(t, []*pb.ChatMessage{expectedMsg1, expectedMsg2, expectedMsg3}, resp.GetMessages())
}

func TestServiceGetSkip(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	te.addData(t, 1, 2, "hello")
	te.addData(t, 2, 1, "hi")
	te.addData(t, 1, 2, "bye bye")
	te.addData(t, 2, 1, "bye")

	in := &pb.GetChatRequest{
		User1: 2,
		User2: 1,
		Total: 100,
		Skip:  2,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	expectedMsg1 := &pb.ChatMessage{
		From:    2,
		To:      1,
		Message: "hi",
		Time:    nil,
	}
	expectedMsg2 := &pb.ChatMessage{
		From:    1,
		To:      2,
		Message: "hello",
		Time:    nil,
	}

	var respMsgs []*pb.ChatMessage
	for _, v := range resp.GetMessages() {
		ts, err := ptypes.Timestamp(v.GetTime())
		assert.NoError(t, err)
		assert.WithinDuration(t, time.Now(), ts, time.Second*20)

		v.Time = nil
		respMsgs = append(respMsgs, v)
	}
	assert.Equal(t, []*pb.ChatMessage{expectedMsg1, expectedMsg2}, resp.GetMessages())
}

func TestServiceGet(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	te.addData(t, 1, 2, "hello")
	te.addData(t, 2, 1, "hi")
	te.addData(t, 1, 2, "bye bye")

	in := &pb.GetChatRequest{
		User1: 1,
		User2: 2,
		Total: 100,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	expectedMsg1 := &pb.ChatMessage{
		From:    1,
		To:      2,
		Message: "bye bye",
		Time:    nil,
	}

	expectedMsg2 := &pb.ChatMessage{
		From:    2,
		To:      1,
		Message: "hi",
		Time:    nil,
	}

	expectedMsg3 := &pb.ChatMessage{
		From:    1,
		To:      2,
		Message: "hello",
		Time:    nil,
	}

	var respMsgs []*pb.ChatMessage
	for _, v := range resp.GetMessages() {
		ts, err := ptypes.Timestamp(v.GetTime())
		assert.NoError(t, err)
		assert.WithinDuration(t, time.Now(), ts, time.Second*20)

		v.Time = nil
		respMsgs = append(respMsgs, v)
	}
	assert.Equal(t, []*pb.ChatMessage{expectedMsg1, expectedMsg2, expectedMsg3}, resp.GetMessages())
}

func TestServiceTLSInvalidPath(t *testing.T) {
	cfg := config.Config{
		Cassandra: config.CassandraConfig{"127.0.0.1", 9042, 10},
		Port:      14000,
		Metrics:   config.MetricsConfig{14001, "/metrics"},
		TLS:       config.TLSConfig{"crt.pem", "key.pem"},
		NSQ:       config.NSQConfig{"127.0.0.1:4150", "notification_test"},
	}

	_, err := NewService(cfg)
	assert.Error(t, err)
}

func TestServiceTLS(t *testing.T) {
	cfg := config.Config{
		Cassandra: config.CassandraConfig{"127.0.0.1", 9042, 10},
		Port:      14000,
		Metrics:   config.MetricsConfig{14001, "/metrics"},
		TLS:       config.TLSConfig{"../../../../tests/tls/localhost.crt", "../../../../tests/tls/localhost.key"},
		NSQ:       config.NSQConfig{"127.0.0.1:4150", "notification_test"},
	}

	service, err := NewService(cfg)
	assert.NoError(t, err)
	service.Close()
}

type testNSQMessageHandler struct {
	ch chan Message
	t  testing.TB
}

func (h *testNSQMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	var msg Message
	err := json.Unmarshal(m.Body, &msg)
	assert.NoError(h.t, err)
	h.ch <- msg

	return err
}
