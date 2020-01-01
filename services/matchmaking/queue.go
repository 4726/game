package main

type QueueData struct {
	UserID uint64
 	Rating uint64
	FoundCh chan uint64
	StartTime time.Time
}

type Queue struct {
	sync.Mutex
	data []QueueData
}

var ErrAlreadyInQueue = errors.New("user already in queue")
var ErrDoesNotExist = errors.New("does not exist")
var ErrNotAllExists = errors.New("not all exists")

func (q *Queue) Enqueue(d QueueData) error {
	q.Lock()
	defer q.Unlock()

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

func (q *Queue) ForEach(fn func(QueueData) bool) {
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
				q.data = []QueueData{}
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
			q.data = []QueueData{}
		} else {
			q.data = append(q.data[:v], q.data[v+1:]...)
		}
	}
}

func (Q *Queue) Clear() {
	q.Lock()
	defer q.Unlock()

	q.data = []QueueData{}
}

func (q *Queue) FindAndDeleteWithinRatingRangeOrEnqueue(d QueueData, ratingRange, total int) ([]QueueData, error) {
	q.Lock()
	defer q.Unlock()

	indexes := []int{}
	for i, v := ragne q.Data {
		if d.UserID == v.UserID {
			return []QueueData, ErrAlreadyInQueue
		}

		if math.Abs(float64(d.Rating - v.Rating)) <= float64(ratingRange) {
			indexes = append(indexes, i)
		}
	}

	if len(indexes) < total {
		q.data = append(q.data, d)
		return []QueueData, nil
	}

	users := []QueueData{}
	for _, v := range indexes {
		users = append(users, q.data[v])
	}

	for _, v := range indexes {
		if len(q.data) == 1 {
			q.data = []QueueData{}
		} else {
			q.data = append(q.data[:v], q.data[v+1:]...)
		}
	}

	return users, nil
}