package main

import (
	"time"
)

type QueueData struct {
	UserID uint64
 	Rating uint64
	FoundCh chan QueueStatus
	StartTime time.Time
	MatchFound bool
}

type QueueStatus struct {
	MatchID uint64
	MatchStarted bool
}

type Queue struct {
	sync.Mutex
	data []*QueueData
	limit int
}

var ErrAlreadyInQueue = errors.New("user already in queue")
var ErrDoesNotExist = errors.New("does not exist")
var ErrNotAllExists = errors.New("not all exists")
var ErrQueueFull = errors.New("queue full")

func NewQueue() *Queue {
	return &Queue{
		data: []*QueueData{},
		limit: 32766,
	}
}

func (q *Queue) Enqueue(d *QueueData) error {
	q.Lock()
	defer q.Unlock()

	if q.Len() >= q.limit {
		return ErrQueueFull
	}

	for _, v := range q.data {
		if d.UserID == v.UserID {
			return ErrAlreadyInQueue
		}
	}
	q.data = append(q.data, d)
	return nil
}

func (q *Queue) Dequeue() {
	q.Lock()
	defer q.Unlock()

	q.data = q.data[1:]
}

func (q *Queue) ForEach(fn func(*QueueData) bool) {
	q.Lock()
	defer q.Unlock()

	for _, v := range q.data {
		if !fn(v) {
			return
		}
	}
}

func (q* Queue) DeleteOne(userID uint64) error {
	q.Lock()
	defer q.Unlock()

	var found bool
	for i, v := range q.data {
		if d.UserID == v.UserID {
			if len(q.data) == 1 {
				q.data = []*QueueData{}
			} else {
				q.data = append(q.data[:i], q.data[i+1:]...)
			}
			return nil
		}
	}
	if !found {
		return ErrDoesNotExist
	}
	return nil
}

func (q *Queue) DeleteMultiple(userIDs []uint64, force bool) error {
	q.Lock()
	defer q.Unlock()

	indexes := []int{}
	for i, v := range q.data {
		for _, v2 := range userIDs {
			if v2.UserID == v.UserID {
				index = append(index, i)
			}
		}
	}

	if !force && len(userIDs) != len(indexes) {
		return ErrNotAllExists
	}

	for _, v := range indexes {
		if len(q.data) == 1 {
			q.data = []*QueueData{}
		} else {
			q.data = append(q.data[:v], q.data[v+1:]...)
		}
	}
}

func (q *Queue) Clear() {
	q.Lock()
	defer q.Unlock()

	q.data = []QueueData{}
}

func (q *Queue) Len() {
	q.Lock()
	defer q.Unlock()

	return len(q.data)
}

func (q *Queue) MarkMatchFound(userID uint64, found bool) {
	q.Lock()
	defer q.Unlock()

	for _, v := range q.data {
		if userID == v.UserID {
			v.MatchFound = found
		}
	}
}

func (q *Queue) MarkMatchFoundMultiple(userIDs []uint64, found bool) {
	q.Lock()
	defer q.Unlock()

	for i, v := range q.data {
		for _, v2 := range userIDs {
			if v2.UserID == v.UserID {
				v.MatchFound = found
			}
		}
	}
}

func (q *Queue) SetMatchStartedAndDelete(userID, matchID uint64, started bool) {
	q.Lock()
	defer q.Unlock()

	indexes := []int{}
	for i, v := range q.data {
		if userID == v.UserID {
			v.FoundCh <- QueueStatus{matchID, started}

			if len(q.data) == 1 {
				q.data = []*QueueData{}
			} else {
				q.data = append(q.data[:i], q.data[i+1:]...)
			}
			return
		}
	}
}

func (q *Queue) FindAndMarkMatchFoundWithinRatingRangeOrEnqueue(d *QueueData, ratingRange uint64, total int) ([]*QueueData, error) {
	q.Lock()
	defer q.Unlock()

	indexes := []int{}
	for i, v := range q.data {
		if d.UserID == v.UserID {
			return []*QueueData{}, ErrAlreadyInQueue
		}

		if math.Abs(float64(d.Rating - v.Rating)) <= float64(ratingRange / 2) {
			if v.MatchFound { 
				//player already found a match, ignore
				continue
			}
			indexes = append(indexes, i)
		}
	}

	if len(indexes) < total {
		q.data = append(q.data, d)
		return []*QueueData{}, nil
	}

	users := []*QueueData{}
	for _, v := range indexes {
		q.data[v].MatchFound = true
		users = append(users, q.data[v])
	}

	return users, nil
}