package main

import (
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type MysqlQueue struct {
	sync.Mutex
	db          *gorm.DB
	limit       int
	subscribers []chan PubSubMessage
	matchID     uint64
}

func NewMysqlQueue(limit int) (*MysqlQueue, error) {
	s := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True", "root", "pass", "127.0.0.1:3306", "tempname")

	db, err := gorm.Open("mysql", s)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&QueueData{})

	return &MysqlQueue{
		db:          db,
		limit:       limit,
		subscribers: []chan PubSubMessage{},
		matchID:     uint64(0),
	}, err
}

func (q *MysqlQueue) Enqueue(userID, rating uint64) error {
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

func (q *MysqlQueue) DeleteOne(userID uint64) error {
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

//Len returns length of queue
func (q *MysqlQueue) Len() int {
	var count int
	q.db.Table("queue_data").Count(&count)
	return count
}

func (q *MysqlQueue) MarkMatchFound(userID uint64, found bool) error {
	var pubSubTopic PubSubTopic
	if found {
		pubSubTopic = PubSubTopicMatchFound
	} else {
		pubSubTopic = PubSubTopicMatchNotFound
	}

	qd := QueueData{UserID: userID}
	if err := q.db.Model(&qd).Update("matchfound", found).Error; err != nil {
		return err
	}
	t := time.Time{}
	if qd.StartTime == t {
		return ErrDoesNotExist
	}
	if qd.MatchFound == found {
		//already set to found, don't set pubsubmessage
		return nil
	}

	q.publish(pubSubTopic, qd)
	return nil
}

func (q *MysqlQueue) EnqueueAndFindMatch(userID, rating, ratingRange uint64, total int) (found bool, matchID uint64, qds []QueueData, err error) {
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

	db := q.db.FirstOrCreate(&qd)
	if db.Error != nil {
		tx.Rollback()
		err = db.Error
		return
	}
	if db.RowsAffected < 1 {
		tx.Commit()
		err = ErrAlreadyInQueue
		return
	}

	qds = []QueueData{}
	ratingLessThan := qd.Rating + ratingRange/2
	ratingGreaterThan := qd.Rating - ratingRange/2

	err = q.db.
		Where("rating <= ? AND rating >= ? && match_found = ?", ratingLessThan, ratingGreaterThan, false).
		Order("starttime asc").
		Limit(total).
		Find(&qds).Error
	if err != nil {
		tx.Rollback()
		return
	}

	if len(qds) < total {
		tx.Commit()
		return
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return
	}

	q.publish(PubSubTopicAdd, qd)

	found = true
	matchID = q.getMatchID()
	for _, v := range qds {
		q.publish(PubSubTopicMatchFound, v)
	}

	return
}

func (q *MysqlQueue) Subscribe(ch chan PubSubMessage) {
	q.subscribers = append(q.subscribers, ch)
}

func (q *MysqlQueue) publish(topic PubSubTopic, qd QueueData) {
	for _, v := range q.subscribers {
		go func(ch chan PubSubMessage) {
			select {
			case ch <- PubSubMessage{topic, qd}:
			case <-time.After(time.Second * 10):
			}
		}(v)
	}
}

func (q *MysqlQueue) getMatchID() uint64 {
	q.matchID++
	return q.matchID
}
