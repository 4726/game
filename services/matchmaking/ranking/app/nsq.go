package app

import (
	"bytes"
	"context"
	"time"
	"encoding/json"

	"github.com/4726/game/services/matchmaking/ranking/pb"
	"github.com/nsqio/go-nsq"
	"github.com/go-redis/redis/v7"
)

// nsqMessageHandler implements nsq.Handler
type nsqMessageHandler struct {
	db                   *redis.Client
	zSetName string
}

// Rating is the object passed over nsq
type Rating struct {
	UserID uint64
	Rating uint64
}

// ToRedisZ returns a *redis.Z from Rating
func (r Rating) ToRedisZ() *redis.Z {
	return &redis.Z{float64(r.Rating), rating.UserId}
}

func (h *nsqMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	buffer := bytes.NewBuffer(m.Body)
	var rating Rating
	if err := jsonpb.Unmarshal(buffer, &rating); err != nil {
		return err
	}

	return h.db.ZAdd(h.zSetName, rating.ToRedisZ()).Err()
}
