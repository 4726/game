package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func assertMatchStatusEqual(t testing.TB, expected, actual MatchStatus) {
	assert.Equal(t, expected.TotalAccepted, actual.TotalAccepted)
	assert.Equal(t, expected.TotalNeeded, actual.TotalNeeded)
	assert.Equal(t, expected.Cancelled, actual.Cancelled)
	assert.ElementsMatch(t, expected.Players, actual.Players)
}

func TestMatchAcceptUserDoesNotExist(t *testing.T) {
	m := NewMatch([]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, time.Second*20)
	ch := make(chan MatchStatus, 1)
	assert.Equal(t, ErrUserNotInMatch, m.Accept(11, ch))
	assert.Empty(t, ch)
}

func TestMatchAcceptUserAlreadyAccepted(t *testing.T) {
	m := NewMatch([]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, time.Second*20)
	ch := make(chan MatchStatus, 1)
	m.Accept(1, ch)
	<-ch
	assert.Equal(t, ErrUserAlreadyAccepted, m.Accept(1, ch))
	assert.Empty(t, ch)
}

func TestMatchAcceptUserAlreadyDeclined(t *testing.T) {
	m := NewMatch([]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, time.Second*20)
	ch := make(chan MatchStatus, 1)
	m.Decline(1)
	assert.Equal(t, ErrUserAlreadyDeclined, m.Accept(1, ch))
	assert.Empty(t, ch)
}

func TestMatchAcceptTimeoutBefore(t *testing.T) {
	players := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	m := NewMatch(players, time.Second*2)
	time.Sleep(time.Second * 3)

	ch := make(chan MatchStatus, 1)
	assert.NoError(t, m.Accept(1, ch))
	status := <-ch
	expectedStatus := MatchStatus{1, 10, true, players}
	assertMatchStatusEqual(t, expectedStatus, status)
	assert.Empty(t, ch)
}

func TestMatchAcceptTimeoutAfter(t *testing.T) {
	players := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	m := NewMatch(players, time.Second*2)
	ch := make(chan MatchStatus, 1)
	assert.NoError(t, m.Accept(1, ch))
	<-ch
	time.Sleep(time.Second * 3)

	status := <-ch
	expectedStatus := MatchStatus{1, 10, true, players}
	assertMatchStatusEqual(t, expectedStatus, status)
	assert.Empty(t, ch)
}

func TestMatchAccept(t *testing.T) {
	players := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	m := NewMatch(players, time.Second*20)
	ch := make(chan MatchStatus, 1)
	assert.NoError(t, m.Accept(1, ch))

	status := <-ch
	expectedStatus := MatchStatus{1, 10, false, players}
	assertMatchStatusEqual(t, expectedStatus, status)
	assert.Empty(t, ch)

	ch2 := make(chan MatchStatus, 1)
	m.Accept(2, ch2)

	status = <-ch
	expectedStatus = MatchStatus{2, 10, false, players}
	assertMatchStatusEqual(t, expectedStatus, status)
	assert.Empty(t, ch)

	status = <-ch2
	assertMatchStatusEqual(t, expectedStatus, status)
	assert.Empty(t, ch2)
}

func TestMatchDeclineUserDoesNotExist(t *testing.T) {
	m := NewMatch([]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, time.Second*20)
	assert.Equal(t, ErrUserNotInMatch, m.Decline(11))
}

func TestMatchDeclineUserAlreadyAccepted(t *testing.T) {
	m := NewMatch([]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, time.Second*20)
	ch := make(chan MatchStatus, 1)
	m.Accept(1, ch)
	<-ch
	assert.Equal(t, ErrUserAlreadyAccepted, m.Decline(1))
	assert.Empty(t, ch)
}

func TestMatchDeclineUserAlreadyDeclined(t *testing.T) {
	m := NewMatch([]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, time.Second*20)
	m.Decline(1)
	assert.Equal(t, ErrUserAlreadyDeclined, m.Decline(1))
}

func TestMatchDecline(t *testing.T) {
	players := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	m := NewMatch(players, time.Second*20)
	assert.NoError(t, m.Decline(1))

	ch := make(chan MatchStatus, 1)
	assert.NoError(t, m.Accept(2, ch))

	status := <-ch
	expectedStatus := MatchStatus{1, 10, true, players}
	assertMatchStatusEqual(t, expectedStatus, status)
	assert.Empty(t, ch)
}
