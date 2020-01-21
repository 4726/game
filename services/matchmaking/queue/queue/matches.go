package queue

import (
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type MatchAcceptStatus int

const (
	MatchAccepted MatchAcceptStatus = iota
	MatchDeclined
	MatchUnknown
)

type Matches struct {
	r           *redis.Client
	subscribers []chan MatchPubSubMessage
}

type MatchStatus struct {
	TotalAccepted, TotalNeeded int
	Cancelled                  bool //a user declined
	Players                    map[uint64]MatchAcceptStatus
	Expired                    bool
}

type MatchPubSubMessage struct {
	MatchID uint64
	State   MatchStatus
}

var ErrUserNotInMatch = errors.New("user is not in this match")
var ErrUserAlreadyAccepted = errors.New("user already accepted")
var ErrUserAlreadyDeclined = errors.New("user already declined")
var ErrMatchCancelled = errors.New("match is cancelled") //a user declined

func NewMatches(ch chan MatchPubSubMessage) (*Matches, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if err := r.Ping().Err(); err != nil {
		return nil, err
	}

	if err := r.ConfigSet("notify-keyspace-events", "Eghx").Err(); err != nil {
		return nil, err
	}

	m := &Matches{r: r, subscribers: []chan MatchPubSubMessage{ch}}

	pubsub := r.PSubscribe("__keyevent@0__:*", "custom:match:*")
	if _, err := pubsub.Receive(); err != nil {
		return nil, err
	}
	go func(ch <-chan *redis.Message) {
		for msg := range ch {
			switch msg.Channel {
			case "__keyevent@0__:expired":
				matchID, err := strconv.ParseUint(msg.Payload, 10, 64)
				if err != nil {
					continue
				}
				state := MatchStatus{Expired: true}
				m.sendState(matchID, state)
			case "custom:match:accepted":
				matchID, err := strconv.ParseUint(msg.Payload, 10, 64)
				if err != nil {
					continue
				}
				state, _ := m.getState(matchID)
				m.sendState(matchID, state)
			case "custom:match:cancelled":
				matchID, err := strconv.ParseUint(msg.Payload, 10, 64)
				if err != nil {
					continue
				}
				state := MatchStatus{Cancelled: true}
				m.sendState(matchID, state)
			}

		}
	}(pubsub.Channel())

	return m, nil
}

func (m *Matches) AddUsers(matchID uint64, users []uint64, expire time.Duration) error {
	_, err := m.r.TxPipelined(func(pipe redis.Pipeliner) error {
		for _, v := range users {
			pipe.HSet(strconv.FormatUint(matchID, 10), strconv.FormatUint(v, 10), "false")
		}
		pipe.HSet(strconv.FormatUint(matchID, 10), "deleted", "0")
		pipe.Expire(strconv.FormatUint(matchID, 10), expire)
		return nil
	})
	return err
}

func (m *Matches) Accept(matchID, userID uint64) error {
	res, err := m.r.HMGet(strconv.FormatUint(matchID, 10), strconv.FormatUint(userID, 10), "deleted").Result()
	if err != nil {
		return err
	}
	if res[0] == nil {
		return ErrUserNotInMatch
	}
	if res[0] == "true" {
		return ErrUserAlreadyAccepted
	}
	if res[1] == "1" {
		//a user declined
		return ErrMatchCancelled
	}

	_, err = m.r.TxPipelined(func(pipe redis.Pipeliner) error {
		if err := pipe.HSet(strconv.FormatUint(matchID, 10), strconv.FormatUint(userID, 10), "true").Err(); err != nil {
			return err
		}

		return pipe.Publish("custom:match:accepted", strconv.FormatUint(matchID, 10)).Err()
	})

	return err
}

func (m *Matches) Decline(matchID, userID uint64) error {
	res, err := m.r.HMGet(strconv.FormatUint(matchID, 10), strconv.FormatUint(userID, 10), "deleted").Result()
	if err != nil {
		return err
	}
	if res[0] == nil {
		return ErrUserNotInMatch
	}
	if res[0] == "true" {
		return ErrUserAlreadyAccepted
	}
	if res[1] == "1" {
		//a user declined
		return ErrMatchCancelled
	}

	_, err = m.r.TxPipelined(func(pipe redis.Pipeliner) error {
		if err := pipe.HSet(strconv.FormatUint(matchID, 10), "deleted", "1").Err(); err != nil {
			return err
		}

		return pipe.Publish("custom:match:cancelled", strconv.FormatUint(matchID, 10)).Err()
	})

	return err
}

func (m *Matches) TimeUntilExpire(matchID uint64) (time.Duration, error) {
	return m.r.TTL(strconv.FormatUint(matchID, 10)).Result()
}

func (m *Matches) State(matchID uint64) (MatchStatus, error) {
	return m.getState(matchID)
}

func (m *Matches) getState(matchID uint64) (MatchStatus, error) {
	var state MatchStatus
	state.Players = map[uint64]MatchAcceptStatus{}

	res, err := m.r.HGetAll(strconv.FormatUint(matchID, 10)).Result()
	if err != nil {
		return state, err
	}

	for k, v := range res {
		if k == "deleted" {
			if v == "1" {
				state.Cancelled = true
			}
			continue
		}
		userID, err := strconv.ParseUint(k, 10, 64)
		if err != nil {
			return state, err
		}
		if v == "true" {
			state.Players[userID] = MatchAccepted
			state.TotalAccepted++
		} else {
			state.Players[userID] = MatchUnknown
		}
		state.TotalNeeded++
	}

	return state, nil
}

func (m *Matches) sendState(matchID uint64, state MatchStatus) {
	msg := MatchPubSubMessage{matchID, state}
	for _, v := range m.subscribers {
		v <- msg
	}
}
