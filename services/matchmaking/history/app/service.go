package app

import (
	"context"
	"errors"
	"time"

	"github.com/4726/game/services/matchmaking/history/pb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nsqio/go-nsq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	consumer *nsq.Consumer
	db       *mongo.Client
	cfg      Config
}

type Config struct {
	DB  DBConfig
	NSQ NSQConfig
}

type DBConfig struct {
	Name, Collection string
}

type NSQConfig struct {
	Addr, Topic, Channel string
}

const maxMatchResponses = 100

func NewService(c Config) (*Service, error) {
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

	consumer, err := nsq.NewConsumer(c.NSQ.Topic, c.NSQ.Channel, nsq.NewConfig())
	if err != nil {
		return nil, err
	}
	consumer.AddHandler(&nsqMessageHandler{db, c.DB.Name, c.DB.Collection})
	if err := consumer.ConnectToNSQD(c.NSQ.Addr); err != nil {
		return nil, err
	}

	return &Service{consumer, db, c}, nil
}

func (s *Service) Get(ctx context.Context, in *pb.GetHistoryRequest) (*pb.GetHistoryResponse, error) {
	if in.GetTotal() > maxMatchResponses {
		return nil, errors.New("invalid total")
	}
	matches := []*pb.MatchHistoryInfo{}
	collection := s.db.Database(s.cfg.DB.Name).Collection(s.cfg.DB.Collection)
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"end_time": -1})
	findOptions.SetLimit(int64(in.GetTotal()))
	cur, err := collection.Find(context.Background(), bson.D{{}}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	if err := cur.All(context.Background(), &matches); err != nil {
		return nil, err
	}

	return &pb.GetHistoryResponse{
		Match: matches,
	}, nil
}

func (s *Service) GetUser(ctx context.Context, in *pb.GetUserHistoryRequest) (*pb.GetUserHistoryResponse, error) {
	if in.GetTotal() > maxMatchResponses {
		return nil, errors.New("invalid total")
	}

	matches := []*pb.MatchHistoryInfo{}
	collection := s.db.Database(s.cfg.DB.Name).Collection(s.cfg.DB.Collection)
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"end_time": -1})
	findOptions.SetLimit(int64(in.GetTotal()))
	//mongo shell version: {$or: [{"winner.users": 1}, {"loser.users": 1}]}
	filter := bson.M{"$or": bson.A{
		bson.M{"winner.users": in.GetUserId()},
		bson.M{"loser.users": in.GetUserId()},
	},
	}
	cur, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	if err := cur.All(context.Background(), &matches); err != nil {
		return nil, err
	}

	return &pb.GetUserHistoryResponse{
		Match:  matches,
		UserId: in.GetUserId(),
	}, nil
}

func (s *Service) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	s.db.Disconnect(ctx)
	s.consumer.Stop()
	<-s.consumer.StopChan
}
