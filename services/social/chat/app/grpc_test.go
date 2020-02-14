package app

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/4726/game/services/social/chat/config"
	"github.com/4726/game/services/social/chat/pb"
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
	}
	service, err := NewService(cfg)
	fmt.Println(err)
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
}

func TestServiceGetNone(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.GetChatRequest{
		User1:    1,
		User2:      2,
		Total: 100,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.GetChatResponse{
		Messages: []*pb.ChatMessage{},
	}
	assert.Equal(t, expectedResp, resp)
}

func TestServiceGetTotal(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	te.addData(t, 1, 2, "hello")
	te.addData(t, 2, 1, "hi")
	te.addData(t, 1, 2, "bye bye")
	te.addData(t, 2, 1, "bye")

	in := &pb.GetChatRequest{
		User1:    1,
		User2:      2,
		Total: 2,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	expectedMsg1 := &pb.ChatMessage{
		From: 2,
		To: 1,
		Message: "bye",
		Time: time.Time{},
	}
	
	expectedMsg2 := &pb.ChatMessage{
		From: 1,
		To: 2,
		Message: "bye bye",
		Time: time.Time{},
	}
	
	var respMsgs []*pb.ChatMessage 
	for _, v := range resp.GetMessages() {
		t, err := ptypes.Timestamp(v.GetTime())
		assert.NoError(t, err)
		assert.WithinDuration(t, time.Now(), t, time.Second * 20)

		v.Time == time.Time{}
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
		User1:    2,
		User2:      1,
		Total: 100,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	expectedMsg1 := &pb.ChatMessage{
		From: 1,
		To: 2,
		Message: "bye bye",
		Time: time.Time{},
	}
	
	expectedMsg2 := &pb.ChatMessage{
		From: 2,
		To: 1,
		Message: "hi",
		Time: time.Time{},
	}

	expectedMsg2 := &pb.ChatMessage{
		From: 1,
		To: 2,
		Message: "hello",
		Time: time.Time{},
	}
	
	var respMsgs []*pb.ChatMessage 
	for _, v := range resp.GetMessages() {
		t, err := ptypes.Timestamp(v.GetTime())
		assert.NoError(t, err)
		assert.WithinDuration(t, time.Now(), t, time.Second * 20)

		v.Time == time.Time{}
		respMsgs = append(respMsgs, v)
	}
	assert.Equal(t, []*pb.ChatMessage{expectedMsg1, expectedMsg2, expectedMsg3}, resp.GetMessages())
}

func TestServiceGet(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	te.addData(t, 1, 2, "hello")
	te.addData(t, 2, 1, "hi")
	te.addData(t, 1, 2, "bye bye")

	in := &pb.GetChatRequest{
		User1:    1,
		User2:      2,
		Total: 100,
	}
	resp, err := te.c.Get(context.Background(), in)
	assert.NoError(t, err)
	expectedMsg1 := &pb.ChatMessage{
		From: 1,
		To: 2,
		Message: "bye bye",
		Time: time.Time{},
	}
	
	expectedMsg2 := &pb.ChatMessage{
		From: 2,
		To: 1,
		Message: "hi",
		Time: time.Time{},
	}

	expectedMsg2 := &pb.ChatMessage{
		From: 1,
		To: 2,
		Message: "hello",
		Time: time.Time{},
	}
	
	var respMsgs []*pb.ChatMessage 
	for _, v := range resp.GetMessages() {
		t, err := ptypes.Timestamp(v.GetTime())
		assert.NoError(t, err)
		assert.WithinDuration(t, time.Now(), t, time.Second * 20)

		v.Time == time.Time{}
		respMsgs = append(respMsgs, v)
	}
	assert.Equal(t, []*pb.ChatMessage{expectedMsg1, expectedMsg2, expectedMsg3}, resp.GetMessages())
}


