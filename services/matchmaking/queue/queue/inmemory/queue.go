package inmemory

import (
	"errors"
	"sync"

	"github.com/4726/game/services/matchmaking/queue/queue"
)

type Queue struct {
	sync.Mutex
	data        map[uint64]queue.UserData
	limit       int
	matchID     uint64
	perMatch    int
	ratingRange int
	groups      map[uint64]map[uint64]struct{}
}

var ErrAlreadyInQueue = errors.New("user already in queue")
var ErrDoesNotExist = errors.New("does not exist")
var ErrQueueFull = errors.New("queue full")
var ErrUserNotInMatch = errors.New("user is not in this match")
var ErrUserAlreadyAccepted = errors.New("user already accepted")

func New(limit, perMatch, ratingRange int) *Queue {
	return &Queue{
		data:        map[uint64]queue.UserData{},
		limit:       limit,
		matchID:     uint64(0),
		perMatch:    perMatch,
		ratingRange: ratingRange,
		groups:      map[uint64]map[uint64]struct{}{},
	}
}

func (q *Queue) Join(userID, rating uint64) (<-chan queue.JoinStatus, error) {
	q.Lock()
	defer q.Unlock()

	if len(q.data) >= q.limit {
		return nil, ErrQueueFull
	}

	if _, ok := q.data[userID]; ok {
		return nil, ErrAlreadyInQueue
	}

	joinStatusChannel := make(chan queue.JoinStatus)

	q.data[userID] = queue.UserData{
		Rating:              rating,
		State:               queue.QueueStateInQueue,
		Data:                queue.QueueStateInQueueData{},
		JoinStatusChannel:   joinStatusChannel,
		AcceptStatusChannel: nil,
	}

	go func() {
		joinStatusChannel <- queue.JoinStatus{
			State: queue.JoinStateEntered,
			Data:  queue.JoinStateEnteredData{},
		}
		q.searchMatch(userID)
	}()

	return joinStatusChannel, nil
}

func (q *Queue) Leave(userID uint64) error {
	q.Lock()
	defer q.Unlock()

	if _, exists := q.data[userID]; !exists {
		return ErrDoesNotExist
	}

	delete(q.data, userID)
	return nil
}

func (q *Queue) Accept(userID, matchID uint64) (<-chan queue.AcceptStatus, error) {
	q.Lock()
	defer q.Unlock()

	usersInMatch := q.groups[matchID]
	_, ok := usersInMatch[userID]
	if !ok {
		delete(q.data, userID)
		return nil, ErrUserNotInMatch
	}
	userData := q.data[userID]
	groupData := userData.Data.(queue.QueueStateInGroupData)
	if groupData.Accepted {
		return nil, ErrUserAlreadyAccepted
	}

	ch := make(chan queue.AcceptStatus)
	setQueueStateInGroup(&userData, queue.QueueStateInGroupData{
		Accepted: true,
		Denied:   false,
		MatchID:  matchID,
	})
	userData.AcceptStatusChannel = ch
	q.data[userID] = userData
	go q.sendMatchUpdate(matchID)
	return ch, nil
}

func (q *Queue) Decline(userID, matchID uint64) error {
	q.Lock()
	defer q.Unlock()

	usersInMatch := q.groups[matchID]
	_, ok := usersInMatch[userID]
	if !ok {
		delete(q.data, userID)
		return ErrUserNotInMatch
	}
	userData := q.data[userID]
	groupData := userData.Data.(queue.QueueStateInGroupData)
	if groupData.Accepted {
		return ErrUserAlreadyAccepted
	}

	for k := range usersInMatch {
		userData := q.data[k]
		if userData.AcceptStatusChannel != nil {
			userData.AcceptStatusChannel <- queue.AcceptStatus{
				State: queue.AcceptStateFailed,
				Data:  queue.AcceptStateFailedData{},
			}
			close(userData.AcceptStatusChannel)
		}
		setQueueStateInQueue(&userData, queue.QueueStateInQueueData{})
		userData.AcceptStatusChannel = nil
		q.data[k] = userData
	}
	delete(q.data, userID)
	return nil
}

func (q *Queue) All() (map[uint64]queue.UserData, error) {
	q.Lock()
	defer q.Unlock()

	m := map[uint64]queue.UserData{}
	for k, v := range q.data {
		m[k] = v
	}

	return m, nil
}

func setQueueStateInQueue(userData *queue.UserData, data queue.QueueStateInQueueData) {
	userData.State = queue.QueueStateInQueue
	userData.Data = data
}

func setQueueStateInGroup(userData *queue.UserData, data queue.QueueStateInGroupData) {
	userData.State = queue.QueueStateInGroup
	userData.Data = data
}

func (q *Queue) searchMatch(userID uint64) {
	q.Lock()
	defer q.Unlock()

	userData, ok := q.data[userID]
	if !ok {
		return
	}

	ratingLessThan := userData.Rating + uint64(q.ratingRange/2)
	ratingGreaterThan := userData.Rating - uint64(q.ratingRange/2)

	suitableUsers := map[uint64]struct{}{}
	for k, v := range q.data {
		if v.Rating <= ratingLessThan && v.Rating >= ratingGreaterThan && v.State == queue.QueueStateInQueue {
			if len(suitableUsers) > q.perMatch {
				break
			}
			suitableUsers[k] = struct{}{}
		}
	}

	if len(suitableUsers) < q.perMatch {
		return
	}

	matchID := q.getMatchID()
	for k := range suitableUsers {
		userData := q.data[k]
		userData.JoinStatusChannel <- queue.JoinStatus{
			State: queue.JoinStateGroupFound,
			Data: queue.JoinStateGroupFoundData{
				MatchID: matchID,
			},
		}
		setQueueStateInGroup(&userData, queue.QueueStateInGroupData{
			Accepted: false,
			Denied:   false,
			MatchID:  matchID,
		})
		q.data[k] = userData
	}
	q.groups[matchID] = suitableUsers
}

//need to handle when all users acce[ted]
func (q *Queue) sendMatchUpdate(matchID uint64) {
	q.Lock()
	defer q.Unlock()

	usersInMatch := q.groups[matchID]
	var accepted int
	for k := range usersInMatch {
		userData := q.data[k]
		queueData := userData.Data.(queue.QueueStateInGroupData)
		if queueData.Accepted {
			accepted++
		}
	}

	for k := range usersInMatch {
		userData := q.data[k]
		if userData.AcceptStatusChannel != nil {
			userData.AcceptStatusChannel <- queue.AcceptStatus{
				State: queue.AcceptStateUpdate,
				Data: queue.AcceptStatusUpdateData{
					UsersAccepted: accepted,
					UsersNeeded:   len(usersInMatch),
				},
			}
		}
	}
}

func (q *Queue) getMatchID() uint64 {
	q.matchID++
	return q.matchID
}
