package inmemory

//need to test if all accepts

import (
	"testing"
	"time"

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
	time.Sleep(time.Second)
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
	q := New(1000, 10, 100)
	usersBefore, _ := q.All()
	_, err := q.Accept(1, 1)
	assert.Equal(t, ErrUserNotInMatch, err)
	usersAfter, _ := q.All()
	assert.Equal(t, usersBefore, usersAfter)
}

//group exists but this user is not in the group
func TestAcceptNotInMatch2(t *testing.T) {
	q := New(1000, 10, 100)
	for i := 1; i < 11; i++ {
		ch, _ := q.Join(uint64(i), 1000)
		<-ch
		go func() { <-ch }()
	}
	time.Sleep(time.Second * 2)

	usersBefore, _ := q.All()
	_, err := q.Accept(11, 1)
	assert.Equal(t, ErrUserNotInMatch, err)
	usersAfter, _ := q.All()
	assert.Equal(t, usersBefore, usersAfter)
}

func TestAcceptAlreadyAccepted(t *testing.T) {
	q := New(1000, 10, 100)
	for i := 1; i < 11; i++ {
		ch, _ := q.Join(uint64(i), 1000)
		<-ch
		go func() { <-ch }()
	}
	time.Sleep(time.Second * 2)
	ch, _ := q.Accept(1, 1)
	<-ch

	usersBefore, _ := q.All()
	_, err := q.Accept(1, 1)
	assert.Equal(t, ErrUserAlreadyAccepted, err)
	usersAfter, _ := q.All()
	assert.Equal(t, usersBefore, usersAfter)
	time.Sleep(time.Second)
	assert.Empty(t, ch)
}

func TestAccept(t *testing.T) {
	q := New(1000, 10, 100)
	for i := 1; i < 11; i++ {
		ch, _ := q.Join(uint64(i), 1000)
		<-ch
		go func() { <-ch }()
	}
	usersBefore, _ := q.All()
	time.Sleep(time.Second * 2)
	ch, err := q.Accept(1, 1)
	assert.NoError(t, err)
	expectedMsg := queue.AcceptStatus{
		State: queue.AcceptStateUpdate,
		Data: queue.AcceptStatusUpdateData{
			UsersAccepted: 1,
			UsersNeeded:   10,
		},
	}
	assert.Equal(t, expectedMsg, <-ch)

	usersAfter, _ := q.All()
	expectedData := usersBefore[1].Data.(queue.QueueStateInGroupData)
	expectedData.Accepted = true
	assert.Equal(t, usersBefore[1].Rating, usersAfter[1].Rating)
	assert.Equal(t, usersBefore[1].State, usersAfter[1].State)
	assert.Equal(t, expectedData, usersAfter[1].Data)
	assert.NotNil(t, usersAfter[1].JoinStatusChannel)
	assert.NotNil(t, usersAfter[1].AcceptStatusChannel)

	delete(usersAfter, 1)
	delete(usersBefore, 1)
	assert.Equal(t, usersBefore, usersAfter)
	time.Sleep(time.Second)
	assert.Empty(t, ch)
}

func TestDeclineNotInMatch(t *testing.T) {
	q := New(1000, 10, 100)
	usersBefore, _ := q.All()
	err := q.Decline(1, 1)
	assert.Equal(t, ErrUserNotInMatch, err)
	usersAfter, _ := q.All()
	assert.Equal(t, usersBefore, usersAfter)
}

func TestDeclineNotInMatch2(t *testing.T) {
	q := New(1000, 10, 100)
	for i := 1; i < 11; i++ {
		ch, _ := q.Join(uint64(i), 1000)
		<-ch
		go func() { <-ch }()
	}
	time.Sleep(time.Second * 2)

	usersBefore, _ := q.All()
	err := q.Decline(11, 1)
	assert.Equal(t, ErrUserNotInMatch, err)
	usersAfter, _ := q.All()
	assert.Equal(t, usersBefore, usersAfter)
}

func TestDeclineAlreadyAccepted(t *testing.T) {
	q := New(1000, 10, 100)
	for i := 1; i < 11; i++ {
		ch, _ := q.Join(uint64(i), 1000)
		<-ch
		go func() { <-ch }()
	}
	time.Sleep(time.Second * 2)
	ch, _ := q.Accept(1, 1)
	<-ch

	usersBefore, _ := q.All()
	err := q.Decline(1, 1)
	assert.Equal(t, ErrUserAlreadyAccepted, err)
	usersAfter, _ := q.All()
	assert.Equal(t, usersBefore, usersAfter)
	time.Sleep(time.Second)
	assert.Empty(t, ch)
}

func TestDecline(t *testing.T) {
	q := New(1000, 10, 100)
	for i := 1; i < 11; i++ {
		ch, _ := q.Join(uint64(i), 1000)
		<-ch
		go func() { <-ch }()
	}
	usersBefore, _ := q.All()
	time.Sleep(time.Second * 2)
	err := q.Decline(1, 1)
	assert.NoError(t, err)

	usersAfter, _ := q.All()
	var i uint64
	for i = 2; i < 11; i++ {
		assert.Equal(t, usersBefore[i].Rating, usersAfter[i].Rating)
		assert.Equal(t, queue.QueueStateInQueue, usersAfter[i].State)
		assert.Equal(t, queue.QueueStateInQueueData{}, usersAfter[i].Data)
		assert.NotNil(t, usersAfter[i].JoinStatusChannel)
		assert.Nil(t, usersAfter[i].AcceptStatusChannel)
	}
	assert.Len(t, usersAfter, len(usersBefore)-1)
}
