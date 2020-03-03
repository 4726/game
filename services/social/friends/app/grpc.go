package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/4726/game/services/social/friends/config"
	"github.com/4726/game/services/social/friends/pb"
	"github.com/cenkalti/backoff"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errCannotAddSelf           = errors.New("cannot add self")
	errCannotSendFriendRequest = errors.New("cannot send friend request")
	errNotFriend               = errors.New("not friends")
	errNoRequest               = errors.New("no request")
)

type friendsServer struct {
	db  *gorm.DB
	cfg config.Config
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
		db:  db,
		cfg: cfg,
	}, nil
}

func (s *friendsServer) Add(ctx context.Context, in *pb.AddFriendRequest) (*pb.AddFriendResponse, error) {
	if in.GetUserId() == in.GetFriendId() {
		return nil, status.Error(codes.FailedPrecondition, errCannotAddSelf.Error())
	}

	request := Request{
		From:     in.GetUserId(),
		To:       in.GetFriendId(),
		Accepted: false,
	}
	res := s.db.Exec(`INSERT INTO requests ("from", "to", "accepted") 
		SELECT ?, ?, ? FROM requests 
		WHERE ("from" = ? AND "to" = ?) OR ("from" = ? AND "to" = ?) 
		HAVING count(*) < 1;`,
		request.From, request.To, false, request.From, request.To, request.To, request.From)
	if res.Error != nil {
		return nil, status.Error(codes.Internal, res.Error.Error())
	}

	if res.RowsAffected < 1 {
		return nil, status.Error(codes.FailedPrecondition, errCannotSendFriendRequest.Error())
	}

	return &pb.AddFriendResponse{}, nil
}

func (s *friendsServer) Delete(ctx context.Context, in *pb.DeleteFriendRequest) (*pb.DeleteFriendResponse, error) {
	res := s.db.Where(`("from" = ? AND "to" = ? AND accepted = ?) OR ("from" = ? AND "to" = ? AND accepted = ?)`, in.GetUserId(), in.GetFriendId(), true, in.GetFriendId(), in.GetUserId(), true).Delete(Request{})
	if res.Error != nil {
		return nil, status.Error(codes.Internal, res.Error.Error())
	}
	if res.RowsAffected < 1 {
		return nil, status.Error(codes.FailedPrecondition, errNotFriend.Error())
	}

	return &pb.DeleteFriendResponse{}, nil
}

func (s *friendsServer) Get(ctx context.Context, in *pb.GetFriendRequest) (*pb.GetFriendResponse, error) {
	var requests []Request
	if err := s.db.Model(Request{}).Where(`("from" = ? AND "accepted" = ?) OR ("to" = ? AND "accepted" = ?)`, in.GetUserId(), true, in.GetUserId(), true).Find(&requests).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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
	if err := s.db.Model(Request{}).Where(`"to" = ? AND "accepted" = ?`, in.GetUserId(), false).Find(&requests).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var userIDs []uint64
	for _, v := range requests {
		userIDs = append(userIDs, v.From)
	}

	return &pb.GetRequestsFriendResponse{
		Requests: userIDs,
	}, nil
}

func (s *friendsServer) Accept(ctx context.Context, in *pb.AcceptFriendRequest) (*pb.AcceptFriendResponse, error) {
	res := s.db.Model(Request{}).Where(`"from" = ? AND "to" = ? AND accepted = ?`, in.GetFriendId(), in.GetUserId(), false).Update("accepted", true)
	if res.Error != nil {
		return nil, status.Error(codes.Internal, res.Error.Error())
	}
	if res.RowsAffected < 1 {
		return nil, status.Error(codes.FailedPrecondition, errNoRequest.Error())
	}

	return &pb.AcceptFriendResponse{}, nil
}

func (s *friendsServer) Deny(ctx context.Context, in *pb.DenyFriendRequest) (*pb.DenyFriendResponse, error) {
	res := s.db.Where(`"from" = ? AND "to" = ? AND accepted = ?`, in.GetFriendId(), in.GetUserId(), false).Delete(Request{})
	if res.Error != nil {
		return nil, status.Error(codes.Internal, res.Error.Error())
	}
	if res.RowsAffected < 1 {
		return nil, status.Error(codes.FailedPrecondition, errNoRequest.Error())
	}

	return &pb.DenyFriendResponse{}, nil
}
