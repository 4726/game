package queue

type Queue interface {
	Join(userID, rating uint64) (<-chan JoinStatus, error)
	Leave(userID uint64) error
	Accept(userID, matchID uint64) (<-chan AcceptStatus, error)
	Decline(userID, matchID uint64) error
	All() (map[uint64]UserData, error)
	Channel() <-chan Match
}

type UserData struct {
	Rating              uint64
	State               QueueState
	Data                interface{}
	JoinStatusChannel   chan JoinStatus
	AcceptStatusChannel chan AcceptStatus
}

type QueueState int

const (
	QueueStateInQueue QueueState = iota
	QueueStateInGroup
)

type QueueStateInQueueData struct{}

type QueueStateInGroupData struct {
	Accepted, Denied bool
	MatchID          uint64
}

type JoinStatus struct {
	State JoinState
	Data  interface{}
}

type JoinState int

const (
	JoinStateEntered JoinState = iota
	JoinStateLeft
	JoinStateGroupFound
)

type JoinStateEnteredData struct{}

type JoinStateLeftData struct{}

type JoinStateGroupFoundData struct {
	MatchID uint64
}

type AcceptStatus struct {
	State AcceptState
	Data  interface{}
}

type AcceptState int

const (
	AcceptStateUpdate AcceptState = iota
	AcceptStateFailed
	AcceptStateExpired
	AcceptStateSuccess
)

type AcceptStatusUpdateData struct {
	UsersNeeded, UsersAccepted int
}

type AcceptStateFailedData struct{}

type AcceptStateExpiredData struct{}

type AcceptStateSuccessData struct {
	UserCount int
	MatchID   uint64
}

type Match struct {
	Users   map[uint64]uint64
	MatchID uint64
}
