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
	startTime  time.Time
	MatchFound bool
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
	qd := QueueData{userID, rating, time.Now(), false}
	q.data = append(q.data, qd)
	q.Publish(PubSubTopicAdd, qd)
	return nil
}

//Dequeue removes a user from the queue and notifies all subscribers
func (q *Queue) Dequeue() {
	q.Lock()
	defer q.Unlock()

	deleted := q.data[0]
	q.data = q.data[1:]
	q.Publish(PubSubTopicDelete, deleted)
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
	q.Publish(PubSubTopicDelete, qd)
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
	for _, v := range q.data {
		if userID == v.UserID {
			exists = true
			v.MatchFound = found
			qd = v
			break
		}
	}

	if !exists {
		return ErrDoesNotExist
	}

	q.Publish(pubSubTopic, qd)
	return nil
}

// FindMatchOrEnqueue tries to find a match for user,
// if successful, will return []QueueData of other users in the match,
// else will add the user to the queue
func (q *Queue) FindMatchOrEnqueue(userID, rating, ratingRange uint64, total int) ([]QueueData, error) {
	q.Lock()
	defer q.Unlock()

	indexes := []int{}
	for i, v := range q.data {
		if userID == v.UserID {
			return []QueueData{}, ErrAlreadyInQueue
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
		//Enqueue
		qd := QueueData{userID, rating, time.Now(), false}
		q.data = append(q.data, qd)
		q.Publish(PubSubTopicDelete, qd)
		return []QueueData{}, nil
	}

	users := []QueueData{}
	for _, v := range indexes {
		//MarkMatchFound true
		q.data[v].MatchFound = true
		q.Publish(PubSubTopicMatchFound, q.data[v])
		users = append(users, q.data[v])
	}

	return users, nil
}

func (q *Queue) Publish(topic PubSubTopic, qd QueueData) {
	for _, v := range q.subscribers {
		go func() {
			select {
			case v <- PubSubMessage{topic, qd}:
			case <-time.After(time.Second * 10):
			}
		}()
	}
}
