package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/4726/game/services/matchmaking/custom-match/config"
	"github.com/4726/game/services/matchmaking/custom-match/pb"
	"github.com/cenkalti/backoff"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type customMatchServer struct {
	db     *gorm.DB
	cfg    config.Config
	groups map[int64][]chan PubSubMessage
	sync.Mutex
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

var (
	ErrNoLeaderPrivileges = errors.New("user is not leader of group")
	ErrGroupFull          = errors.New("group full")
	ErrGroupDoesNotExist  = errors.New("group does not exist or password is invalid")
	ErrUserAlreadyInGroup = errors.New("user already in group")
	ErrUserNotInGroup     = errors.New("user not in group")
)

func newCustomMatchServer(cfg config.Config) (*customMatchServer, error) {
	dbParams := fmt.Sprintf("user=%v password=%v dbname=%v host=%v port=%v",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.Host,
		cfg.DB.Port,
	)

	var db *gorm.DB

	op := func() error {
		logEntry.Infof("connecting to postgres: %v:%v", cfg.DB.Host, cfg.DB.Port)
		var err error
		db, err = gorm.Open("postgres", dbParams)
		if err != nil {
			logEntry.Warn("could not connect to postgres, retrying")
		}
		return err
	}

	if err := backoff.Retry(op, backoff.NewExponentialBackOff()); err != nil {
		logEntry.Error("could not connect to postgres, max retries reached")
		return nil, fmt.Errorf("could not connect to redis: %v", err)
	}
	logEntry.Infof("connected to postgres: %v:%v", cfg.DB.Host, cfg.DB.Port)

	if err := db.AutoMigrate(&Group{}, &User{}).Error; err != nil {
		return nil, err
	}

	db.Model(&User{}).AddForeignKey("group_id", "groups(id)", "CASCADE", "CASCADE")

	s := &customMatchServer{
		db:     db,
		cfg:    cfg,
		groups: map[int64][]chan PubSubMessage{},
	}

	go s.handleNotifications("groups")

	return s, nil
}

func (s *customMatchServer) Add(in *pb.AddCustomMatchRequest, outStream pb.CustomMatch_AddServer) error {
	g := Group{
		Leader:   in.GetLeader(),
		Name:     in.GetName(),
		Password: in.GetPassword(),
		MaxUsers: in.GetMaxUsers(),
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&g).Error; err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		u := User{
			UserID:  in.GetLeader(),
			GroupID: g.ID,
		}

		return tx.Create(&u).Error
	}); err != nil {
		if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
			return status.Error(codes.FailedPrecondition, ErrUserAlreadyInGroup.Error())
		}
		return status.Error(codes.Internal, err.Error())
	}

	resp := &pb.AddCustomMatchResponse{
		Users:    []uint64{in.GetLeader()},
		MaxUsers: g.MaxUsers,
		Leader:   g.Leader,
		GroupId:  g.ID,
	}

	if err := outStream.Send(resp); err != nil {
		return err
	}

	ch := make(chan PubSubMessage, 1)
	s.Lock()
	s.groups[g.ID] = []chan PubSubMessage{ch}
	s.Unlock()

	for {
		psm, ok := <-ch
		if !ok {
			return nil
		}

		out := &pb.AddCustomMatchResponse{
			Users:    psm.Users,
			MaxUsers: psm.MaxUsers,
			Leader:   psm.Leader,
			GroupId:  psm.GroupID,
			Started:  psm.Started,
		}

		if err := outStream.Send(out); err != nil {
			return err
		}
	}
}

func (s *customMatchServer) Delete(ctx context.Context, in *pb.DeleteCustomMatchRequest) (*pb.DeleteCustomMatchResponse, error) {
	res := s.db.Where("id = ? AND leader = ?", in.GetGroupId(), in.GetUserId()).Delete(Group{})
	if res.Error != nil {
		return nil, status.Error(codes.Internal, res.Error.Error())
	}

	if res.RowsAffected != 1 {
		return nil, status.Error(codes.FailedPrecondition, ErrNoLeaderPrivileges.Error())
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
	s.Lock()
	group, ok := s.groups[in.GetGroupId()]
	if !ok {
		s.groups[in.GetGroupId()] = []chan PubSubMessage{ch}
	} else {
		group = append(group, ch)
		s.groups[in.GetGroupId()] = group
	}
	s.Unlock()
	defer func() {
		s.Lock()
		defer s.Unlock()
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

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		var group Group
		res := tx.First(&group, "id = ? AND password = ?", in.GetGroupId(), in.GetGroupPassword())
		if res.Error != nil {
			if res.Error == gorm.ErrRecordNotFound {
				return status.Error(codes.FailedPrecondition, ErrGroupDoesNotExist.Error())
			}
			return status.Error(codes.Internal, res.Error.Error())
		}

		res = tx.Exec("INSERT INTO users (user_id, group_id) SELECT ?, ? FROM users WHERE group_id = ? HAVING count(*) < ?;", in.GetUserId(), in.GetGroupId(), in.GetGroupId(), group.MaxUsers)
		if res.Error != nil {
			return status.Error(codes.Internal, res.Error.Error())
		}

		if res.RowsAffected != 1 {
			return status.Error(codes.FailedPrecondition, ErrGroupFull.Error())
		}
		return nil
	}); err != nil {
		return err
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
			GroupId:  psm.GroupID,
			Started:  psm.Started,
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
	res := s.db.Delete(&u)
	if res.Error != nil {
		return nil, status.Error(codes.Internal, res.Error.Error())
	}
	if res.RowsAffected != 1 {
		return nil, status.Error(codes.FailedPrecondition, ErrUserNotInGroup.Error())
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

		res := tx.Preload("Users").First(&deletedGroup, "id = ? AND leader = ?", in.GetGroupId(), in.GetUserId())
		if res.Error != nil {
			if res.Error == gorm.ErrRecordNotFound {
				return status.Error(codes.FailedPrecondition, ErrNoLeaderPrivileges.Error())
			}
			return status.Error(codes.Internal, res.Error.Error())
		}
		if res.RowsAffected != 1 {
			return status.Error(codes.FailedPrecondition, ErrNoLeaderPrivileges.Error())
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
			return status.Error(codes.Internal, err.Error())
		}

		if err := tx.Exec(fmt.Sprintf("NOTIFY groups, '%v'", string(b))).Error; err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		g := Group{
			ID: in.GetGroupId(),
		}
		if err := tx.Delete(&g).Error; err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		return nil
	})
	if err != nil {
		return nil, err
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
		res := tx.Preload("Users").First(&g, "id = ?", groupID)
		if res.Error != nil {
			if res.Error == gorm.ErrRecordNotFound {
				psm := PubSubMessage{
					GroupID:   groupID,
					Started:   started,
					Disbanded: disbanded,
				}
				b, err := json.Marshal(&psm)
				if err != nil {
					return err
				}

				return tx.Exec(fmt.Sprintf("NOTIFY groups, '%v'", string(b))).Error
			}
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

		return tx.Exec(fmt.Sprintf("NOTIFY groups, '%v'", string(b))).Error
	})
}

func (s *customMatchServer) handleNotifications(channelName string) {
	ln := pq.NewListener("user=postgres password=postgres", time.Second*10, time.Minute, func(event pq.ListenerEventType, err error) {})
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

		s.Lock()
		chs, ok := s.groups[psm.GroupID]
		if !ok {
			continue
		}
		for _, v := range chs {
			v <- psm
		}
		s.Unlock()

	}
}
