package memory

import (
	"sync"
	"testing"
	"time"

	"github.com/4726/game/services/matchmaking/queue/engine"
	"github.com/stretchr/testify/assert"
)

func teardown(t testing.TB, q *Queue) {
	for _, v := range q.groupTimers {
		v.Stop()
	}
}

func TestJoinFull(t *testing.T) {
	q := New(1, 10, 100, time.Second*20)
	defer teardown(t, q)
	q.Join(1, 1000)

	usersBefore, _ := q.All()

	_, err := q.Join(2, 1000)
	assert.Equal(t, ErrQueueFull, err)
	usersAfter, _ := q.All()
	assert.Equal(t, usersBefore, usersAfter)
}

func TestJoinAlreadyInQueue(t *testing.T) {
	q := New(1000, 10, 100, time.Second*20)
	defer teardown(t, q)
	q.Join(1, 1000)

	usersBefore, _ := q.All()

	_, err := q.Join(1, 1000)
	assert.Equal(t, ErrAlreadyInQueue, err)
	usersAfter, _ := q.All()
	assert.Equal(t, usersBefore, usersAfter)
}

func TestJoin(t *testing.T) {
	q := New(1000, 10, 100, time.Second*20)
	defer teardown(t, q)
	ch, err := q.Join(1, 1000)
	assert.NoError(t, err)
	expectedMsg := engine.JoinStatus{
		State: engine.JoinStateEntered,
		Data:  engine.JoinStateEnteredData{},
	}
	assert.Equal(t, expectedMsg, <-ch)
	time.Sleep(time.Second)
	assert.Empty(t, ch)

	users, _ := q.All()
	assert.Len(t, users, 1)
	assert.Equal(t, uint64(1000), users[1].Rating)
	assert.Equal(t, engine.QueueStateInQueue, users[1].State)
	assert.Equal(t, engine.QueueStateInQueueData{}, users[1].Data)
	assert.NotNil(t, users[1].JoinStatusChannel)
	assert.Nil(t, users[1].AcceptStatusChannel)
}

func TestLeaveDoesNotExist(t *testing.T) {
	q := New(1000, 10, 100, time.Second*20)
	defer teardown(t, q)
	usersBefore, _ := q.All()

	err := q.Leave(1)
	assert.Equal(t, ErrDoesNotExist, err)
	usersAfter, _ := q.All()
	assert.Equal(t, usersBefore, usersAfter)
}

func TestLeave(t *testing.T) {
	q := New(1000, 10, 100, time.Second*20)
	defer teardown(t, q)
	usersBefore, _ := q.All()
	ch, _ := q.Join(1, 1000)
	time.Sleep(time.Second)
	var msgs []engine.JoinStatus
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
	expectedMsg := engine.JoinStatus{
		State: engine.JoinStateLeft,
		Data:  engine.JoinStateLeftData{},
	}
	assert.Equal(t, expectedMsg, msgs[1])
}

func TestAcceptNotInMatch(t *testing.T) {
	q := New(1000, 10, 100, time.Second*20)
	defer teardown(t, q)
	usersBefore, _ := q.All()
	_, err := q.Accept(1, 1)
	assert.Equal(t, ErrUserNotInMatch, err)
	usersAfter, _ := q.All()
	assert.Equal(t, usersBefore, usersAfter)
}

