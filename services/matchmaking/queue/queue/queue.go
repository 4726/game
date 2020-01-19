package queue

import (
	"errors"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type PubSubTopic int

const (
	PubSubTopicAdd PubSubTopic = iota
	PubSubTopicDelete
	PubSubTopicMatchFound
	PubSubTopicMatchNotFound
)

type QueueData struct {
	UserID     uint64 `gorm:"PRIMARY_KEY;NOT NULL"`
	Rating     uint64
	StartTime  time.Time
	MatchFound bool
	MatchID    uint64
}

type PubSubMessage struct {
	Topic PubSubTopic
	Data  QueueData
}

type Queue struct {
	sync.Mutex
	db          *gorm.DB
	limit       int
	subscribers []chan PubSubMessage
	matchID     uint64
}

var ErrAlreadyInQueue = errors.New("user already in queue")
var ErrDoesNotExist = errors.New("does not exist")
var ErrNotAllExists = errors.New("not all exists")
var ErrQueueFull = errors.New("queue full")

func New(limit int) (*Queue, error) {
	s := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True", "root", "pass", "127.0.0.1:3306", "tempname")

	db, err := gorm.Open("mysql", s)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&QueueData{})

	return &Queue{
		db:          db,
		limit:       limit,
		subscribers: []chan PubSubMessage{},
		matchID:     uint64(0),
	}, err
}

func (q *Queue) Enqueue(userID, rating uint64) error {
	// if len(q.data) >= q.limit {
	// 	return ErrQueueFull
	// }

	qd := QueueData{userID, rating, time.Now(), false, 0}

	db := q.db.FirstOrCreate(&qd)
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected < 1 {
		return ErrAlreadyInQueue
	}

	q.publish(PubSubTopicAdd, qd)
	return nil
}

func (q *Queue) DeleteOne(userID uint64) error {
	qd := QueueData{UserID: userID}
	db := q.db.Delete(&qd)
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected < 1 {
		return ErrDoesNotExist
	}
	q.publish(PubSubTopicDelete, qd)
	return nil
}

func (q *Queue) All() ([]QueueData, error) {
	var qd []QueueData

	if err := q.db.Find(&qd).Error; err != nil {
		return qd, err
	}

	return qd, nil
}

//Len returns length of queue
func (q *Queue) Len() int {
	var count int
	q.db.Table("queue_data").Count(&count)
	return count
}

func (q *Queue) MarkMatchFound(userID uint64, found bool) error {
	var pubSubTopic PubSubTopic
	if found {
		pubSubTopic = PubSubTopicMatchFound
	} else {
		pubSubTopic = PubSubTopicMatchNotFound
	}

	qd := QueueData{UserID: userID}
	res := q.db.Model(&qd).Update("match_found", found)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected < 1 {
		//either does not exist or matchfound was already set to found
		return nil

	}

	q.publish(pubSubTopic, qd)
	return nil
}

func (q *Queue) EnqueueAndFindMatch(userID, rating, ratingRange uint64, total int) (found bool, matchID uint64, qds []QueueData, err error) {
	// if len(q.data) >= q.limit {
	// 	err = ErrQueueFull
	// 	return
	// }

	qd := QueueData{userID, rating, time.Now(), false, 0}

	tx := q.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Error; err != nil {
		return
	}
	res := tx.FirstOrCreate(&qd)
	if res.Error != nil {
		tx.Rollback()
		err = res.Error
		return
	}
	if res.RowsAffected < 1 {
		tx.Commit()
		err = ErrAlreadyInQueue
		return
	}

	ratingLessThan := qd.Rating + ratingRange/2
	ratingGreaterThan := qd.Rating - ratingRange/2

	qds = []QueueData{}
	err = tx.
		Where("rating <= ? AND rating >= ? AND match_found = ?", ratingLessThan, ratingGreaterThan, false).
		Order("start_time asc").
		Limit(total).
		Find(&qds).Error
	if err != nil {
		tx.Rollback()
		return
	}

	var tempMatchID uint64
	if len(qds) >= total {
		tempMatchID = q.getMatchID()
		for i, v := range qds {
			err = tx.Model(&v).Update("match_found", true).Error
			if err != nil {
				tx.Rollback()
				return
			}
			v.MatchFound = true
			v.MatchID = tempMatchID
			qds[i] = v
		}
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return
	}

	q.publish(PubSubTopicAdd, qd)

	if len(qds) < total {
		qds = []QueueData{}
		return
	}

	found = true
	matchID = tempMatchID
	for _, v := range qds {
		q.publish(PubSubTopicMatchFound, v)
	}

	return
}

func (q *Queue) Subscribe(ch chan PubSubMessage) {
	q.subscribers = append(q.subscribers, ch)
}

func (q *Queue) publish(topic PubSubTopic, qd QueueData) {
	for _, v := range q.subscribers {
		go func(ch chan PubSubMessage) {
			select {
			case ch <- PubSubMessage{topic, qd}:
			case <-time.After(time.Second * 10):
			}
		}(v)
	}
}

func (q *Queue) getMatchID() uint64 {
	q.matchID++
	return q.matchID
}
