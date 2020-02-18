package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/4726/game/services/matchmaking/custom-match/config"
	"github.com/4726/game/services/matchmaking/custom-match/pb"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type customMatchServer struct {
	db     *pg.DB
	cfg    config.Config
	groups map[int64][]chan PubSubMessage
}

type PubSubMessage struct {
	GroupID   int64
	Leader    uint64
	Name      string
	Password  string
	MaxUsers  uint32
	Users     []uint64
	Started   bool
	Disbanded bool
}

func newCustomMatchServer(cfg config.Config) (*customMatchServer, error) {
	db := pg.Connect(&pg.Options{
		User: "postgres",
	})

	if err := db.CreateTable(Group{}, &orm.CreateTableOptions{
		IfNotExists: true,
	}); err != nil {
		return nil, err
	}

	if err := db.CreateTable(User{}, &orm.CreateTableOptions{
		IfNotExists: true,
	}); err != nil {
		return nil, err
	}

	s := &customMatchServer{db, cfg, map[int64][]chan PubSubMessage{}}

	go s.handleNotifications("default")

	return s, nil
}

func (s *customMatchServer) Add(in *pb.AddCustomMatchRequest, outStream pb.CustomMatch_AddServer) error {
	g := Group{
		Leader:     in.GetLeader(),
		Name:       in.GetName(),
		Password:   in.GetPassword(),
		MaxUsers:   in.GetMaxUsers(),
		TotalUsers: 1,
	}

	u := User{
		ID:      in.GetLeader(),
		GroupID: g.ID,
	}

	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		if err := tx.Insert(&u); err != nil {
			return err
		}
		return tx.Insert(&g)
	})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	ch := make(chan PubSubMessage, 1)
	s.groups[g.ID] = []chan PubSubMessage{ch}

	for {
		psm, ok := <-ch
		if !ok {
			return nil
		}

		out := &pb.AddCustomMatchResponse{
			Users:    psm.Users,
			MaxUsers: psm.MaxUsers,
			Leader:   psm.Leader,
		}

		if err := outStream.Send(out); err != nil {
			return err
		}
	}
}

func (s *customMatchServer) Delete(ctx context.Context, in *pb.DeleteCustomMatchRequest) (*pb.DeleteCustomMatchResponse, error) {
	res, err := s.db.Model(&Group{}).Where("id = ? AND leader = ?", in.GetGroupId(), in.GetUserId()).Delete()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if res.RowsAffected() != 1 {
		return nil, errors.New("todo")
	}

	if err := s.sendNotification(in.GetGroupId(), false, false); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteCustomMatchResponse{}, nil

}

func (s *customMatchServer) GetAll(ctx context.Context, in *pb.GetAllCustomMatchRequest) (*pb.GetAllCustomMatchResponse, error) {
	var groups []Group
	_, err := s.db.Query(&groups, "SELECT * FROM groups")
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())

	}

	var pbGroups []*pb.CustomMatchGroup

	for _, v := range groups {
		pbGroup := &pb.CustomMatchGroup{
			GroupId:          v.ID,
			Name:             v.Name,
			Leader:           v.Leader,
			PasswordRequired: v.Password != "",
			MaxUsers:         v.MaxUsers,
			TotalUsers:       v.TotalUsers,
		}
		pbGroups = append(pbGroups, pbGroup)
	}

	return &pb.GetAllCustomMatchResponse{
		Groups: pbGroups,
	}, nil
}

func (s *customMatchServer) Join(in *pb.JoinCustomMatchRequest, outStream pb.CustomMatch_JoinServer) error {
	//do this first because not sure if can get row affected when running the update query
	ch := make(chan PubSubMessage, 1)
	group, ok := s.groups[in.GetGroupId()]
	if !ok {
		s.groups[in.GetGroupId()] = []chan PubSubMessage{ch}
	} else {
		group = append(group, ch)
		s.groups[in.GetGroupId()] = group
	}
	defer func() {
		group, ok := s.groups[in.GetGroupId()]
		if ok {
			for i, v := range group {
				if v == ch {
					group[i] = group[len(group) - 1]
					group = group[:len(group) - 1]
				}
			}
			s.groups[in.GetGroupId()] = group
		}
	}()

	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		u := User{
			ID:      in.GetUserId(),
			GroupID: in.GetGroupId(),
		}
		if err := tx.Insert(&u); err != nil {
			return err
		}

		_, err := tx.Exec("UPDATE groups SET total_users = total_users + 1 WHERE id = ?", in.GetGroupId())
		if err != nil {
			return err
		}

		var g Group
		_, err = tx.QueryOne(&g, "SELECT * FROM groups WHERE id = ?", in.GetGroupId())
		return err
	})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if err := s.sendNotification(in.GetGroupId(), false, false); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for {
		psm, ok := <-ch
		if !ok {
			return nil
		}

		out := &pb.JoinCustomMatchResponse{
			Users:    psm.Users,
			MaxUsers: psm.MaxUsers,
			Leader:   psm.Leader,
		}

		if err := outStream.Send(out); err != nil {
			return err
		}
	}
}