//group exists but this user is not in the group
func TestAcceptNotInMatch2(t *testing.T) {
	q := New(1000, 10, 100, time.Second*20)
	defer teardown(t, q)
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
	q := New(1000, 10, 100, time.Second*20)
	defer teardown(t, q)
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

func TestAcceptDeclinedBefore(t *testing.T) {
	q := New(1000, 10, 100, time.Second*20)
	defer teardown(t, q)
	for i := 1; i < 11; i++ {
		ch, _ := q.Join(uint64(i), 1000)
		go func() {
			for range ch {
			}
		}()
	}
	time.Sleep(time.Second * 2)
	q.Decline(2, 1)
	time.Sleep(time.Second * 2)

	_, err := q.Accept(1, 1)
	assert.Equal(t, ErrUserNotInMatch, err)
}

func TestAcceptDeclinedAfter(t *testing.T) {
	q := New(1000, 10, 100, time.Second*20)
	defer teardown(t, q)
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
	var wg sync.WaitGroup
	wg.Add(1)
	var msgs []engine.AcceptStatus
	go func() {
		for {
			msg, ok := <-ch
			if !ok {
				wg.Done()
				return
			}
			msgs = append(msgs, msg)
		}
	}()
	time.Sleep(time.Second * 2)

	q.Decline(2, 1)
	wg.Wait()
	assert.Len(t, msgs, 2)
	expectedMsg := engine.AcceptStatus{
		State: engine.AcceptStateFailed,
		Data:  engine.AcceptStateFailedData{},
	}
	assert.Equal(t, expectedMsg, msgs[1])

}

func TestAcceptChannelMessage(t *testing.T) {
	q := New(1000, 10, 100, time.Second*20)
	defer teardown(t, q)
	for i := 1; i < 11; i++ {
		ch, _ := q.Join(uint64(i), 1000)
		<-ch
		go func() { <-ch }()
	}
	time.Sleep(time.Second * 2)
	ch, _ := q.Accept(1, 1)
	var msgs []engine.AcceptStatus
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
	expectedAcceptedMsg := engine.AcceptStatus{
		State: engine.AcceptStateUpdate,
		Data: engine.AcceptStatusUpdateData{
			UsersAccepted: 2,
			UsersNeeded:   10,
		},
	}
	assert.Contains(t, msgs, expectedAcceptedMsg)
	assert.Len(t, msgs, 2)
}

func TestAccept(t *testing.T) {
	q := New(1000, 10, 100, time.Second*20)
	defer teardown(t, q)
	for i := 1; i < 11; i++ {
		ch, _ := q.Join(uint64(i), 1000)
		<-ch
		go func() { <-ch }()
	}
	usersBefore, _ := q.All()
	time.Sleep(time.Second * 2)
	ch, err := q.Accept(1, 1)
	assert.NoError(t, err)
	expectedMsg := engine.AcceptStatus{
		State: engine.AcceptStateUpdate,
		Data: engine.AcceptStatusUpdateData{
			UsersAccepted: 1,
			UsersNeeded:   10,
		},
	}
	assert.Equal(t, expectedMsg, <-ch)

	usersAfter, _ := q.All()
	expectedData := usersBefore[1].Data.(engine.QueueStateInGroupData)
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
	q := New(1000, 10, 100, time.Second*20)
	defer teardown(t, q)
	usersBefore, _ := q.All()
	for i := 1; i < 11; i++ {
		ch, _ := q.Join(uint64(i), 1000)
		go func() {
			for range ch {
			}
		}()
	}
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
	expectedMsg := engine.AcceptStatus{
		State: engine.AcceptStateSuccess,
		Data: engine.AcceptStateSuccessData{
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
	expectedFoundMsg := engine.Match{
		Users:   expectedUsers,
		MatchID: 1,
	}
	foundMsg := <-foundCh
	assert.Equal(t, expectedFoundMsg, foundMsg)

	usersAfter, _ := q.All()
	assert.Equal(t, usersBefore, usersAfter)
	assert.Empty(t, q.groups)
	time.Sleep(time.Second)
	assert.Empty(t, ch)
	assert.Empty(t, foundCh)
}

func TestAcceptAllAcceptedLater(t *testing.T) {
	q := New(1000, 10, 100, time.Second*20)
	defer teardown(t, q)
	usersBefore, _ := q.All()
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
	var chMsgs []engine.AcceptStatus
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
	expectedMsg := engine.AcceptStatus{
		State: engine.AcceptStateSuccess,
		Data: engine.AcceptStateSuccessData{
			UserCount: 10,
			MatchID:   1,
		},
	}
	assert.Contains(t, chMsgs, expectedMsg)

	foundCh := q.Channel()
	expectedUsers := map[uint64]uint64{}
	for i := 1; i < 11; i++ {
		expectedUsers[uint64(i)] = 1000
	}
	expectedFoundMsg := engine.Match{
		Users:   expectedUsers,
		MatchID: 1,
	}
	foundMsg := <-foundCh
	assert.Equal(t, expectedFoundMsg, foundMsg)

	usersAfter, _ := q.All()
	assert.Equal(t, usersBefore, usersAfter)
	assert.Empty(t, q.groups)
	time.Sleep(time.Second)
}

func TestDeclineDoesNotExist(t *testing.T) {
	q := New(1000, 10, 100, time.Second*20)
	defer teardown(t, q)
	usersBefore, _ := q.All()
	err := q.Decline(1, 1)
	assert.Equal(t, ErrUserNotInMatch, err)
	usersAfter, _ := q.All()
	assert.Equal(t, usersBefore, usersAfter)
}

func TestDeclineMatchDoesNotExist(t *testing.T) {
	q := New(1000, 10, 100, time.Second*20)
	defer teardown(t, q)
	ch, _ := q.Join(1, 1000)
	go func() {
		for range ch {
		}
	}()
	time.Sleep(time.Second * 2)

	usersBefore, _ := q.All()
	err := q.Decline(1, 1)
	assert.Equal(t, ErrUserNotInMatch, err)
	usersAfter, _ := q.All()
	delete(usersBefore, 1)
	assert.Equal(t, usersBefore, usersAfter)
}

func TestDeclineNotInMatch(t *testing.T) {
	q := New(1000, 10, 100, time.Second*20)
	defer teardown(t, q)
	for i := 1; i < 11; i++ {
		ch, _ := q.Join(uint64(i), 1000)
		go func() {
			for range ch {
			}
		}()
	}
	time.Sleep(time.Second * 2)
	ch, _ := q.Join(11, 1000)
	go func() {
		for range ch {
		}
	}()
	time.Sleep(time.Second * 2)
	usersBefore, _ := q.All()
	err := q.Decline(11, 1)
	assert.Equal(t, ErrUserNotInMatch, err)
	usersAfter, _ := q.All()
	delete(usersBefore, 11)
	assert.Equal(t, usersBefore, usersAfter)
}

func TestDeclineAlreadyAccepted(t *testing.T) {
	q := New(1000, 10, 100, time.Second*20)
	defer teardown(t, q)
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
	q := New(1000, 10, 100, time.Second*20)
	defer teardown(t, q)
	for i := 1; i < 11; i++ {
		ch, _ := q.Join(uint64(i), 1000)
		go func() {
			for {
				_, ok := <-ch
				if !ok {
					return
				}
			}
		}()
	}
	usersBefore, _ := q.All()
	time.Sleep(time.Second * 2)
	err := q.Decline(1, 1)
	assert.NoError(t, err)

	usersAfter, _ := q.All()
	var i uint64
	for i = 2; i < 11; i++ {
		assert.Equal(t, usersBefore[i].Rating, usersAfter[i].Rating)
		assert.Equal(t, engine.QueueStateInQueue, usersAfter[i].State)
		assert.Equal(t, engine.QueueStateInQueueData{}, usersAfter[i].Data)
		assert.NotNil(t, usersAfter[i].JoinStatusChannel)
		assert.Nil(t, usersAfter[i].AcceptStatusChannel)
	}
	assert.Len(t, usersAfter, len(usersBefore)-1)
	assert.Len(t, q.groupTimers, 0)
	assert.NotContains(t, usersAfter, uint64(1))
	assert.NotContains(t, q.groups, uint64(1))
}

func TestGroupTimeout(t *testing.T) {
	q := New(1000, 10, 100, time.Second*10)
	defer teardown(t, q)
	for i := 1; i < 11; i++ {
		ch, _ := q.Join(uint64(i), 1000)
		go func() {
			for {
				_, ok := <-ch
				if !ok {
					return
				}
			}
		}()
	}
	time.Sleep(time.Second * 2)

	usersBefore, _ := q.All()
	var user1Msgs []engine.AcceptStatus
	var wg sync.WaitGroup
	wg.Add(1)
	for i := 1; i < 5; i++ {
		ch, err := q.Accept(uint64(i), 1)
		assert.NoError(t, err)
		go func(userID int) {
			for {
				msg, ok := <-ch
				if !ok {
					if userID == 1 {
						wg.Done()
					}
					return
				}
				if userID == 1 {
					user1Msgs = append(user1Msgs, msg)
				}
			}
		}(i)
	}
	time.Sleep(time.Second * 15)
	wg.Wait()

	usersAfter, _ := q.All()
	var i uint64
	for i = 1; i < 5; i++ {
		assert.Equal(t, usersBefore[i].Rating, usersAfter[i].Rating)
		assert.Equal(t, engine.QueueStateInQueue, usersAfter[i].State)
		assert.Equal(t, engine.QueueStateInQueueData{}, usersAfter[i].Data)
		assert.NotNil(t, usersAfter[i].JoinStatusChannel)
		assert.Nil(t, usersAfter[i].AcceptStatusChannel)
	}
	assert.Len(t, usersAfter, len(usersBefore)-6)
	assert.Len(t, q.groupTimers, 0)
	assert.NotContains(t, q.groups, uint64(1))
	assert.NotContains(t, usersAfter, uint64(5))
	assert.NotContains(t, usersAfter, uint64(6))
	assert.NotContains(t, usersAfter, uint64(7))
	assert.NotContains(t, usersAfter, uint64(8))
	assert.NotContains(t, usersAfter, uint64(9))
	assert.NotContains(t, usersAfter, uint64(10))
	expectedExpireMsg := engine.AcceptStatus{
		State: engine.AcceptStateExpired,
		Data:  engine.AcceptStateExpiredData{},
	}
	assert.Contains(t, user1Msgs, expectedExpireMsg)
}

func TestFoundChannelBuffered(t *testing.T) {
	q := New(1000, 10, 100, time.Second*20)
	defer teardown(t, q)
	assert.Equal(t, 10, cap(q.foundCh))
}
