package queue

import (
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
)

type Queue struct {
	r       *redis.Client
	matchID uint64
	ps      *pubSub
	matches *Matches
}

var ErrAlreadyInQueue = errors.New("user already in queue")
var ErrDoesNotExist = errors.New("does not exist")
var ErrNotAllExists = errors.New("not all exists")
var ErrQueueFull = errors.New("queue full")

func New(subscribers []chan interface{}) (*Queue, error) {
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

	q := &Queue{
		r:       r,
		matchID: uint64(0),
		ps:      newPubSub(subscribers),
	}

	matchCh := make(chan MatchPubSubMessage, 1)
	go func(ch chan MatchPubSubMessage) {
		for msg := range ch {
			if msg.State.Cancelled {
				for k, v := range msg.State.Players {
					if v != MatchDeclined {
						q.ps.publish(PubSubTopicGroupFailed{
							UserID:     k,
							Rating:     0,
							Expired:    false,
							UserDenied: true,
						})
						res, err := q.r.HGet("queue:groupfound", strconv.FormatUint(k, 10)).Result()
						if err != nil {
							continue
						}
						rating, err := strconv.ParseUint(res, 10, 64)
						if err != nil {
							continue
						}
						if err := q.r.HDel("queue:groupfound", strconv.FormatUint(k, 10)).Err(); err != nil {
							continue
						}
						q.EnqueueAndFindMatch(k, rating, 100, 10)
					} else {
						q.DeleteOne(k)
					}
				}
				return
			}
			(userID, rating, ratingRange uint64, total int)
			if msg.State.Expired {
				for k, v := range msg.State.Players {
					if v != MatchAccepted {
						q.DeleteOne(k)
					} else {
						q.ps.publish(PubSubTopicGroupFailed{
							UserID:     k,
							Rating:     0,
							Expired:    true,
							UserDenied: false,
						})
						res, err := q.r.HGet("queue:groupfound", strconv.FormatUint(k, 10)).Result()
						if err != nil {
							continue
						}
						rating, err := strconv.ParseUint(res, 10, 64)
						if err != nil {
							continue
						}
						if err := q.r.HDel("queue:groupfound", strconv.FormatUint(k, 10)).Err(); err != nil {
							continue
						}
						q.EnqueueAndFindMatch(k, rating, 100, 10)
					}
				}
				return
			}

			var usersAccepted, usersInGroup []uint64
			for k, v := range msg.State.Players {
				if v == MatchAccepted {
					usersAccepted = append(usersAccepted, k)
				}
				usersInGroup = append(usersInGroup, k)
			}

			for k, v := range msg.State.Players {
				q.ps.publish(PubSubTopicGroupUpdate{
					UserID:        k,
					Rating:        0,
					UsersAccepted: usersAccepted,
					UsersInGroup:  usersInGroup,
				})
			}
		}
	}(matchCh)
	matches, err := NewMatches(matchCh)
	if err != nil {
		return nil, err
	}

	q.matches = matches

	if err := r.ConfigSet("notify-keyspace-events", "Eghx").Err(); err != nil {
		return nil, err
	}

	pubsub := r.PSubscribe("custom:queue:*")
	if _, err := pubsub.Receive(); err != nil {
		return nil, err
	}
	go func(ch <-chan *redis.Message) {
		for msg := range ch {
			switch msg.Channel {
			case "custom:queue:groupfound":
				userID, err := strconv.ParseUint(msg.Payload, 10, 64)
				if err != nil {
					continue
				}
				q.ps.publish(PubSubTopicGroupFound{
					UserID: userID,
					Rating: 0,
					Users:  usersInGroup,
				})
			case "custom:queue:delete":
				userID, err := strconv.ParseUint(msg.Payload, 10, 64)
				if err != nil {
					continue
				}
				q.ps.publish(PubSubTopicLeftQueue{
					UserID: userID,
					Rating: 0,
				})
			case "custom:queue:add":
				userID, err := strconv.ParseUint(msg.Payload, 10, 64)
				if err != nil {
					continue
				}
				q.ps.publish(PubSubTopicJoinQueue{
					UserID: userID,
					Rating: 0,
				})
			}

		}
	}(pubsub.Channel())

	return q, nil
}

