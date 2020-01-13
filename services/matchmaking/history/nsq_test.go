package main

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/4726/game/services/matchmaking/history/pb"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestHandleMessage(t *testing.T) {
	opts := options.Client().ApplyURI("mongodb://localhost:27017")
	db, err := mongo.Connect(context.Background(), opts)
	assert.NoError(t, err)
	assert.NoError(t, db.Ping(context.Background(), nil))
	collection := db.Database("history_test").Collection("collection_test")
	assert.NoError(t, collection.Drop(context.Background()))

	consumer, err := nsq.NewConsumer("matches_test", "db_test", nsq.NewConfig())
	assert.NoError(t, err)
	defer consumer.Stop()
	consumer.AddHandler(&nsqMessageHandler{db, "history_test", "collection_test"})
	assert.NoError(t, consumer.ConnectToNSQD("127.0.0.1:4150"))

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

	time.Sleep(time.Second * 5)

	matches := []*pb.MatchHistoryInfo{}
	cur, err := collection.Find(context.Background(), bson.D{{}}, options.Find())
	assert.NoError(t, err)
	defer cur.Close(context.Background())
	assert.NoError(t, cur.All(context.Background(), &matches))
	fmt.Println("testMatches: ", len(testMatches))
	fmt.Println("matches: ", len(matches))
	assert.ElementsMatch(t, testMatches, matches)
}

func matchHistoryInfoTestData() []*pb.MatchHistoryInfo {
	winner := &pb.TeamHistoryInfo{
		Users: []uint64{1, 2, 3, 4, 5},
		Score: 20,
	}
	loser := &pb.TeamHistoryInfo{
		Users: []uint64{6, 7, 8, 9, 10},
		Score: 10,
	}
	m1 := &pb.MatchHistoryInfo{
		Id:     1,
		Winner: winner,
		Loser:  loser,
		Time:   ptypes.TimestampNow(),
	}

	return []*pb.MatchHistoryInfo{m1}
}
