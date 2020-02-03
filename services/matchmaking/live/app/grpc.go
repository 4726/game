package app

import (
	"context"

	"github.com/4726/game/services/matchmaking/live/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type liveServer struct {
	db *mongo.Client
}

func newLiveServer() (*liveServer, error) {
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
	
	return &liveServer{db}, nil
}

func (s *liveServer) Get(ctx context.Context, in *pb.GetLiveRequest) (*pb.GetLiveResponse, error) {
	var result pb.GetLiveResponse
	collection := s.db.Database("db").Collection("collection")
	if err := collection.FindOne(context.Background(), bson.M{"match_id": in.GetMatchId()}).Decode(&result); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &result, nil
}

func (s *liveServer) GetTotal(ctx context.Context, in *pb.GetTotalLiveRequest) (*pb.GetTotalLiveResponse, error) {
	collection := s.db.Database("db").Collection("collection")
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
	collection := s.db.Database("db").Collection("collection")
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"start_time": -1})
	findOptions.SetLimit(int64(in.GetTotal()))
	cur, err := collection.Find(context.Background(), bson.M{}, findOptions)
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
	collection := s.db.Database("db").Collection("collection")
	filter := bson.M{"$or": bson.A{
		bson.M{"team1.users": in.GetUserId()},
		bson.M{"team2.users": in.GetUserId()},
	},
	}
	if err := collection.FindOne(context.Background(), filter).Decode(&result); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.FindUserLiveResponse{
		UserId: in.GetUserId(),
		MatchId: result.GetMatchId(),
	}, nil
}

func (s *liveServer) Add(ctx context.Context, in *pb.AddLiveRequest) (*pb.AddLiveResponse, error) {
	collection := s.db.Database("db").Collection("collection")
	filter := bson.M{"match_id": in.GetMatchId()}
	update := bson.M{"$set": in}
	opts := options.Update().SetUpsert(true)
	if _, err := collection.UpdateOne(context.Background(), filter, update, opts); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.AddLiveResponse{}, nil
}

func (s *liveServer) Remove(ctx context.Context, in *pb.RemoveLiveRequest) (*pb.RemoveLiveResponse, error) {
	collection := s.db.Database("db").Collection("collection")
	filter := bson.M{"match_id": in.GetMatchId()}
	if _, err := collection.DeleteOne(context.Background(), filter); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RemoveLiveResponse{}, nil
}
