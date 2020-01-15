package app

import (
	"context"
	"fmt"
	"time"

	"github.com/4726/game/services/matchmaking/history/config"
	"github.com/4726/game/services/matchmaking/history/pb"
	"github.com/nsqio/go-nsq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//historyServer implements pb.HistoryServer
type historyServer struct {
	consumer *nsq.Consumer
	db       *mongo.Client
	cfg      config.Config
}

func newHistoryServer(c config.Config) (*historyServer, error) {
	opts := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	db, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("could not connect to mongo: " + err.Error())
	}
	pingCtx, pingCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer pingCancel()
	if err = db.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("ping mongo error: %v", err)
	}

	consumer, err := nsq.NewConsumer(c.NSQ.Topic, c.NSQ.Channel, nsq.NewConfig())
	if err != nil {
		return nil, fmt.Errorf("could not create nsq consumer: %v", err)
	}
	consumer.AddHandler(&nsqMessageHandler{db, c.DB.Name, c.DB.Collection})
	if err := consumer.ConnectToNSQD(c.NSQ.Addr); err != nil {
		return nil, fmt.Errorf("could not connect to nsqd: %v", err)
	}

	return &historyServer{consumer, db, c}, nil
}

func (s *historyServer) Get(ctx context.Context, in *pb.GetHistoryRequest) (*pb.GetHistoryResponse, error) {
	total := in.GetTotal()
	if in.GetTotal() > s.cfg.MaxMatchResponses {
		total = s.cfg.MaxMatchResponses
	}

	matches := []*pb.MatchHistoryInfo{}
	collection := s.db.Database(s.cfg.DB.Name).Collection(s.cfg.DB.Collection)
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"end_time": -1})
	findOptions.SetLimit(int64(total))
	cur, err := collection.Find(context.Background(), bson.D{{}}, findOptions)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer cur.Close(context.Background())
	if err := cur.All(context.Background(), &matches); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetHistoryResponse{
		Match: matches,
	}, nil
}

func (s *historyServer) GetUser(ctx context.Context, in *pb.GetUserHistoryRequest) (*pb.GetUserHistoryResponse, error) {
	total := in.GetTotal()
	if in.GetTotal() > s.cfg.MaxMatchResponses {
		total = s.cfg.MaxMatchResponses
	}

	matches := []*pb.MatchHistoryInfo{}
	collection := s.db.Database(s.cfg.DB.Name).Collection(s.cfg.DB.Collection)
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"end_time": -1})
	findOptions.SetLimit(int64(total))
	//mongo shell version: {$or: [{"winner.users": 1}, {"loser.users": 1}]}
	filter := bson.M{"$or": bson.A{
		bson.M{"winner.users": in.GetUserId()},
		bson.M{"loser.users": in.GetUserId()},
	},
	}
	cur, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer cur.Close(context.Background())
	if err := cur.All(context.Background(), &matches); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetUserHistoryResponse{
		Match:  matches,
		UserId: in.GetUserId(),
	}, nil
}

func (s *historyServer) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	s.db.Disconnect(ctx)
	s.consumer.Stop()
	<-s.consumer.StopChan
}
