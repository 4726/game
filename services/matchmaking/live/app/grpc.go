package app

import (
	"github.com/4726/game/services/matchmaking/live/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type liveServer struct {
	db       *mongo.Client
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

	return uint64(count), nil
}

func (s *liveServer) FindMultiple(ctx context.Context, in *pb.FindMultipleLiveRequest) (*pb.FindMultipleLiveResponse, error) {

}

func (s *liveServer) FindUser(ctx context.Context, in *pb.FindUserLiveRequest) (*pb.FindUserLiveResponse, error) {

}

func (s *liveServer) Add(ctx context.Context, in *pb.AddLiveRequest) (*pb.AddLiveResponse, error) {

}

func (s *liveServer) Remove(ctx context.Context, in *pb.RemoveLiveRequest) (*pb.RemoveLiveResponse, error) {

}
