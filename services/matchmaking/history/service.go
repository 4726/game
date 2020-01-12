package main

import (
	"context"
	"fmt"

	"github.com/4726/game/services/matchmaking/history/pb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nsqio/go-nsq"
)

type Service struct {
	consumer *nsq.Consumer
}

type Config struct {
	DB DBConfig
}

type DBConfig struct {
	Username, Password, Name, Address string
}

func NewService(c Config) (*Service, error) {

	consumer, err := nsq.NewConsumer("matches", "db", nsq.NewConfig())
	if err != nil {
		return nil, err
	}
	consumer.AddHandler(&nsqMessageHandler{db})
	err := consumer.ConnectToNSQD()
	if err != nil {
		return nil, err
	}

	return &Service{db, consumer}, nil
}

func (s *Service) Get(ctx context.Context, in *pb.GetHistoryRequest) (*pb.GetHistoryResponse, error) {
	matches := &pb.MatchHistoryInfo{}
	return &pb.GetHistoryResponse{
		Match: matches
	}, nil
}

//a bit difficult because have to search inside of array, might have to exectue raw query
func (s *Service) GetUser(ctx context.Context, in *pb.GetUserHistoryRequest) (*pb.GetUserHistoryResponse, error) {
	return nil, nil
}

func (s *Service) Close() {
	s.consumer.Stop()
	<-s.consumer.StopChan
}
