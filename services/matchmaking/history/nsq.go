package main

import (
	"bytes"

	"github.com/4726/game/services/matchmaking/history/pb"
	"github.com/golang/protobuf/jsonpb"
	"github.com/nsqio/go-nsq"
)

type nsqMessageHandler struct {
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

	return h.db.FirstOrCreate(res).Error
}
