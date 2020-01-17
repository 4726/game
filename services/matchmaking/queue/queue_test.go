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
	q := NewQueue(100)
	q.Enqueue(1, 1000)
	q.Enqueue(2, 1000)
	q.Enqueue(3, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)
	assert.Equal(t, ErrDoesNotExist, q.DeleteOne(4))
	assert.Empty(t, ch)
}

func TestQueueDeleteOne(t *testing.T) {
	q := NewQueue(100)
	q.Enqueue(1, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	assert.NoError(t, q.DeleteOne(1))
	msg := <-ch
	expectedQD := QueueData{1, 1000, time.Now(), false, 0}
	expectedMsg := PubSubMessage{PubSubTopicDelete, expectedQD}
	assertPubSubMessageEqual(t, expectedMsg, msg)

	assert.Empty(t, ch)
}

func TestQueueLen(t *testing.T) {
	q := NewQueue(100)
	q.Enqueue(1, 1000)
	q.Enqueue(2, 1000)
	q.Enqueue(3, 1000)
	assert.Equal(t, 3, q.Len())
}

func TestQueueMarkMatchFoundDoesNotExist(t *testing.T) {
	q := NewQueue(100)
	q.Enqueue(1, 1000)
	q.Enqueue(2, 1000)
	q.Enqueue(3, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)
	assert.Equal(t, ErrDoesNotExist, q.MarkMatchFound(4, true))
	assert.Empty(t, ch)
}

func TestQueueMarkMatchFoundNoChange(t *testing.T) {
	q := NewQueue(100)
	q.Enqueue(1, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	assert.NoError(t, q.MarkMatchFound(1, true))
	<-ch
	assert.NoError(t, q.MarkMatchFound(1, true))
	assert.Empty(t, ch)
}

func TestQueueMarkMatchFound(t *testing.T) {
	q := NewQueue(100)
	q.Enqueue(1, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	assert.NoError(t, q.MarkMatchFound(1, true))
	msg := <-ch
	expectedQD := QueueData{1, 1000, time.Now(), true, 0}
	expectedMsg := PubSubMessage{PubSubTopicMatchFound, expectedQD}
	assertPubSubMessageEqual(t, expectedMsg, msg)
	assert.Empty(t, ch)
}

func TestQueueMarkMatchNotFoundDoesNotExist(t *testing.T) {
	q := NewQueue(100)
	q.Enqueue(1, 1000)
	q.Enqueue(2, 1000)
	q.Enqueue(3, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)
	assert.Equal(t, ErrDoesNotExist, q.MarkMatchFound(4, false))
	assert.Empty(t, ch)
}

func TestQueueMarkMatchNotFoundNoChange(t *testing.T) {
	q := NewQueue(100)
	q.Enqueue(1, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	assert.NoError(t, q.MarkMatchFound(1, false))
	assert.Empty(t, ch)
}

func TestQueueMarkMatchNotFound(t *testing.T) {
	q := NewQueue(100)
	q.Enqueue(1, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	assert.NoError(t, q.MarkMatchFound(1, true))
	<-ch

	assert.NoError(t, q.MarkMatchFound(1, false))
	msg := <-ch
	expectedQD := QueueData{1, 1000, time.Now(), false, 0}
	expectedMsg := PubSubMessage{PubSubTopicMatchNotFound, expectedQD}
	assertPubSubMessageEqual(t, expectedMsg, msg)

	assert.Empty(t, ch)
}

func TestQueueEnqueueAndFindMatchAlreadyInQueue(t *testing.T) {
	q := NewQueue(100)
	q.Enqueue(1, 1000)
	q.Enqueue(2, 1000)
	q.Enqueue(3, 1000)
	q.Enqueue(4, 1000)
	q.MarkMatchFound(4, true)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)
	_, _, _, err := q.EnqueueAndFindMatch(4, 1000, 100, 5)
	assert.Equal(t, ErrAlreadyInQueue, err)
	assert.Empty(t, ch)
}

func TestQueueEnqueueAndFindMatchNotEnoughPlayers(t *testing.T) {
	q := NewQueue(100)
	q.Enqueue(1, 1000)
	q.Enqueue(2, 1000)
	q.Enqueue(3, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)
	found, matchID, qds, err := q.EnqueueAndFindMatch(4, 1000, 100, 5)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), matchID)
	assert.False(t, found)
	assert.Empty(t, qds)
	msg := <-ch
	expectedQD := QueueData{4, 1000, time.Now(), false, 0}
	expectedMsg := PubSubMessage{PubSubTopicAdd, expectedQD}
	assertPubSubMessageEqual(t, expectedMsg, msg)
	assert.Empty(t, ch)
}

//tests if matchfound == true will not count as available for match
func TestQueueEnqueueAndFindMatchNotEnoughPlayers2(t *testing.T) {
	q := NewQueue(100)
	q.Enqueue(1, 1000)
	q.Enqueue(2, 1000)
	q.Enqueue(3, 1000)
	q.MarkMatchFound(3, true)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)
	found, matchID, qds, err := q.EnqueueAndFindMatch(4, 1000, 100, 4)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), matchID)
	assert.False(t, found)
	assert.Empty(t, qds)
	msg := <-ch
	expectedQD := QueueData{4, 1000, time.Now(), false, 0}
	expectedMsg := PubSubMessage{PubSubTopicAdd, expectedQD}
	assertPubSubMessageEqual(t, expectedMsg, msg)
	assert.Empty(t, ch)
}