func (q *Queue) Enqueue(userID, rating uint64) error {
	rows, err := q.r.ZAddNX("queue", &redis.Z{float64(rating), userID}).Result()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrAlreadyInQueue
	}

	return q.r.Publish("custom:queue:add", userID).Err()
}

func (q *Queue) DeleteOne(userID uint64) error {
	var zScoreRes *redis.FloatCmd
	var zRemRes *redis.IntCmd
	_, err := q.r.TxPipelined(func(pipe redis.Pipeliner) error {
		zScoreRes = pipe.ZScore("queue", strconv.FormatUint(userID))
		zRemRes = pipe.ZRem("queue", userID)
		return nil
	})
	if err != nil {
		return err
	}
	rows, err := zRemRes.Result()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrDoesNotExist
	}
	rating, err := zScoreRes.Result()
	if err != nil {
		return err
	}

	return q.r.Publish("custom:queue:delete", userID).Err()
}

func (q *Queue) All() (map[uint64]uint64, error) {
	m := map[uint64]uint64{}

	res, err := q.r.ZRangeWithScores("queue", 0, -1).Result()
	if err != nil {
		return m, err
	}

	for _, v := range res {
		userID, err := strconv.ParseUint(v.Member.(string), 10, 64)
		if err != nil {
			return m, err
		}
		m[userID] = uint64(v.Score)
	}

	return m, nil
}

func (q *Queue) Len() (int, error) {
	count, err := q.r.ZCard("queue").Result()
	return int(count), err
}

func (q *Queue) EnqueueAndFindMatch(userID, rating, ratingRange uint64, total int) error {
	var enterQueueRes *redis.FloatCmd
	_, err := q.r.TxPipelined(func(pipe redis.Pipeliner) error {
		//need to also check if is in matchfound hash
		enterQueueRes = pipe.ZAddNX("queue", &redis.Z{float64(rating), userID})
		return nil
	})
	if err != nil {
		return err
	}

	rows, err := enterQueueRes.Result()
	if err != nil {
		return err
	}
	if rows == 0 {
		return  ErrAlreadyInQueue
	}

	if err := q.r.Publish("custom:queue:add", userID).Err(); errr != nil {
		return err
	}

	ratingLessThan := rating + ratingRange/2
	ratingGreaterThan := rating - ratingRange/2

	//need transaction start
	res, err := q.r.ZRangeByScoreWithScores("queue", &redis.ZRangeBy{
		Min:    strconv.FormatUint(ratingGreaterThan, 10),
		Max:    strconv.FormatUint(ratingLessThan, 10),
		Offset: 0,
		Count:  int64(total),
	}).Result()
	if err != nil {
		return err
	}

	if len(res) < total {
		return
	}

	matchID := q.getMatchID()
	var usersInGroup []uint64
	for _, v := range res {
		userID, err := strconv.ParseUint(v.Member.(string), 10, 64)
		if err != nil {
			return err
		}
		rating := uint64(v.Score)
		_, err = q.r.ZRem("queue", userID).Result()
		if err != nil {
			return err
		}
		_, err = q.r.HSet("queue:groupfound", v.Member.(string), rating).Result()
		if err != nil {
			return err
		}
		usersInGroup = append(usersInGroup, userID)
		err = q.r.Publish("custom:groupfound", v.Member.(string)).Err()
		if err != nil {
			return err
		}
	}
	//need transaction end

	return q.matches.AddUsers(matchID, usersInGroup, time.Second*20)
}

func (q *Queue) getMatchID() uint64 {
	q.matchID++
	return q.matchID
}
