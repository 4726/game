package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/4726/game/services/social/friends/config"
	"github.com/4726/game/services/social/friends/pb"
	"github.com/cenkalti/backoff"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type friendsServer struct {
	db     *gorm.DB
	cfg    config.Config
}

func newFriendsServer(cfg config.Config) (*friendsServer, error) {
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

	if err := db.AutoMigrate(&Request{}).Error; err != nil {
		return nil, err
	}

	return &friendsServer{
		db:     db,
		cfg:    cfg,
	}, nil
}

func (s *friendsServer) Add(ctx context.Context, in *pb.AddFriendRequest) (*pb.AddFriendResponse, error) {
	request := Request{
		From: in.GetUserId(),
		To: in.GetFriendId(),
		Accepted: false,
	}
	if err := s.db.Create(&request); err != nil {
		return nil, status.Error(codes.Internal, err)
	}

	return &pb.AddFriendResponse{}, nil
}

func (s *friendsServer) Delete(ctx context.Context, in *pb.DeleteFriendRequest) (*pb.DeleteFriendResponse, error) {
	if err := s.db.Where("(from = ? AND to = ?) OR (from = ? AND to = ?)", in.GetUserId(), in.GetFriendId(), in.GetFriendId(), in.GetUserId()).Delete(Request{}); err != nil {
		return nil, status.Error(codes.Internal, err)
	}

	return &pb.DeleteFriendResponse{}, nil
}

func (s *friendsServer) Get(ctx context.Context, in *pb.GetFriendRequest) (*pb.GetFriendResponse, error) {
	var requests []Request
	if err := s.db.Model(Request{}).Where("(from = ? AND accepted = ?) OR (to = ? AND accepted = ?)", in.GetUserId(), true, in.GetUserId(), true).Find(&requests); err != nil {
		return nil, status.Error(codes.Internal, err)
	}

	var userIDs []uint64
	for _, v := range requests {
		if v.From != in.GetUserId() {
			userIDs = append(userIDs, v.From)
		} else {
			userIDs = append(userIDs, v.To)
		}
	}

	return &pb.GetFriendResponse{
		Friends: userIDs,
	}, nil
}

func (s *friendsServer) GetRequests(ctx context.Context, in *pb.GetRequestsFriendRequest) (*pb.GetRequestsFriendResponse, error) {
	var requests []Request
	if err := s.db.Model(Request{}).Where("to = ? AND accepted = ?", in.GetUserId(), false).Find(&requests); err != nil {
		return nil, status.Error(codes.Internal, err)
	}

	var userIDs []uint64
	for _, v := range requests {
		userIDs = append(userIDs, v.To)
	}

	return &pb.GetRequestsFriendResponse{
		Requests: userIds,
	}, nil
}

func (s *friendsServer) Accept(ctx context.Context, in *pb.AcceptFriendRequest) (*pb.AcceptFriendResponse, error) {
	if err := s.db.Model(Request{}).Where("from = ? AND to =?", in.GetFriendId(), in.GetUserId()).Update("accepted", true); err != nil {
		return nil, status.Error(codes.Internal, err)
	}

	return &pb.AcceptFriendResponse{}, nil
}

func (s *friendsServer) Deny(ctx context.Context, in *pb.DenyFriendRequest) (*pb.DenyFriendResponse, error) {
	if err := s.db.Where("(from = ? AND to = ?) OR (from = ? AND to = ?)", in.GetUserId(), in.GetFriendId(), in.GetFriendId(), in.GetUserId()).Delete(Request{}); err != nil {
		return nil, status.Error(codes.Internal, err)
	}

	return &pb.DenyFriendResponse{}, nil
}