//have to use assertPubSubMessageContains() because channel messages
//come in different orders
func TestQueueEnqueueAndFindMatchMatchFound(t *testing.T) {
	q := NewQueue(100)
	q.Enqueue(1, 950)
	q.Enqueue(2, 1050)
	q.Enqueue(3, 1030)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	found, matchID, qds, err := q.EnqueueAndFindMatch(4, 1000, 100, 4)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, matchID)
	assert.True(t, found)
	assert.Len(t, qds, 4)

	msgs := []PubSubMessage{}
	msgs = append(msgs, <-ch)
	msgs = append(msgs, <-ch)
	msgs = append(msgs, <-ch)
	msgs = append(msgs, <-ch)
	msgs = append(msgs, <-ch)

	expectedQD := QueueData{4, 1000, time.Now(), false, 0}
	expectedMsg := PubSubMessage{PubSubTopicAdd, expectedQD}
	assertPubSubMessageContains(t, expectedMsg, msgs)

	expectedQD = QueueData{1, 950, time.Now(), true, 1}
	expectedMsg = PubSubMessage{PubSubTopicMatchFound, expectedQD}
	assertPubSubMessageContains(t, expectedMsg, msgs)

	expectedQD = QueueData{2, 1050, time.Now(), true, 1}
	expectedMsg = PubSubMessage{PubSubTopicMatchFound, expectedQD}
	assertPubSubMessageContains(t, expectedMsg, msgs)

	expectedQD = QueueData{3, 1030, time.Now(), true, 1}
	expectedMsg = PubSubMessage{PubSubTopicMatchFound, expectedQD}
	assertPubSubMessageContains(t, expectedMsg, msgs)

	expectedQD = QueueData{4, 1000, time.Now(), true, 1}
	expectedMsg = PubSubMessage{PubSubTopicMatchFound, expectedQD}
	assertPubSubMessageContains(t, expectedMsg, msgs)
	assert.Empty(t, ch)
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

func assertPubSubMessageContains(t testing.TB, expected PubSubMessage, list []PubSubMessage) {
	expected.Data.StartTime = time.Time{}
	updatedList := []PubSubMessage{}
	for _, v := range list {
		v.Data.StartTime = time.Time{}
		updatedList = append(updatedList, v)
	}
	assert.Contains(t, updatedList, expected)
}
