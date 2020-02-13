package app

import (
	"context"
	"time"

	"github.com/4726/game/services/social/chat/config"
	"github.com/4726/game/services/social/chat/pb"
	"github.com/gocql/gocql"
	"github.com/golang/protobuf/ptypes"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type chatServer struct {
	db  *gocql.Session
	cfg config.Config
}

type Message struct {
	From, To uint64
	Message  string
	Time     time.Time
}

func newChatServer(c config.Config) (*chatServer, error) {
	cluster := gocql.NewCluster("192.168.1.1")
	cluster.Keyspace = "chat"
	cluster.ConnectTimeout = time.Second * time.Duration(c.Cassandra.DialTimeout)
	cluster.Port = c.Cassandra.Port
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return &chatServer{session, c}, nil
}

func (s *chatServer) Send(ctx context.Context, in *pb.SendChatRequest) (*pb.SendChatResponse, error) {
	msg := Message{in.GetFrom(), in.GetTo(), in.GetMessage(), time.Now()}
	stmt, names := qb.Insert("table.chat").Columns("from", "to", "message").ToCql()
	q := gocqlx.Query(s.db.Query(stmt), names).BindStruct(msg)

	if err := q.ExecRelease(); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SendChatResponse{}, nil
}

func (s *chatServer) Get(ctx context.Context, in *pb.GetChatRequest) (*pb.GetChatResponse, error) {
	var msgs []Message
	stmt, names := qb.Select("table.chat").Where(qb.Eq("from"), qb.Eq("to")).OrderBy("time", qb.Order(false)).Limit(uint(in.GetTotal())).ToCql()
	q := gocqlx.Query(s.db.Query(stmt), names).BindMap(qb.M{
		"from": in.GetUser1(),
		"to":   in.GetUser2(),
	})
	if err := q.SelectRelease(&msgs); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var pbMsgs []*pb.ChatMessage
	for _, v := range msgs {
		pbTimestamp, err := ptypes.TimestampProto(v.Time)
		if err != nil {
			continue
		}
		pbMsg := &pb.ChatMessage{
			From:    v.From,
			To:      v.To,
			Message: v.Message,
			Time:    pbTimestamp,
		}
		pbMsgs = append(pbMsgs, pbMsg)
	}

	return &pb.GetChatResponse{
		Messages: pbMsgs,
	}, nil
}

//Close gracefully stops the server
func (s *chatServer) Close() {
	s.db.Close()
}
