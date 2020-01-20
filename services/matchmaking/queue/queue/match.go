package queue

//can use a seperate hash to store if match expired or someone denied

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

type Match struct {
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
var ErrMatchCancelled = errors.New("match is cancelled")

func NewMatch(ch chan MatchPubSubMessage) (*Match, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if err := r.ConfigSet("notify-keyspace-events", "Eghx").Err(); err != nil {
		return nil, err
	}

	m := &Match{r: r, subscribers: []chan MatchPubSubMessage{ch}}

	pubsub := r.PSubscribe("__keyevent@0__:*", "custom:*")
	if _, err := pubsub.Receive(); err != nil {
		return nil, err
	}
	go func(ch <-chan *redis.Message) {
		for msg := range ch {
			switch msg.Channel {
			case "__keyevent@0__:del":
				matchID, err := strconv.ParseUint(msg.Payload, 10, 64)
				if err != nil {
					continue
				}
				state := MatchStatus{Cancelled: true}
				m.sendState(matchID, state)
			case "__keyevent@0__:expired":
				matchID, err := strconv.ParseUint(msg.Payload, 10, 64)
				if err != nil {
					continue
				}
				state := MatchStatus{Expired: true}
				m.sendState(matchID, state)
			case "custom:accepted":
				matchID, err := strconv.ParseUint(msg.Payload, 10, 64)
				if err != nil {
					continue
				}
				state, _ := m.getState(matchID)
				m.sendState(matchID, state)
			}

		}
	}(pubsub.Channel())

	return m, nil
}

func (m *Match) AddUsers(matchID uint64, users []uint64, expire time.Duration) error {
	_, err := m.r.TxPipelined(func(pipe redis.Pipeliner) error {
		for _, v := range users {
			pipe.HSet(strconv.FormatUint(matchID, 10), strconv.FormatUint(v, 10), "false")
		}
		pipe.Expire(strconv.FormatUint(matchID, 10), expire)
		return nil
	})
	return err
}

func (m *Match) Accept(matchID, userID uint64) error {
	res, err := m.r.HGet(strconv.FormatUint(matchID, 10), strconv.FormatUint(userID, 10)).Result()
	if err != nil {
		if err == redis.Nil {
			//should prob remove from queue cause maybe expired
			return ErrUserNotInMatch
		}
		return err
	}
	if res == "true" {
		return ErrUserAlreadyAccepted
	}

	_, err = m.r.TxPipelined(func(pipe redis.Pipeliner) error {
		if err := pipe.HSet(strconv.FormatUint(matchID, 10), strconv.FormatUint(userID, 10), "true").Err(); err != nil {
			return err
		}

		return pipe.Publish("custom:accepted", strconv.FormatUint(matchID, 10)).Err()
	})

	return err
}

func (m *Match) Decline(matchID, userID uint64) error {
	res, err := m.r.HGet(strconv.FormatUint(matchID, 10), strconv.FormatUint(userID, 10)).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrUserNotInMatch
		}
		return err
	}
	if res == "true" {
		return ErrUserAlreadyAccepted
	}

	return m.r.Del(strconv.FormatUint(matchID, 10)).Err()
}

func (m *Match) TimeUntilExpire(matchID uint64) (time.Duration, error) {
	return m.r.TTL(strconv.FormatUint(matchID, 10)).Result()
}

func (m *Match) State(matchID uint64) (MatchStatus, error) {
	return m.getState(matchID)
}

func (m *Match) getState(matchID uint64) (MatchStatus, error) {
	var state MatchStatus
	state.Players = map[uint64]MatchAcceptStatus{}

	res, err := m.r.HGetAll(strconv.FormatUint(matchID, 10)).Result()
	if err != nil {
		return state, err
	}

	for k, v := range res {
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

func (m *Match) sendState(matchID uint64, state MatchStatus) {
	msg := MatchPubSubMessage{matchID, state}
	for _, v := range m.subscribers {
		v <- msg
	}
}
