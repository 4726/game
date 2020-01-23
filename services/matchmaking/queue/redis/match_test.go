package queue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMatchAcceptUserDoesNotExist(t *testing.T) {
	ch := make(chan MatchPubSubMessage, 1)
	m, err := NewMatches(ch)
	assert.NoError(t, err)
	assert.NoError(t, m.r.FlushAll().Err())
	assert.NoError(t, m.AddUsers(1, []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, time.Second*20))
	assert.Equal(t, ErrUserNotInMatch, m.Accept(1, 11))
	assert.Empty(t, ch)
}

func TestMatchAcceptUserAlreadyAccepted(t *testing.T) {
	ch := make(chan MatchPubSubMessage, 1)
	m, err := NewMatches(ch)
	assert.NoError(t, err)
	assert.NoError(t, m.r.FlushAll().Err())
	assert.NoError(t, m.AddUsers(1, []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, time.Second*20))
	m.Accept(1, 1)
	<-ch
	assert.Equal(t, ErrUserAlreadyAccepted, m.Accept(1, 1))
	assert.Empty(t, ch)
}

func TestMatchAcceptNoLongerExists(t *testing.T) {
	ch := make(chan MatchPubSubMessage, 1)
	m, err := NewMatches(ch)
	assert.NoError(t, err)
	assert.NoError(t, m.r.FlushAll().Err())
	assert.NoError(t, m.AddUsers(1, []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, time.Second*20))
	m.Decline(1, 2)
	<-ch
	assert.Equal(t, ErrMatchCancelled, m.Accept(1, 1))
	assert.Empty(t, ch)
}

func TestMatchAcceptExpired(t *testing.T) {
	ch := make(chan MatchPubSubMessage, 1)
	m, err := NewMatches(ch)
	assert.NoError(t, err)
	assert.NoError(t, m.r.FlushAll().Err())
	assert.NoError(t, m.AddUsers(1, []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, time.Second))
	<-ch
	assert.Equal(t, ErrUserNotInMatch, m.Accept(1, 1))
	assert.Empty(t, ch)
}

func TestMatchAccept(t *testing.T) {
	ch := make(chan MatchPubSubMessage, 1)
	m, err := NewMatches(ch)
	assert.NoError(t, err)
	assert.NoError(t, m.r.FlushAll().Err())
	assert.NoError(t, m.AddUsers(1, []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, time.Second*20))
	m.Accept(1, 1)
	msg := <-ch
	expectedMsg := MatchPubSubMessage{1, MatchStatus{
		TotalAccepted: 1,
		TotalNeeded:   10,
		Players:       testDefaultMatchPlayers(),
		Expired:       false,
		Cancelled:     false,
	}}
	expectedMsg.State.Players[uint64(1)] = MatchAccepted
	assert.Equal(t, expectedMsg, msg)
}

func TestMatchDeclineUserDoesNotExist(t *testing.T) {
	ch := make(chan MatchPubSubMessage, 1)
	m, err := NewMatches(ch)
	assert.NoError(t, err)
	assert.NoError(t, m.r.FlushAll().Err())
	assert.NoError(t, m.AddUsers(1, []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, time.Second*20))
	assert.Equal(t, ErrUserNotInMatch, m.Decline(1, 11))
	assert.Empty(t, ch)
}

func TestMatchDeclineUserAlreadyAccepted(t *testing.T) {
	ch := make(chan MatchPubSubMessage, 1)
	m, err := NewMatches(ch)
	assert.NoError(t, err)
	assert.NoError(t, m.r.FlushAll().Err())
	assert.NoError(t, m.AddUsers(1, []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, time.Second*20))
	m.Accept(1, 1)
	<-ch
	assert.Equal(t, ErrUserAlreadyAccepted, m.Decline(1, 1))
	assert.Empty(t, ch)
}

func TestMatchDeclineNoLongerExists(t *testing.T) {
	ch := make(chan MatchPubSubMessage, 1)
	m, err := NewMatches(ch)
	assert.NoError(t, err)
	assert.NoError(t, m.r.FlushAll().Err())
	assert.NoError(t, m.AddUsers(1, []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, time.Second*20))
	m.Decline(1, 2)
	<-ch
	assert.Equal(t, ErrMatchCancelled, m.Decline(1, 1))
	assert.Empty(t, ch)
}

func TestMatchDeclineExpired(t *testing.T) {
	ch := make(chan MatchPubSubMessage, 1)
	m, err := NewMatches(ch)
	assert.NoError(t, err)
	assert.NoError(t, m.r.FlushAll().Err())
	assert.NoError(t, m.AddUsers(1, []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, time.Second))
	<-ch
	assert.Equal(t, ErrUserNotInMatch, m.Decline(1, 1))
	assert.Empty(t, ch)
}

func TestMatchDecline(t *testing.T) {
	ch := make(chan MatchPubSubMessage, 1)
	m, err := NewMatches(ch)
	assert.NoError(t, err)
	assert.NoError(t, m.r.FlushAll().Err())
	assert.NoError(t, m.AddUsers(1, []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, time.Second*20))
	m.Decline(1, 1)
	msg := <-ch
	expectedMsg := MatchPubSubMessage{1, MatchStatus{
		Cancelled: true,
	}}
	assert.Equal(t, expectedMsg, msg)
}

//for testing
func testDefaultMatchPlayers(userIDs ...uint64) map[uint64]MatchAcceptStatus {
	m := map[uint64]MatchAcceptStatus{}

	for _, v := range userIDs {
		m[v] = MatchUnknown
	}
	if len(m) > 0 {
		return m
	}

	for i := 1; i < 11; i++ {
		m[uint64(i)] = MatchUnknown
	}
	return m
}
