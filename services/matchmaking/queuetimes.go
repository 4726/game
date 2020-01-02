package main

import (
	"math"
	"sync"
	"time"
)

type QueueDuration struct {
	Rating   uint64
	Duration time.Duration
}

type QueueTimes struct {
	sync.Mutex
	times []QueueDuration
	limit int
}

func (qt *QueueTimes) Add(qd QueueDuration) {
	qt.Lock()
	defer qt.Unlock()

	qt.times = append(qt.times, qd)
	if len(qt.times) > qt.limit {
		qt.times = qt.times[1:]
	}
}

func (qt *QueueTimes) EstimatedWaitTime(rating, ratingRange uint64) time.Duration {
	qt.Lock()
	defer qt.Unlock()

	totalWithinRange := 0
	totalDuration := 0
	for _, v := range qt.times {
		if math.Abs(float64(rating-v.Rating)) <= float64(ratingRange) {
			totalWithinRange++
			totalDuration += v.Duration
		}
	}

	return totalDuration / totalWithinRange
}
