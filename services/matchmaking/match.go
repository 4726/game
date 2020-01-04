package main

import (
	"errors"
	"sync"
	"time"
)

type MatchAcceptStatus int

const (
	MatchAccepted MatchAcceptStatus = iota
	MatchDenied
	MatchUnknown
)

type Match struct {
	sync.Mutex
	players map[*QueueData]MatchAcceptStatus
	chs     []chan MatchStatus
}

type MatchStatus struct {
	TotalAccepted, TotalNeeded int
	Cancelled                  bool
}

var ErrUserNotInMatch = errors.New("user is not in this match")

func NewMatch(users []uint64, timeout time.Duration) *Match {
	players := map[*QueueData]MatchAcceptStatus{}
	for _, v := range queueData {
		players[v] = MatchUnknown
	}

	m := &Match{
		players: players,
		chs:     []chan MatchStatus{},
	}

	time.AfterFunc(timeout, func() {
		m.Lock()
		defer m.Unlock()

		state := m.getState()
		if state.TotalAccepted != state.TotalNeeded {
			state.Cancelled = true
			m.sendState(state)
		}
	})

	return m
}

//Accept accepts the match request and subscribes to match status updates
func (m *Match) Accept(userID uint64, ch chan MatchStatus) error {
	m.Lock()
	defer m.Unlock()

	m.chs = append(m.chs, ch)

	var found bool
	for k, v := range m.players {
		if userID == k.UserID {
			found = true
			if v == MatchUnknown {
				m.players[k] = MatchAccepted
				break
			} else if v == MatchDenied {
				//user already denied
				return nil
			}
		}
	}
	if !found {
		return ErrUserNotInMatch
	} else {
		m.sendState(m.getState())
		return nil
	}
}

func (m *Match) Decline(userID uint64) error {
	m.Lock()
	defer m.Unlock()

	var found bool
	for k, v := range m.players {
		if userID == k.UserID {
			if v == MatchUnknown {
				m.players[k] = MatchDenied
				break
			} else if v == MatchAccepted {
				//user already accepted
				return nil
			}
		}
	}
	if !found {
		return ErrUserNotInMatch
	} else {
		m.sendState(m.getState())
		return nil
	}
}

func (m *Match) getState() MatchStatus {
	accepted := 0
	var cancelled bool

	for _, v := range m.players {
		if v == MatchAccepted {
			accepted++
		} else if v == MatchDenied {
			cancelled = true
		}
	}

	return MatchStatus{
		accepted,
		len(m.players),
		cancelled,
	}
}

//sendState sends MatchStatus to all users in the match
func (m *Match) sendState(state MatchStatus) {
	for _, v := range m.chs {
		v <- state
	}
}
