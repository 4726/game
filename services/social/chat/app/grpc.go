package app

// CREATE TABLE messages(
// 	messages_from varint,
// 	messages_message text,
// 	messages_user1 varint,
// 	messages_user2 varint,
// 	messages_time timeuuid,
// 	PRIMARY KEY ((messages_user1, messages_user2), messages_time)
// ) WITH CLUSTERING ORDER BY (messages_time ASC);

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
	MessagesFrom    uint64
	MessagesMessage string
	MessagesUser1   uint64
	MessagesUser2   uint64
	MessagesTime    gocql.UUID
}

func newChatServer(c config.Config) (*chatServer, error) {
	cluster := gocql.NewCluster(c.Cassandra.Host)
	cluster.Keyspace = "chat"
	cluster.ConnectTimeout = time.Second * time.Duration(c.Cassandra.DialTimeout)
	cluster.Port = c.Cassandra.Port
	cluster.Timeout = time.Second * 5
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return &chatServer{session, c}, nil
}

func (s *chatServer) Send(ctx context.Context, in *pb.SendChatRequest) (*pb.SendChatResponse, error) {
	var sortedUsers []uint64
	if in.GetFrom() < in.GetTo() {
		sortedUsers = []uint64{in.GetFrom(), in.GetTo()}
	} else {
		sortedUsers = []uint64{in.GetTo(), in.GetFrom()}
	}
	msg := Message{in.GetFrom(), in.GetMessage(), sortedUsers[0], sortedUsers[1], gocql.TimeUUID()}
	stmt, names := qb.Insert("messages").Columns("messages_from", "messages_message", "messages_user1", "messages_user2", "messages_time").ToCql()
	q := gocqlx.Query(s.db.Query(stmt), names).BindStruct(msg)

	if err := q.ExecRelease(); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SendChatResponse{}, nil
}

func (s *chatServer) Get(ctx context.Context, in *pb.GetChatRequest) (*pb.GetChatResponse, error) {
	var msgs []Message
	stmt, names := qb.Select("messages").Where(qb.Eq("messages_user1")).Where(qb.Eq("messages_user2")).OrderBy("messages_time", qb.Order(false)).Limit(uint(in.GetTotal())).ToCql()
	var usersQuery []uint64
	if in.GetUser1() < in.GetUser2() {
		usersQuery = []uint64{in.GetUser1(), in.GetUser2()}
	} else {
		usersQuery = []uint64{in.GetUser2(), in.GetUser1()}
	}
	q := gocqlx.Query(s.db.Query(stmt), names).BindMap(qb.M{
		"messages_user1": usersQuery[0],
		"messages_user2": usersQuery[1],
	})
	if err := q.SelectRelease(&msgs); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var pbMsgs []*pb.ChatMessage
	for _, v := range msgs {
		t := v.MessagesTime.Time()
		pbTimestamp, err := ptypes.TimestampProto(t)
		if err != nil {
			continue
		}
		var to uint64
		if v.MessagesUser1 == v.MessagesFrom {
			to = v.MessagesUser2
		} else {
			to = v.MessagesUser1
		}
		pbMsg := &pb.ChatMessage{
			From:    v.MessagesFrom,
			To:      to,
			Message: v.MessagesMessage,
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
