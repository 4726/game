package app

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
)

func TestNSQHandleMessage(t *testing.T) {
	var (
		dbSetName      = "ranking_test"
		nsqTopicName   = "ranking_test"
		nsqChannelName = "test"
		nsqAddr        = "127.0.0.1:4150"
	)

	db := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	assert.NoError(t, db.Ping().Err())
	assert.NoError(t, db.Del(dbSetName).Err())

	consumer, err := nsq.NewConsumer(nsqTopicName, nsqChannelName, nsq.NewConfig())
	assert.NoError(t, err)
	defer consumer.Stop()
	consumer.AddHandler(&nsqMessageHandler{db, dbSetName})
	assert.NoError(t, consumer.ConnectToNSQD(nsqAddr))

	producer, err := nsq.NewProducer(nsqAddr, nsq.NewConfig())
	assert.NoError(t, err)
	defer producer.Stop()

	ratings := ratingTestData()
	for _, v := range ratings {
		b, err := json.Marshal(v)
		assert.NoError(t, err)
		assert.NoError(t, producer.Publish(nsqTopicName, b))
	}

	time.Sleep(time.Second * 5)

	var actualRatings []Rating
	res, err := db.ZRangeWithScores(dbSetName, 0, -1).Result()
	assert.NoError(t, err)
	for _, v := range res {
		userID, err := strconv.ParseUint(v.Member.(string), 10, 64)
		assert.NoError(t, err)
		actualRatings = append(actualRatings, Rating{userID, uint64(v.Score)})
	}
	assert.ElementsMatch(t, ratings, actualRatings)
}

func ratingTestData() []Rating {
	r1 := Rating{1, 1000}
	r2 := Rating{2, 4000}
	r3 := Rating{3, 4505}

	return []Rating{r1, r2, r3}
}
