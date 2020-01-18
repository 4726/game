package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newMongoQueueTest(t testing.TB) *MongoQueue {
	q, err := NewMongoQueue(100)
	assert.NoError(t, err)
	collection := q.db.Database(q.dbName).Collection(q.dbCollection)
	assert.NoError(t, collection.Drop(context.Background()))
	return q
}

func TestMongoQueueEnqueueFull(t *testing.T) {

}

func TestMongoQueueEnqueueAlreadyInQueue(t *testing.T) {
	q := newMongoQueueTest(t)
	q.Enqueue(1, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	assert.Equal(t, ErrAlreadyInQueue, q.Enqueue(1, 1000))
	assert.Empty(t, ch)
}

func TestMongoQueueEnqueue(t *testing.T) {
	q := newMongoQueueTest(t)

	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	assert.NoError(t, q.Enqueue(1, 1000))
	msg := <-ch
	expectedQD := QueueData{1, 1000, time.Now(), false, 0}
	expectedMsg := PubSubMessage{PubSubTopicAdd, expectedQD}
	assertPubSubMessageEqual(t, expectedMsg, msg)

	assert.Empty(t, ch)
}

func TestMongoQueueDequeue(t *testing.T) {
	q := newMongoQueueTest(t)
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

func TestMongoQueueDeleteOneDoesNotExist(t *testing.T) {
	q := newMongoQueueTest(t)
	q.Enqueue(1, 1000)
	q.Enqueue(2, 1000)
	q.Enqueue(3, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)
	assert.Equal(t, ErrDoesNotExist, q.DeleteOne(4))
	assert.Empty(t, ch)
}

func TestMongoQueueDeleteOne(t *testing.T) {
	q := newMongoQueueTest(t)
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

func TestMongoQueueLen(t *testing.T) {
	q := newMongoQueueTest(t)
	q.Enqueue(1, 1000)
	q.Enqueue(2, 1000)
	q.Enqueue(3, 1000)
	assert.Equal(t, 3, q.Len())
}

func TestMongoQueueMarkMatchFoundDoesNotExist(t *testing.T) {
	q := newMongoQueueTest(t)
	q.Enqueue(1, 1000)
	q.Enqueue(2, 1000)
	q.Enqueue(3, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)
	assert.Equal(t, ErrDoesNotExist, q.MarkMatchFound(4, true))
	assert.Empty(t, ch)
}

func TestMongoQueueMarkMatchFoundNoChange(t *testing.T) {
	q := newMongoQueueTest(t)
	q.Enqueue(1, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	assert.NoError(t, q.MarkMatchFound(1, true))
	<-ch
	assert.NoError(t, q.MarkMatchFound(1, true))
	assert.Empty(t, ch)
}

func TestMongoQueueMarkMatchFound(t *testing.T) {
	q := newMongoQueueTest(t)
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

func TestMongoQueueMarkMatchNotFoundDoesNotExist(t *testing.T) {
	q := newMongoQueueTest(t)
	q.Enqueue(1, 1000)
	q.Enqueue(2, 1000)
	q.Enqueue(3, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)
	assert.Equal(t, ErrDoesNotExist, q.MarkMatchFound(4, false))
	assert.Empty(t, ch)
}

func TestMongoQueueMarkMatchNotFoundNoChange(t *testing.T) {
	q := newMongoQueueTest(t)
	q.Enqueue(1, 1000)
	ch := make(chan PubSubMessage, 1)
	q.Subscribe(ch)

	assert.NoError(t, q.MarkMatchFound(1, false))
	assert.Empty(t, ch)
}

func TestMongoQueueMarkMatchNotFound(t *testing.T) {
	q := newMongoQueueTest(t)
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

func TestMongoQueueEnqueueAndFindMatchAlreadyInQueue(t *testing.T) {
	q := newMongoQueueTest(t)
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

func TestMongoQueueEnqueueAndFindMatchNotEnoughPlayers(t *testing.T) {
	q := newMongoQueueTest(t)
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
func TestMongoQueueEnqueueAndFindMatchNotEnoughPlayers2(t *testing.T) {
	q := newMongoQueueTest(t)
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
func TestMongoQueueEnqueueAndFindMatchMatchFound(t *testing.T) {
	q := newMongoQueueTest(t)
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