func (s *customMatchServer) Leave(ctx context.Context, in *pb.LeaveCustomMatchRequest) (*pb.LeaveCustomMatchResponse, error) {
	var groupID int64
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		u := User{
			ID: in.GetUserId(),
		}
		if err := tx.Delete(&u); err != nil {
			return err
		}

		_, err := tx.Model(&u).Returning("group_id").Delete()
		if err != nil {
			return err
		}

		groupID = u.GroupID

		_, err = tx.Exec("UPDATE groups SET total_users = total_users - 1 WHERE id = ?", u.GroupID)
		return err
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := s.sendNotification(groupID, false, false); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.LeaveCustomMatchResponse{}, nil
}

func (s *customMatchServer) Start(ctx context.Context, in *pb.StartCustomMatchRequest) (*pb.StartCustomMatchResponse, error) {
	var users []User

	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		var deletedGroup Group

		res, err := tx.QueryOne(&deletedGroup, "SELECT * FROM groups WHERE id = ? AND leader = ?", in.GetGroupId(), in.GetUserId())
		if err != nil {
			return err
		}
		if res.RowsAffected() != 1 {
			//user is not leader or group does not exist
			return errors.New("todo")
		}

		_, err = tx.Query(&users, "SELECT * FROM users WHERE group_id = ?", in.GetGroupId())
		if err != nil {
			return err
		}

		var userIDs []uint64
		for _, v := range users {
			userIDs = append(userIDs, v.ID)
		}
		psm := PubSubMessage{
			GroupID:   in.GetGroupId(),
			Leader:    deletedGroup.Leader,
			Name:      deletedGroup.Name,
			Password:  deletedGroup.Password,
			MaxUsers:  deletedGroup.MaxUsers,
			Users:     userIDs,
			Started:   true,
			Disbanded: false, 
		}
		b, err := json.Marshal(&psm)
		if err != nil {
			return err
		}

		_, err = tx.Exec(fmt.Sprintf("NOTIFY default, '%v'", string(b)))
		if err != nil {
			return err
		}

		g := Group{
			ID: in.GetGroupId(),
		}
		if err := tx.Delete(&g); err != nil {
			return err
		}

		_, err = tx.Model(&User{}).Where("group_id = ?", in.GetGroupId()).Delete()
		return err
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var userIDs []uint64
	for _, v := range users {
		userIDs = append(userIDs, v.ID)
	}

	return &pb.StartCustomMatchResponse{
		Users:   userIDs,
		Leader:  in.GetUserId(),
		GroupId: in.GetGroupId(),
	}, nil
}

func (s *customMatchServer) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *customMatchServer) sendNotification(groupID int64, disbanded, started bool) error {
	return s.db.RunInTransaction(func(tx *pg.Tx) error {
		var g Group
		_, err := tx.QueryOne(&g, "SELECT * FROM groups WHERE id = ?", groupID)
		if err != nil {
			return err
		}

		var users []User
		_, err = tx.Query(&users, "SELECT * FROM users WHERE group_id = ?", groupID)
		if err != nil {
			return err
		}

		var userIDs []uint64
		for _, v := range users {
			userIDs = append(userIDs, v.ID)
		}

		psm := PubSubMessage{
			GroupID:   groupID,
			Leader:    g.Leader,
			Name:      g.Name,
			Password:  g.Password,
			MaxUsers:  g.MaxUsers,
			Users:     userIDs,
			Started:   started,
			Disbanded: disbanded,
		}
		b, err := json.Marshal(&psm)
		if err != nil {
			return err
		}

		_, err = tx.Exec(fmt.Sprintf("NOTIFY default, '%v'", string(b)))
		return err
	})
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

		var psm PubSubMessage
		if err := json.Unmarshal([]byte(notif.Payload), &psm); err != nil {
			continue
		}

		chs, ok := s.groups[psm.GroupID]
		if !ok {
			continue
		}
		for _, v := range chs {
			v <- psm
		}
	}
}
