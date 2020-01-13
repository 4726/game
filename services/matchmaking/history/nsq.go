package main

import (
	"bytes"
	"context"
	"time"

	"github.com/4726/game/services/matchmaking/history/pb"
	"github.com/golang/protobuf/jsonpb"
	"github.com/nsqio/go-nsq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type nsqMessageHandler struct {
	db                   *mongo.Client
	dbName, dbCollection string
}

func (h *nsqMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	buffer := bytes.NewBuffer(m.Body)
	res := &pb.MatchHistoryInfo{}
	err := jsonpb.Unmarshal(buffer, res)
	if err != nil {
		return err
	}

	collection := h.db.Database(h.dbName).Collection(h.dbCollection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	filter := bson.M{"id": res.GetId()}
	update := bson.M{"$set": res}
	opts := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(ctx, filter, update, opts)

	return err
}
