package app

import (
	"context"
	"fmt"
	"time"

	"github.com/4726/game/services/matchmaking/live/config"
	"github.com/4726/game/services/matchmaking/live/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type liveServer struct {
	db  *mongo.Client
	cfg config.Config
}

func newLiveServer(cfg config.Config) (*liveServer, error) {
	opts := options.Client().ApplyURI(cfg.DB.Addr)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.DB.DialTimeout))
	defer cancel()
	db, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("could not connect to mongo: " + err.Error())
	}
	logEntry.Info("connection to mongodb: ", cfg.DB.Addr)
	pingCtx, pingCancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.DB.DialTimeout))
	defer pingCancel()
	if err = db.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("ping mongo error: %v", err)
	}
	logEntry.Info("connected to mongodb: ", cfg.DB.Addr)

	return &liveServer{db, cfg}, nil
}

func (s *liveServer) Get(ctx context.Context, in *pb.GetLiveRequest) (*pb.GetLiveResponse, error) {
	var result pb.GetLiveResponse
	collection := s.db.Database(s.cfg.DB.Name).Collection(s.cfg.DB.Collection)
	if err := collection.FindOne(context.Background(), bson.M{"match_id": in.GetMatchId()}).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return &result, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &result, nil
}

func (s *liveServer) GetTotal(ctx context.Context, in *pb.GetTotalLiveRequest) (*pb.GetTotalLiveResponse, error) {
	collection := s.db.Database(s.cfg.DB.Name).Collection(s.cfg.DB.Collection)
	count, err := collection.CountDocuments(context.Background(), bson.M{}, options.Count())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if count < 0 {
		count = 0
	}

	return &pb.GetTotalLiveResponse{
		Total: uint64(count),
	}, nil
}

func (s *liveServer) FindMultiple(ctx context.Context, in *pb.FindMultipleLiveRequest) (*pb.FindMultipleLiveResponse, error) {
	var matches []*pb.GetLiveResponse
	collection := s.db.Database(s.cfg.DB.Name).Collection(s.cfg.DB.Collection)
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"starttime": -1})
	findOptions.SetLimit(int64(in.GetTotal()))
	filter := bson.M{"$or": bson.A{
		bson.M{"team1.averagerating": bson.D{
			{"$gt", in.GetRatingOver()},
			{"$lt", in.GetRatingUnder()},
		}},
		bson.M{"team2.averagerating": bson.D{
			{"$gt", in.GetRatingOver()},
			{"$lt", in.GetRatingUnder()},
		}},
	}}
	cur, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer cur.Close(context.Background())
	if err := cur.All(context.Background(), &matches); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.FindMultipleLiveResponse{
		Matches: matches,
	}, nil
}

func (s *liveServer) FindUser(ctx context.Context, in *pb.FindUserLiveRequest) (*pb.FindUserLiveResponse, error) {
	var result pb.GetLiveResponse
	collection := s.db.Database(s.cfg.DB.Name).Collection(s.cfg.DB.Collection)
	filter := bson.M{"$or": bson.A{
		bson.M{"team1.users": bson.M{"$in": bson.A{in.GetUserId()}}},
		bson.M{"team2.users": bson.M{"$in": bson.A{in.GetUserId()}}},
	},
	}
	if err := collection.FindOne(context.Background(), filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.FindUserLiveResponse{
				UserId:  in.GetUserId(),
				MatchId: 0,
			}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.FindUserLiveResponse{
		UserId:  in.GetUserId(),
		MatchId: result.GetMatchId(),
	}, nil
}

func (s *liveServer) Add(ctx context.Context, in *pb.AddLiveRequest) (*pb.AddLiveResponse, error) {
	collection := s.db.Database(s.cfg.DB.Name).Collection(s.cfg.DB.Collection)
	filter := bson.M{"match_id": in.GetMatchId()}
	update := bson.M{"$set": in}
	opts := options.Update().SetUpsert(true)
	if _, err := collection.UpdateOne(context.Background(), filter, update, opts); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.AddLiveResponse{}, nil
}

func (s *liveServer) Remove(ctx context.Context, in *pb.RemoveLiveRequest) (*pb.RemoveLiveResponse, error) {
	collection := s.db.Database(s.cfg.DB.Name).Collection(s.cfg.DB.Collection)
	filter := bson.M{"match_id": in.GetMatchId()}
	if _, err := collection.DeleteOne(context.Background(), filter); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RemoveLiveResponse{}, nil
}

func (s *liveServer) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	s.db.Disconnect(ctx)
}
