package inmemory

import (
	"testing"

	"github.com/4726/game/services/matchmaking/queue/queue"
	"github.com/stretchr/testify/assert"
)

func TestJoinFull(t *testing.T) {
	q := New(1, 10, 100)
	q.Join(1, 1000)

	usersBefore, _ := q.All()

	_, err := q.Join(2, 1000)
	assert.Equal(t, ErrQueueFull, err)
	usersAfter, _ := q.All()
	assert.Equal(t, usersBefore, usersAfter)
}

func TestJoinAlreadyInQueue(t *testing.T) {
	q := New(1000, 10, 100)
	q.Join(1, 1000)

	usersBefore, _ := q.All()

	_, err := q.Join(1, 1000)
	assert.Equal(t, ErrAlreadyInQueue, err)
	usersAfter, _ := q.All()
	assert.Equal(t, usersBefore, usersAfter)
}

func TestJoin(t *testing.T) {
	q := New(1000, 10, 100)
	ch, err := q.Join(1, 1000)
	assert.NoError(t, err)
	expectedMsg := queue.JoinStatus{
		State: queue.JoinStateEntered,
		Data:  queue.JoinStateEnteredData{},
	}
	assert.Equal(t, expectedMsg, <-ch)
	assert.Empty(t, ch)

	users, _ := q.All()
	assert.Len(t, users, 1)
	assert.Equal(t, uint64(1000), users[1].Rating)
	assert.Equal(t, queue.QueueStateInQueue, users[1].State)
	assert.Equal(t, queue.QueueStateInQueueData{}, users[1].Data)
	assert.NotNil(t, users[1].JoinStatusChannel)
	assert.Nil(t, users[1].AcceptStatusChannel)
}

func TestLeaveDoesNotExist(t *testing.T) {
	q := New(1000, 10, 100)

	usersBefore, _ := q.All()

	err := q.Leave(1)
	assert.Equal(t, ErrDoesNotExist, err)
	usersAfter, _ := q.All()
	assert.Equal(t, usersBefore, usersAfter)
}

func TestLeave(t *testing.T) {
	q := New(1000, 10, 100)
	usersBefore, _ := q.All()
	q.Join(1, 1000)

	err := q.Leave(1)
	assert.NoError(t, err)
	usersAfter, _ := q.All()
	assert.Equal(t, usersBefore, usersAfter)
}

func TestAcceptNotInMatch(t *testing.T) {

}

func TestAcceptAlreadyAccepted(t *testing.T) {

}

func TestAccept(t *testing.T) {

}

func TestDeclineNotInMatch(t *testing.T) {

}

func TestDeclineAlreadyAccepted(t *testing.T) {

}

func TestDecline(t *testing.T) {

}
