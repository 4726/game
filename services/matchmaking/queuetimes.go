package main

import "time"

type QueueTimes struct {
}

func (qt *QueueTimes) Add(rating uint64, duration time.Duration) {

}

func (qt *QueueTimes) EstimatedWaitTime(rating uint64) time.Duration {
	return time.Minute
}
