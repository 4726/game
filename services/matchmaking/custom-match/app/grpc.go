package app

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/4726/game/services/matchmaking/custom-match/config"
	"github.com/4726/game/services/matchmaking/custom-match/pb"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type customMatchServer struct {
	db     *gorm.DB
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
	db, err := gorm.Open("postgres", "user=postgres password=postgres")
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Group{}, &User{}).Error; err != nil {
		return nil, err
	}

	db.Model(&User{}).AddForeignKey("group_id", "groups(id)", "CASCADE", "CASCADE")

	s := &customMatchServer{db, cfg, map[int64][]chan PubSubMessage{}}

	go s.handleNotifications("default")

	return s, nil
}

func (s *customMatchServer) Add(in *pb.AddCustomMatchRequest, outStream pb.CustomMatch_AddServer) error {
	g := Group{
		Leader:   in.GetLeader(),
		Name:     in.GetName(),
		Password: in.GetPassword(),
		MaxUsers: in.GetMaxUsers(),
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&g).Error; err != nil {
			return err
		}

		u := User{
			UserID:  in.GetLeader(),
			GroupID: g.ID,
		}

		return tx.Create(&u).Error
	})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	out := &pb.AddCustomMatchResponse{
		Users:    []uint64{},
		MaxUsers: g.MaxUsers,
		Leader:   g.Leader,
	}

	if err := outStream.Send(out); err != nil {
		return err
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
	res := s.db.Delete(&Group{}, "id = ? AND leader = ?", in.GetGroupId(), in.GetUserId())
	if res.Error != nil {
		return nil, status.Error(codes.Internal, res.Error.Error())
	}

	if res.RowsAffected != 1 {
		return nil, errors.New("todo")
	}

	if err := s.sendNotification(in.GetGroupId(), false, false); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteCustomMatchResponse{}, nil

}

func (s *customMatchServer) GetAll(ctx context.Context, in *pb.GetAllCustomMatchRequest) (*pb.GetAllCustomMatchResponse, error) {
	var groups []Group
	res := s.db.Preload("Users").Find(&groups)
	if res.Error != nil {
		return nil, status.Error(codes.Internal, res.Error.Error())
	}

	var pbGroups []*pb.CustomMatchGroup

	for _, v := range groups {
		pbGroup := &pb.CustomMatchGroup{
			GroupId:          v.ID,
			Name:             v.Name,
			Leader:           v.Leader,
			PasswordRequired: v.Password != "",
			MaxUsers:         v.MaxUsers,
			TotalUsers:       uint32(len(v.Users)),
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
					group[i] = group[len(group)-1]
					group = group[:len(group)-1]
				}
			}
			s.groups[in.GetGroupId()] = group
		}
	}()

	u := User{
		UserID:  in.GetUserId(),
		GroupID: in.GetGroupId(),
	}
	if err := s.db.Create(&u).Error; err != nil {
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
	u := User{
		UserID: in.GetUserId(),
	}
	if err := s.db.Delete(&u).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	groupID := u.GroupID

	if err := s.sendNotification(groupID, false, false); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.LeaveCustomMatchResponse{}, nil
}

func (s *customMatchServer) Start(ctx context.Context, in *pb.StartCustomMatchRequest) (*pb.StartCustomMatchResponse, error) {
	var userIDs []uint64

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var deletedGroup Group

		res := tx.First(&deletedGroup, "id = ? AND leader = ?", in.GetGroupId(), in.GetUserId())
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			//user is not leader or group does not exist
			return errors.New("todo")
		}

		for _, v := range deletedGroup.Users {
			userIDs = append(userIDs, v.UserID)
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

		res = tx.Exec("NOTIFY default, '%v'", string(b))
		if res.Error != nil {
			return res.Error
		}

		g := Group{
			ID: in.GetGroupId(),
		}
		return tx.Delete(&g).Error
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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
	return s.db.Transaction(func(tx *gorm.DB) error {
		var g Group
		res := tx.First(&g, "id = ?", groupID)
		if res.Error != nil {
			return res.Error
		}

		var userIDs []uint64
		for _, v := range g.Users {
			userIDs = append(userIDs, v.UserID)
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

		return tx.Exec("NOTIFY default, '%v'", string(b)).Error
	})
}

func (s *customMatchServer) handleNotifications(channelName string) {
	ln := pq.NewListener("name", time.Minute, time.Minute, func(event pq.ListenerEventType, err error) {})
	defer ln.Close()

	ln.Listen(channelName)

	ch := ln.NotificationChannel()
	for {
		notif, ok := <-ch
		if !ok {
			return
		}

		var psm PubSubMessage
		if err := json.Unmarshal([]byte(notif.Extra), &psm); err != nil {
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
