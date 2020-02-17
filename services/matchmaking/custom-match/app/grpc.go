package app

import (
	"context"
	"fmt"
	"time"

	"github.com/4726/game/services/matchmaking/custom-match/config"
	"github.com/4726/game/services/matchmaking/custom-match/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type customMatchServer struct {
	db  *mongo.Client
	cfg config.Config
	groups map[int64][]chan Group
}

func newCustomMatchServer(cfg config.Config) (*customMatchServer, error) {
    db := pg.Connect(&pg.Options{
        User: "postgres",
	})
	
	if err := db.CreateTable(Group{}, orm.CreateTableOptions{
		IfNotExists: true.
	}); err != nil {
		return nil, err
	}

	s := &customMatchServer{db, cfg}

	go s.handleNotifications("default")
	
	return s, nil
}

func (s *customMatchServer) Add(in *pb.AddCustomMatchRequest, outStream pb.CustomMatch_AddServer) error {
	g := Group{
		Leader: in.GetLeader(),
		Name: in.GetName(),
		Password: in.GetPassword(),
		MaxUsers: in.GetMaxUsers(),
		Users: []uint64{}
	}
	if err := s.db.Insert(&g); err != nil {
		return err
	}

	ch := make(chan Group, 1)
	s.groups[g.ID] = []chan Group{ch}

	for {
		updatedGroup, ok := <-ch
		if !ok {
			return nil
		}

		out := &pb.AddCustomMatchResponse{
			Users: updatedGroup.Users,
			MaxUsers: updatedGroup.MaxUsers,
			Leader: updatedGroup.Leader,
		}

		if err := outStream.Send(out); err != nil {
			return err
		}
	}
}

func (s *customMatchServer) Delete(ctx context.Context, in *pb.DeleteCustomMatchRequest) (*pb.DeleteCustomMatchResponse, error) {

}

func (s *customMatchServer) GetAll(ctx context.Context, in *pb.GetAllCustomMatchRequest) (*pb.GetAllCustomMatchResponse, error) {

}

func (s *customMatchServer) Join(in *pb.JoinCustomMatchRequest, outStream *pb.CustomMatch_JoinServer) error {
	//do this first because not sure if can get row affected when running the update query
	ch := make(chan Group, 1)
	group, ok := s.groups[in.GetGroupId()]
	if !ok {
		s.groups[in.GetGroupId()] = []chan Group{ch}
	} else {
		group = append(group, ch)
		s.groups[in.GetGroupId()] = group
	}
	defer func() {
		group, ok := s.groups[in.GetGroupId()]
		if ok {
			delete(group, ch)
			s.groups[in.GetGroupId()] = group
		}
	}()

	_, err := s.db.Exec("UPDATE groups SET users = array_append(users, ?) WHERE id = ?", in.GetUserId(), in.GetGroupId())
	if err != nil {
		return status.Code(codes.Internal, err)
	}

	for {
		updatedGroup, ok := <-ch
		if !ok {
			return nil
		}

		out := &pb.JoinCustomMatchResponse{
			Users: updatedGroup.Users,
			MaxUsers: updatedGroup.MaxUsers,
			Leader: updatedGroup.Leader,
		}

		if err := outStream.Send(out); err != nil {
			return err
		}
	}
}

func (s *customMatchServer) Leave(ctx context.Context, in *pb.LeaveCustomMatchRequest) (*pb.LeaveCustomMatchResponse, error) {

}

func (s *customMatchServer) Start(ctx context.Context, in *pb.StartCustomMatchRequest) (*pb.StartCustomMatchRequest, error) {

}

func (s *customMatchServer) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *customMatchServer) handleNotifications(channelName string) {
	ln := s.db.Listen(channelName)
	defer ln.Close()

	ch := ln.Channel()
	for {
		notif, ok := <-ch
		if !ok {
			return
		}

		var g Group
		if err := json.Unmarshal([]byte(notif.Payload), &g); err != nil {
			continue
		}
		
		chs, ok := s.groups[g.ID]
		if !ok {
			continue
		}
		for _, v := range chs {
			v <- g
		}
	}
}