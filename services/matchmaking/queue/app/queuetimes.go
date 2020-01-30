package app

import (
	"math"
	"sync"
	"time"
)

type queueDuration struct {
	Rating   uint64
	Duration time.Duration
}

type queueTimes struct {
	sync.Mutex
	times []queueDuration
	limit int
}

func newQueueTimes(limit int) *queueTimes {
	if limit < 1 {
		limit = 1
	}
	return &queueTimes{
		times: []queueDuration{},
		limit: limit,
	}
}

func (qt *queueTimes) Add(qd queueDuration) {
	qt.Lock()
	defer qt.Unlock()

	qt.times = append(qt.times, qd)
	if len(qt.times) > qt.limit {
		qt.times = qt.times[1:]
	}
}

func (qt *queueTimes) EstimatedWaitTime(rating, ratingRange uint64) time.Duration {
	qt.Lock()
	defer qt.Unlock()

	totalWithinRange := 0
	totalDuration := time.Duration(0)
	for _, v := range qt.times {
		if math.Abs(float64(rating-v.Rating)) <= float64(ratingRange) {
			totalWithinRange++
			totalDuration += v.Duration
		}
	}

	if totalWithinRange == 0 {
		return time.Duration(0)
	}

	return totalDuration / time.Duration(totalWithinRange)
}
