package main

//todo: add multiqueue support(ex. 2 player queue/5 player queue/etc)

import (
	"errors"
	"math"
	"sync"
	"time"
)

type PubSubTopic int

const (
	PubSubTopicAdd PubSubTopic = iota
	PubSubTopicDelete
	PubSubTopicMatchFound
	PubSubTopicMatchNotFound
)

type QueueData struct {
	UserID     uint64
	Rating     uint64
	StartTime  time.Time
	MatchFound bool
	MatchID    uint64
}

type PubSubMessage struct {
	Topic PubSubTopic
	Data  QueueData
}

type Queue struct {
	sync.Mutex
	data        []QueueData
	limit       int
	subscribers []chan PubSubMessage
	matchID     uint64
}

var ErrAlreadyInQueue = errors.New("user already in queue")
var ErrDoesNotExist = errors.New("does not exist")
var ErrNotAllExists = errors.New("not all exists")
var ErrQueueFull = errors.New("queue full")

func NewQueue() *Queue {
	return &Queue{
		data:        []QueueData{},
		limit:       32766,
		subscribers: []chan PubSubMessage{},
		matchID:     uint64(0),
	}
}

//Enqueue adds a user into the queue and notifies all subscribers
func (q *Queue) Enqueue(userID, rating uint64) error {
	q.Lock()
	defer q.Unlock()

	if q.Len() >= q.limit {
		return ErrQueueFull
	}

	for _, v := range q.data {
		if userID == v.UserID {
			return ErrAlreadyInQueue
		}
	}
	qd := QueueData{userID, rating, time.Now(), false, 0}
	q.data = append(q.data, qd)
	q.publish(PubSubTopicAdd, qd)
	return nil
}

//Dequeue removes a user from the queue and notifies all subscribers
func (q *Queue) Dequeue() {
	q.Lock()
	defer q.Unlock()

	deleted := q.data[0]
	q.data = q.data[1:]
	q.publish(PubSubTopicDelete, deleted)
}

//ForEach calls fn on all elements, return false to break early
func (q *Queue) ForEach(fn func(QueueData) bool) {
	q.Lock()
	defer q.Unlock()

	for _, v := range q.data {
		if !fn(v) {
			return
		}
	}
}

//DeleteOne iterates over elements to find element with userID,
//then deletes the element and notifies subscribers
func (q *Queue) DeleteOne(userID uint64) error {
	q.Lock()
	defer q.Unlock()

	var found bool
	var qd QueueData
	for i, v := range q.data {
		if userID == v.UserID {
			qd = v
			found = true
			if len(q.data) == 1 {
				q.data = []QueueData{}
			} else {
				q.data = append(q.data[:i], q.data[i+1:]...)
			}
			break
		}
	}
	if !found {
		return ErrDoesNotExist
	}
	q.publish(PubSubTopicDelete, qd)
	return nil
}

//Len returns length of queue
func (q *Queue) Len() int {
	q.Lock()
	defer q.Unlock()

	return len(q.data)
}

func (q *Queue) MarkMatchFound(userID uint64, found bool) error {
	q.Lock()
	defer q.Unlock()

	var qd QueueData
	var exists bool
	var pubSubTopic PubSubTopic
	if found {
		pubSubTopic = PubSubTopicMatchFound
	} else {
		pubSubTopic = PubSubTopicMatchNotFound
	}
	for i, v := range q.data {
		if userID == v.UserID {
			exists = true
			q.data[i] = QueueData{v.UserID, v.Rating, v.StartTime, found, 0}
			qd = q.data[i]
			break
		}
	}

	if !exists {
		return ErrDoesNotExist
	}

	q.publish(pubSubTopic, qd)
	return nil
}

// EnqueueAndFindMatch enqueues user to the queue tries to find a match for user,
// if successful, will return []QueueData of other users in the match,
func (q *Queue) EnqueueAndFindMatch(userID, rating, ratingRange uint64, total int) (found bool, matchID uint64, qds []QueueData, err error) {
	q.Lock()
	defer q.Unlock()

	indexes := []int{}
	for i, v := range q.data {
		if userID == v.UserID {
			err = ErrAlreadyInQueue
			return
		}

		if math.Abs(float64(rating-v.Rating)) <= float64(ratingRange/2) {
			if v.MatchFound {
				//player already found a match, ignore
				continue
			}
			indexes = append(indexes, i)
		}
	}

	if len(indexes) < total {
		qd := QueueData{userID, rating, time.Now(), false, uint64(0)}
		q.data = append(q.data, qd)
		q.publish(PubSubTopicAdd, qd)
		return
	}

	matchID = q.getMatchID()

	qd := QueueData{userID, rating, time.Now(), true, matchID}
	q.data = append(q.data, qd)
	q.publish(PubSubTopicAdd, qd)
	q.publish(PubSubTopicMatchFound, qd)

	found = true
	qds = []QueueData{}
	for _, v := range indexes[:total] {
		qd := q.data[v]
		updatedQD = QueueData{qd.UserID, qd.Rating, qd.StartTime, true, matchID}
		q.data[v] = updatedQD
		q.publish(PubSubTopicMatchFound, q.data[v])
		qds = append(qds, q.data[v])
	}

	return
}

func (q *Queue) Subscribe(ch chan PubSubMessage) {
	q.subscribers = append(q.subscribers, ch)
}

func (q *Queue) publish(topic PubSubTopic, qd QueueData) {
	for _, v := range q.subscribers {
		go func() {
			select {
			case v <- PubSubMessage{topic, qd}:
			case <-time.After(time.Second * 10):
			}
		}()
	}
}

func (q *Queue) getMatchID() uint64 {
	q.matchID++
	return q.matchID
}
