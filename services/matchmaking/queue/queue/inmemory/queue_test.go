package inmemory

//need to test if all accepts

import (
	"sync"
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
	ch, _ := q.Join(1, 1000)
	time.Sleep(time.Second)
	var msgs []queue.JoinStatus
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			select {
			case <-time.After(time.Second * 3):
				wg.Done()
				return
			case msg, ok := <-ch:
				if !ok {
					wg.Done()
					return
				}
				msgs = append(msgs, msg)
			}
		}

	}()
	err := q.Leave(1)
	assert.NoError(t, err)
	usersAfter, _ := q.All()
	assert.Equal(t, usersBefore, usersAfter)
	wg.Wait()
	assert.Len(t, msgs, 2)
	expectedMsg := queue.JoinStatus{
		State: queue.JoinStateLeft,
		Data:  queue.JoinStateLeftData{},
	}
	assert.Equal(t, expectedMsg, msgs[1])
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

func TestAcceptChannelMessage(t *testing.T) {
	q := New(1000, 10, 100)
	for i := 1; i < 11; i++ {
		ch, _ := q.Join(uint64(i), 1000)
		<-ch
		go func() { <-ch }()
	}
	time.Sleep(time.Second * 2)
	ch, _ := q.Accept(1, 1)
	var msgs []queue.AcceptStatus
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					wg.Done()
					return
				}
				msgs = append(msgs, msg)
			case <-time.After(time.Second * 3):
				wg.Done()
				return
			}
		}
	}()
	time.Sleep(time.Second * 2)

	ch2, _ := q.Accept(2, 1)
	<-ch2

	wg.Wait()
	expectedAcceptedMsg := queue.AcceptStatus{
		State: queue.AcceptStateUpdate,
		Data: queue.AcceptStatusUpdateData{
			UsersAccepted: 2,
			UsersNeeded:   10,
		},
	}
	assert.Contains(t, msgs, expectedAcceptedMsg)
	assert.Len(t, msgs, 2)
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

func TestAcceptAllAccepted(t *testing.T) {
	q := New(1000, 10, 100)
	for i := 1; i < 11; i++ {
		ch, _ := q.Join(uint64(i), 1000)
		go func() {
			for range ch {
			}
		}()
	}
	usersBefore, _ := q.All()
	time.Sleep(time.Second * 2)
	for i := 2; i < 11; i++ {
		ch, err := q.Accept(uint64(i), 1)
		assert.NoError(t, err)
		go func() {
			for range ch {
			}
		}()
	}
	time.Sleep(time.Second * 2)

	ch, err := q.Accept(1, 1)
	assert.NoError(t, err)
	expectedMsg := queue.AcceptStatus{
		State: queue.AcceptStateSuccess,
		Data: queue.AcceptStateSuccessData{
			UserCount: 10,
			MatchID:   1,
		},
	}

	assert.Equal(t, expectedMsg, <-ch)

	foundCh := q.Channel()
	expectedUsers := map[uint64]uint64{}
	for i := 1; i < 11; i++ {
		expectedUsers[uint64(i)] = 1000
	}
	expectedFoundMsg := queue.Match{
		Users:   expectedUsers,
		MatchID: 1,
	}
	foundMsg := <-foundCh
	assert.Equal(t, expectedFoundMsg, foundMsg)

	usersAfter, _ := q.All()

	for i := 1; i < 11; i++ {
		delete(usersBefore, uint64(i))
	}
	assert.Equal(t, usersBefore, usersAfter)
	assert.Empty(t, q.groups)
	time.Sleep(time.Second)
	assert.Empty(t, ch)
	assert.Empty(t, foundCh)
}

func TestAcceptAllAcceptedLater(t *testing.T) {
	q := New(1000, 10, 100)
	for i := 1; i < 11; i++ {
		ch, _ := q.Join(uint64(i), 1000)
		go func() {
			for range ch {
			}
		}()
	}
	time.Sleep(time.Second * 2)
	ch, err := q.Accept(1, 1)
	assert.NoError(t, err)
	var chMsgs []queue.AcceptStatus
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					wg.Done()
					return
				}
				chMsgs = append(chMsgs, msg)
			case <-time.After(time.Second * 5):
				wg.Done()
				return
			}
		}
	}()
	time.Sleep(time.Second * 2)
	for i := 2; i < 11; i++ {
		ch, err := q.Accept(uint64(i), 1)
		assert.NoError(t, err)
		go func() {
			for range ch {
			}
		}()
	}
	time.Sleep(time.Second * 2)
	wg.Wait()
	expectedMsg := queue.AcceptStatus{
		State: queue.AcceptStateSuccess,
		Data: queue.AcceptStateSuccessData{
			UserCount: 10,
			MatchID:   1,
		},
	}
	assert.Contains(t, chMsgs, expectedMsg)
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
	assert.NotContains(t, usersAfter, uint64(1))
	assert.NotContains(t, q.groups, uint64(1))
}
