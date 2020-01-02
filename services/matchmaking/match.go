package main

type MatchAcceptStatus int

const (
	MatchAccepted MatchAcceptStatus = iota
	MatchDenied
	MatchUnknown
)

type Match struct {
	sync.Mutex
	players map[QueueData]MatchAcceptStatus
	queueType pb.QueueType
	chs chan *pb.AcceptQueueResponse
}

var ErrUserNotInMatch = errors.New("user is not in this match")

func NewMatch(queueData []QueueData, qt pb.QueueType) Match {
	players := map[QueueData]MatchAcceptStatus{}
	for _, v := range players {
		players[v] = MatchUnknown
	}
	return &Match{
		players,
		qt,
	}
}

func (m *Match) Accept(userID uint64, ch chan *pb.AcceptQueueResponse) error {
	m.Lock()
	defer m.Unlock()

	m.chs = append(m.chs, ch)

	var found bool
	for k, v := range m.players {
		if userID == k.UserID {
			found = true
			if v == MatchUnknown {
				m.players[k] = MatchAccepted
				sendState(m.getState())
				return nil
			} else if v == MatchDenied {
				//user alreadu denied
				sendState(m.getState())
				return nil
			}
		}
	}
	if !found {
		return ErrUserNotInMatch
	} else {
		sendState(m.getState())
		return nil
	}
}

func (m *Match) Decline(userID uint64) error {
	m.Lock()
	defer m.Unlock()

	m.chs = append(m.chs, ch)

	var found bool
	for k, v := range m.players {
		if userID == k.UserID {
			if v == MatchUnknown {
				m.players[k] = MatchDenied
				return nil
			} else if v == MatchAccepted {
				//user already accepted
				return nil
			}
		}
	}
	if !found {
		return ErrUserNotInMatch
	} else {
		return nil
	}
}

func (m *Match) getState() *pb.AcceptQueueResponse {
	accepted := 0
	var cancelled bool

	for _, v := range m.players {
		if v == MatchAccepted {
			accepted++
		} else if v == MatchDenied {
			cancelled = true
		}
	}

	return &pb.AcceptQueueResponse{
		accepted,
		len(m.players),
		m.queueType,
		cancelled,
	}
}

func (m *Match) sendState(resp *pb.AcceptQueueResponse) {
	for _, v := range m.chs {
		v <- resp
	}
}