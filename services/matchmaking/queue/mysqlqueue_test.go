package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newMysqlQueueTest(t testing.TB, limit int) *MysqlQueue {
	q, err := NewMysqlQueue(limit)
	assert.NoError(t, err)
	q.db.Exec("TRUNCATE queue_data;")
	return q
}

func TestMysqlQueueEnqueueFull(t *testing.T) {}

func TestMysqlQueueEnqueueAlreadyInQueue(t *testing.T) {
	q := newMysqlQueueTest(t, 100)
	q.Enqueue(1, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	assert.Equal(t, ErrAlreadyInQueue, q.Enqueue(1, 1000))
	assert.Empty(t, ch)
}

func TestMysqlQueueEnqueue(t *testing.T) {
	q := newMysqlQueueTest(t, 100)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	assert.NoError(t, q.Enqueue(1, 1000))
	msg := <-ch
	expectedQD := QueueData{1, 1000, time.Now(), false, 0}
	expectedMsg := PubSubMessage{PubSubTopicAdd, expectedQD}
	assertPubSubMessageEqual(t, expectedMsg, msg)

	assert.Empty(t, ch)
}

func TestMysqlQueueDeleteOneDoesNotExist(t *testing.T) {
	q := newMysqlQueueTest(t, 100)
	q.Enqueue(1, 1000)
	q.Enqueue(2, 1000)
	q.Enqueue(3, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)
	assert.Equal(t, ErrDoesNotExist, q.DeleteOne(4))
	assert.Empty(t, ch)
}

func TestMysqlQueueDeleteOne(t *testing.T) {
	q := newMysqlQueueTest(t, 100)
	q.Enqueue(1, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	assert.NoError(t, q.DeleteOne(1))
	msg := <-ch
	assert.Equal(t, PubSubTopicDelete, msg.Topic)
	assert.Equal(t, uint64(1), msg.Data.UserID)

	assert.Empty(t, ch)
}

func TestMysqlQueueLen(t *testing.T) {
	q := newMysqlQueueTest(t, 100)
	q.Enqueue(1, 1000)
	q.Enqueue(2, 1000)
	q.Enqueue(3, 1000)
	assert.Equal(t, 3, q.Len())
}

func TestMysqlQueueMarkMatchFoundDoesNotExist(t *testing.T) {
	q := newMysqlQueueTest(t, 100)
	q.Enqueue(1, 1000)
	q.Enqueue(2, 1000)
	q.Enqueue(3, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)
	assert.NoError(t, q.MarkMatchFound(4, true))
	assert.Empty(t, ch)
}

func TestMysqlQueueMarkMatchFoundNoChange(t *testing.T) {
	q := newMysqlQueueTest(t, 100)
	q.Enqueue(1, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	assert.NoError(t, q.MarkMatchFound(1, true))
	<-ch
	assert.NoError(t, q.MarkMatchFound(1, true))
	assert.Empty(t, ch)
}

func TestMysqlQueueMarkMatchFound(t *testing.T) {
	q := newMysqlQueueTest(t, 100)
	q.Enqueue(1, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	assert.NoError(t, q.MarkMatchFound(1, true))
	msg := <-ch
	assert.Equal(t, PubSubTopicMatchFound, msg.Topic)
	assert.Equal(t, uint64(1), msg.Data.UserID)
	assert.Empty(t, ch)
}

func TestMysqlQueueMarkMatchNotFoundDoesNotExist(t *testing.T) {
	q := newMysqlQueueTest(t, 100)
	q.Enqueue(1, 1000)
	q.Enqueue(2, 1000)
	q.Enqueue(3, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)
	assert.NoError(t, q.MarkMatchFound(4, false))
	assert.Empty(t, ch)
}

func TestMysqlQueueMarkMatchNotFoundNoChange(t *testing.T) {
	q := newMysqlQueueTest(t, 100)
	q.Enqueue(1, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	assert.NoError(t, q.MarkMatchFound(1, false))
	assert.Empty(t, ch)
}

func TestMysqlQueueMarkMatchNotFound(t *testing.T) {
	q := newMysqlQueueTest(t, 100)
	q.Enqueue(1, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	assert.NoError(t, q.MarkMatchFound(1, true))
	<-ch

	assert.NoError(t, q.MarkMatchFound(1, false))
	msg := <-ch
	assert.Equal(t, PubSubTopicMatchNotFound, msg.Topic)
	assert.Equal(t, uint64(1), msg.Data.UserID)

	assert.Empty(t, ch)
}

func TestMysqlQueueEnqueueAndFindMatchAlreadyInQueue(t *testing.T) {
	q := newMysqlQueueTest(t, 100)
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

func TestMysqlQueueEnqueueAndFindMatchNotEnoughPlayers(t *testing.T) {
	q := newMysqlQueueTest(t, 100)
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
func TestMysqlQueueEnqueueAndFindMatchNotEnoughPlayers2(t *testing.T) {
	q := newMysqlQueueTest(t, 100)
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
func TestMysqlQueueEnqueueAndFindMatchMatchFound(t *testing.T) {
	q := newMysqlQueueTest(t, 100)
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
