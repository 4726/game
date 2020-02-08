package app

import (
	"encoding/json"

	"github.com/go-redis/redis/v7"
	"github.com/nsqio/go-nsq"
)

// nsqMessageHandler implements nsq.Handler
type nsqMessageHandler struct {
	db       *redis.Client
	zSetName string
}

// Rating is the object passed over nsq
type Rating struct {
	UserID uint64
	Rating uint64
}

// ToRedisZ returns a *redis.Z from Rating
func (r Rating) ToRedisZ() *redis.Z {
	return &redis.Z{float64(r.Rating), r.UserID}
}

func (h *nsqMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	var rating Rating
	if err := json.Unmarshal(m.Body, &rating); err != nil {
		return err
	}

	return h.db.ZAdd(h.zSetName, rating.ToRedisZ()).Err()
}
