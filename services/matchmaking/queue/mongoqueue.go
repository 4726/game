package main

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoQueue struct {
	sync.Mutex
	db                   *mongo.Client
	limit                int
	subscribers          []chan PubSubMessage
	matchID              uint64
	dbName, dbCollection string
}

func NewMongoQueue(limit int) (*MongoQueue, error) {
	opts := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	db, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	pingCtx, pingCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer pingCancel()
	if err = db.Ping(pingCtx, nil); err != nil {
		return nil, err
	}

	return &MongoQueue{
		db:           db,
		limit:        limit,
		subscribers:  []chan PubSubMessage{},
		matchID:      uint64(0),
		dbName:       "tempname",
		dbCollection: "tempcollname",
	}, err
}

func (q *MongoQueue) Enqueue(userID, rating uint64) error {
	// if len(q.data) >= q.limit {
	// 	return ErrQueueFull
	// }

	qd := QueueData{userID, rating, time.Now(), false, 0}

	collection := q.db.Database(q.dbName).Collection(q.dbCollection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	filter := bson.M{"userid": userID}
	update := bson.M{"$set": qd}
	opts := options.Update().SetUpsert(true)
	res, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	if res.MatchedCount == 1 {
		return ErrAlreadyInQueue
	}

	q.publish(PubSubTopicAdd, qd)
	return nil
}

func (q *MongoQueue) Dequeue() {
	collection := q.db.Database(q.dbName).Collection(q.dbCollection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := options.FindOneAndDelete()
	opts.SetSort(bson.M{"_id": 1})
	res := collection.FindOneAndDelete(ctx, bson.D{{}}, opts)
	var qd QueueData
	res.Decode(&qd)
	q.publish(PubSubTopicDelete, qd)
}

func (q *MongoQueue) DeleteOne(userID uint64) error {
	collection := q.db.Database(q.dbName).Collection(q.dbCollection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	filter := bson.M{"userid": userID}
	opts := options.FindOneAndDelete()
	opts.SetSort(bson.M{"_id": 1})
	res := collection.FindOneAndDelete(ctx, filter, opts)
	var qd QueueData
	res.Decode(&qd)

	if qd.UserID == 0 {
		return ErrDoesNotExist
	}
	q.publish(PubSubTopicDelete, qd)
	return nil
}

//Len returns length of queue
func (q *MongoQueue) Len() int {
	collection := q.db.Database(q.dbName).Collection(q.dbCollection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	count, _ := collection.CountDocuments(ctx, bson.D{{}})
	return int(count)
}

func (q *MongoQueue) MarkMatchFound(userID uint64, found bool) error {
	var pubSubTopic PubSubTopic
	if found {
		pubSubTopic = PubSubTopicMatchFound
	} else {
		pubSubTopic = PubSubTopicMatchNotFound
	}

	collection := q.db.Database(q.dbName).Collection(q.dbCollection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	filter := bson.M{"userid": userID}
	update := bson.M{"$set": bson.M{"matchfound": found}}
	res := collection.FindOneAndUpdate(ctx, filter, update)
	var oldQD QueueData
	res.Decode(&oldQD)
	if oldQD.UserID == 0 {
		return ErrDoesNotExist
	}
	if oldQD.MatchFound == found {
		//already set to found, don't need to send pubsub message
		return nil
	}

	oldQD.MatchFound = found

	q.publish(pubSubTopic, oldQD)
	return nil
}

func (q *MongoQueue) EnqueueAndFindMatch(userID, rating, ratingRange uint64, total int) (found bool, matchID uint64, qds []QueueData, err error) {
	// if len(q.data) >= q.limit {
	// 	err = ErrQueueFull
	// 	return
	// }

	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	var sess mongo.Session
	sess, err = q.db.StartSession(opts)
	if err != nil {
		return
	}
	defer sess.EndSession(context.TODO())

	qd := QueueData{userID, rating, time.Now(), false, 0}

	txnOpts := options.Transaction().SetReadPreference(readpref.PrimaryPreferred())
	var queueDatas interface{}
	queueDatas, err = sess.WithTransaction(context.TODO(), func(sessCtx mongo.SessionContext) (interface{}, error) {
		collection := q.db.Database(q.dbName).Collection(q.dbCollection)
		filter := bson.M{"userid": userID}
		update := bson.M{"$set": qd}
		opts := options.Update().SetUpsert(true)
		res, err := collection.UpdateOne(sessCtx, filter, update, opts)
		if err != nil {
			return nil, err
		}
		if res.MatchedCount == 1 {
			return nil, ErrAlreadyInQueue
		}

		ratingLessThan := qd.Rating + ratingRange/2
		ratingGreaterThan := qd.Rating - ratingRange/2

		findOptions := options.Find()
		findOptions.SetSort(bson.M{"_id": 1})
		findOptions.SetLimit(int64(total))
		findFilter := bson.M{"$and": bson.A{
			bson.M{"rating": bson.M{"$gt": ratingGreaterThan}},
			bson.M{"rating": bson.M{"$lt": ratingLessThan}},
			bson.M{"matchfound": false},
		},
		}
		cur, err := collection.Find(sessCtx, findFilter, findOptions)
		if err != nil {
			return nil, err
		}
		defer cur.Close(sessCtx)

		var qds []QueueData
		if err := cur.All(sessCtx, &qds); err != nil {
			return nil, err
		}
		if len(qds) < total {
			return []QueueData{}, nil
		}
		return qds, nil
	}, txnOpts)
	if err != nil {
		return
	}

	qds = queueDatas.([]QueueData)
	q.publish(PubSubTopicAdd, qd)

	if len(qds) == 0 {
		//user added to queue, no match found
		return
	}

	found = true
	matchID = q.getMatchID()
	for _, v := range qds {
		q.publish(PubSubTopicMatchFound, v)
	}

	return
}

func (q *MongoQueue) Subscribe(ch chan PubSubMessage) {
	q.subscribers = append(q.subscribers, ch)
}

func (q *MongoQueue) publish(topic PubSubTopic, qd QueueData) {
	for _, v := range q.subscribers {
		go func(ch chan PubSubMessage) {
			select {
			case ch <- PubSubMessage{topic, qd}:
			case <-time.After(time.Second * 10):
			}
		}(v)
	}
}

func (q *MongoQueue) getMatchID() uint64 {
	q.matchID++
	return q.matchID
}
