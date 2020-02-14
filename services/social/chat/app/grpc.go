package app


// CREATE TABLE messages(
// 	messages_id text PRIMARY KEY,
// 	messages_from varint,
// 	messages_to varint,
// 	messages_message text,
// 	messages_users set<varint>,
// 	messages_time bigint
// );

import (
	"context"
	"time"

	"github.com/4726/game/services/social/chat/config"
	"github.com/4726/game/services/social/chat/pb"
	"github.com/gocql/gocql"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
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
	MessagesId               string
	MessagesFrom, MessagesTo uint64
	MessagesMessage          string
	MessagesUsers            []uint64
	MessagesTime             int64
}

func newChatServer(c config.Config) (*chatServer, error) {
	cluster := gocql.NewCluster(c.Cassandra.Host)
	cluster.Keyspace = "chat"
	cluster.ConnectTimeout = time.Second * time.Duration(c.Cassandra.DialTimeout)
	cluster.Port = c.Cassandra.Port
	cluster.Consistency = gocql.Any
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
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	msg := Message{id.String(), in.GetFrom(), in.GetTo(), in.GetMessage(), sortedUsers, time.Now().Unix()}
	stmt, names := qb.Insert("messages").Columns("messages_id", "messages_from", "messages_to", "messages_message", "messages_users", "messages_time").ToCql()
	q := gocqlx.Query(s.db.Query(stmt), names).BindStruct(msg)

	if err := q.ExecRelease(); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SendChatResponse{}, nil
}

func (s *chatServer) Get(ctx context.Context, in *pb.GetChatRequest) (*pb.GetChatResponse, error) {
	var msgs []Message
	stmt, names := qb.Select("messages").Where(qb.Eq("messages.users")).OrderBy("time", qb.Order(false)).Limit(uint(in.GetTotal())).ToCql()
	var usersQuery []uint64
	if in.GetUser1() < in.GetUser2() {
		usersQuery = []uint64{in.GetUser1(), in.GetUser2()}
	} else {
		usersQuery = []uint64{in.GetUser2(), in.GetUser1()}
	}
	q := gocqlx.Query(s.db.Query(stmt), names).BindMap(qb.M{
		"users": usersQuery,
	})
	if err := q.SelectRelease(&msgs); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var pbMsgs []*pb.ChatMessage
	for _, v := range msgs {
		t := time.Unix(v.MessagesTime, 0)
		pbTimestamp, err := ptypes.TimestampProto(t)
		if err != nil {
			continue
		}
		pbMsg := &pb.ChatMessage{
			From:    v.MessagesFrom,
			To:      v.MessagesTo,
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
