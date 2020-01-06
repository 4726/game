package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueueEnqueueFull(t *testing.T) {
	q := NewQueue(1)
	q.Enqueue(1, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	assert.Equal(t, ErrQueueFull, q.Enqueue(2, 1000))
	assert.Empty(t, ch)
}

func TestQueueEnqueueAlreadyInQueue(t *testing.T) {
	q := NewQueue(100)
	q.Enqueue(1, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	assert.Equal(t, ErrAlreadyInQueue, q.Enqueue(1, 1000))
	assert.Empty(t, ch)
}

func TestQueueEnqueue(t *testing.T) {
	q := NewQueue(100)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	assert.NoError(t, q.Enqueue(1, 1000))
	msg := <-ch
	expectedQD := QueueData{1, 1000, time.Now(), false, 0}
	expectedMsg := PubSubMessage{PubSubTopicAdd, expectedQD}
	assertPubSubMessageEqual(t, expectedMsg, msg)

	assert.Empty(t, ch)
}

func TestQueueDequeue(t *testing.T) {
	q := NewQueue(100)
	q.Enqueue(1, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	q.Dequeue()
	msg := <-ch
	expectedQD := QueueData{1, 1000, time.Now(), false, 0}
	expectedMsg := PubSubMessage{PubSubTopicDelete, expectedQD}
	assertPubSubMessageEqual(t, expectedMsg, msg)

	assert.Empty(t, ch)
}

func TestQueueForEachExitEarly(t *testing.T) {
	q := NewQueue(100)
	q.Enqueue(1, 1000)
	q.Enqueue(2, 1000)
	q.Enqueue(3, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	var i uint64 = 1
	q.ForEach(func(qd QueueData) bool {
		expectedQD := QueueData{i, 1000, time.Now(), false, 0}
		assertQueueDataEqual(t, expectedQD, qd)
		i++
		if i == 3 {
			return false
		}
		return true
	})
	assert.Equal(t, i, uint64(3))
	assert.Empty(t, ch)
}

func TestQueueForEach(t *testing.T) {
	q := NewQueue(100)
	q.Enqueue(1, 1000)
	q.Enqueue(2, 1000)
	q.Enqueue(3, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	var i uint64 = 1
	q.ForEach(func(qd QueueData) bool {
		expectedQD := QueueData{i, 1000, time.Now(), false, 0}
		assertQueueDataEqual(t, expectedQD, qd)
		i++
		return true
	})
	assert.Equal(t, i, uint64(4))
	assert.Empty(t, ch)
}

func TestQueueDeleteOneDoesNotExist(t *testing.T) {

}

func TestQueueDeleteOne(t *testing.T) {

}

func TestQueueLen(t *testing.T) {

}

func TestQueueMarkMatchFoundDoesNotExist(t *testing.T) {

}

func TestQueueMarkMatchFound(t *testing.T) {

}

func TestQueueMarkMatchNotFoundDoesNotExist(t *testing.T) {

}

func TestQueueMarkMatchNotFound(t *testing.T) {

}

func TestQueueEnqueueAndFindMatchAlreadyInQueue(t *testing.T) {

}

func TestQueueEnqueueAndFindMatchCannotFindMatch(t *testing.T) {

}

func TestQueueEnqueueAndFindMatchMatchFound(t *testing.T) {

}

func assertPubSubMessageEqual(t testing.TB, expected, actual PubSubMessage) {
	assert.Equal(t, expected.Topic, actual.Topic)
	assertQueueDataEqual(t, expected.Data, actual.Data)
}

func assertQueueDataEqual(t testing.TB, expected, actual QueueData) {
	assert.WithinDuration(t, expected.StartTime, actual.StartTime, time.Minute)
	expected.StartTime = time.Time{}
	actual.StartTime = time.Time{}
	assert.Equal(t, expected, actual)
}
