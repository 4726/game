package main

import (
	"errors"
	"sync"
	"time"
)

type MatchAcceptStatus int

const (
	MatchAccepted MatchAcceptStatus = iota
	MatchDeclined
	MatchUnknown
)

type Match struct {
	sync.Mutex
	players     map[uint64]MatchAcceptStatus
	subscribers []chan MatchStatus
	cancelled   bool
	startTime   time.Time
}

type MatchStatus struct {
	TotalAccepted, TotalNeeded int
	Cancelled                  bool
	Players                    []uint64
}

var ErrUserNotInMatch = errors.New("user is not in this match")
var ErrUserAlreadyAccepted = errors.New("user already accepted")
var ErrUserAlreadyDeclined = errors.New("user already declined")
var ErrMatchCancelled = errors.New("match is cancelled")

func NewMatch(users []uint64) *Match {
	players := map[uint64]MatchAcceptStatus{}
	for _, v := range users {
		players[v] = MatchUnknown
	}

	m := &Match{
		players:     players,
		subscribers: []chan MatchStatus{},
		cancelled:   false,
		startTime:   time.Now(),
	}

	return m
}

//Accept accepts the match request and subscribes to match status updates
func (m *Match) Accept(userID uint64, ch chan MatchStatus) error {
	m.Lock()
	defer m.Unlock()

	if m.cancelled {
		return ErrMatchCancelled
	}

	var found bool
	for k, v := range m.players {
		if userID == k {
			found = true
			if v == MatchUnknown {
				m.players[k] = MatchAccepted
				break
			} else if v == MatchDeclined {
				//should not go here
				return ErrUserAlreadyDeclined
			} else if v == MatchAccepted {
				return ErrUserAlreadyAccepted
			}
		}
	}
	if !found {
		return ErrUserNotInMatch
	} else {
		m.subscribers = append(m.subscribers, ch)
		m.sendState(m.getState())
		return nil
	}
}

func (m *Match) Decline(userID uint64) error {
	m.Lock()
	defer m.Unlock()

	if m.cancelled {
		return ErrMatchCancelled
	}

	var found bool
	for k, v := range m.players {
		if userID == k {
			found = true
			if v == MatchUnknown {
				m.players[k] = MatchDeclined
				break
			} else if v == MatchAccepted {
				return ErrUserAlreadyAccepted
			} else if v == MatchDeclined {
				//should not go here
				return ErrUserAlreadyDeclined
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

func (m *Match) TimeSince() time.Duration {
	return time.Since(m.startTime)
}

func (m *Match) getState() MatchStatus {
	accepted := 0
	var cancelled bool

	players := []uint64{}

	for k, v := range m.players {
		if v == MatchAccepted {
			accepted++
		} else if v == MatchDeclined {
			cancelled = true
		}
		players = append(players, k)
	}

	return MatchStatus{
		accepted,
		len(m.players),
		cancelled,
		players,
	}
}

//sendState sends MatchStatus to all users in the match
func (m *Match) sendState(state MatchStatus) {
	if m.cancelled {
		return
	}

	for _, v := range m.subscribers {
		v <- state
		if state.Cancelled {
			close(v)
		}
	}

	if state.Cancelled {
		m.cancelled = true
	}
}
