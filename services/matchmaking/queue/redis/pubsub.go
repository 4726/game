package queue

type PubSubTopicJoinQueue struct {
	UserID, Rating uint64
}

type PubSubTopicLeftQueue struct {
	UserID, Rating uint64
}

type PubSubTopicGroupFound struct {
	UserID, Rating uint64
}

type PubSubTopicGroupUpdate struct {
	UserID, Rating uint64
	UsersAccepted  []uint64
	UsersInGroup   []uint64
}

type PubSubTopicGroupFailed struct {
	UserID, Rating      uint64
	Expired, UserDenied bool
}

type PubSubTopicMatchFound struct {
	UserID, Rating uint64
	Users          []uint64
}

type pubSub struct {
	subscribers []chan interface{}
}

func newPubSub(subscribers []chan interface{}) *pubSub {
	return &pubSub{subscribers}
}

func (ps *pubSub) publish(msg interface{}) {
	for _, v := range ps.subscribers {
		v <- msg
	}
}
