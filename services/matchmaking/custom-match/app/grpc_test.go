package app

import (
	"context"
	"testing"
	"time"

	"sync"

	"github.com/4726/game/services/matchmaking/custom-match/config"
	"github.com/4726/game/services/matchmaking/custom-match/pb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type test struct {
	c       pb.CustomMatchClient
	service *Service
}

func newTest(t testing.TB, conf ...config.Config) *test {
	var cfg config.Config
	if len(conf) < 1 {
		cfg = config.Config{
			Port:    14000,
			Metrics: config.MetricsConfig{14001, "/metrics"},
		}
	} else {
		cfg = conf[0]
		cfg.Port = 14000
	}

	service, err := NewService(cfg)
	assert.NoError(t, err)

	go service.Run()
	time.Sleep(time.Second * 2)

	conn, err := grpc.Dial("127.0.0.1:14000", grpc.WithInsecure())
	assert.NoError(t, err)
	c := pb.NewCustomMatchClient(conn)

	service.s.db.Exec("TRUNCATE users, groups")
	service.s.db.Exec("ALTER SEQUENCE groups_id_seq RESTART WITH 1")

	return &test{c, service}
}

func (te *test) teardown() {
	te.service.Close()
}

func (te *test) add(t testing.TB, in *pb.AddCustomMatchRequest) ([]*pb.AddCustomMatchResponse, error) {
	var resps []*pb.AddCustomMatchResponse
	outStream, err := te.c.Add(context.Background(), in)
	assert.NoError(t, err)
	ch := make(chan map[string]interface{}, 1)

	go func() {
		for {
			resp, err := outStream.Recv()
			ch <- map[string]interface{}{"resp": resp, "err": err}
		}
	}()

	for {
		select {
		case <-time.After(time.Second * 5):
			return resps, nil
		case msg := <-ch:
			if err, ok := msg["err"].(error); ok {
				return resps, err
			}
			resps = append(resps, msg["resp"].(*pb.AddCustomMatchResponse))
		}
	}
}

func (te *test) join(t testing.TB, userID uint64, groupID int64, pass string) ([]*pb.JoinCustomMatchResponse, error) {
	var resps []*pb.JoinCustomMatchResponse
	joinIn := &pb.JoinCustomMatchRequest{
		UserId:        userID,
		GroupId:       groupID,
		GroupPassword: pass,
	}
	joinStream, err := te.c.Join(context.Background(), joinIn)
	assert.NoError(t, err)
	ch := make(chan map[string]interface{}, 1)

	go func() {
		for {
			resp, err := joinStream.Recv()
			ch <- map[string]interface{}{"resp": resp, "err": err}
		}
	}()

	for {
		select {
		case <-time.After(time.Second * 5):
			return resps, nil
		case msg := <-ch:
			if err, ok := msg["err"].(error); ok {
				return resps, err
			}
			resps = append(resps, msg["resp"].(*pb.JoinCustomMatchResponse))
		}
	}
}

func TestServiceAddLeaderAlreadyHasGroup(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.AddCustomMatchRequest{
		Leader:   1,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}
	te.add(t, in)

	_, err := te.add(t, in)
	assert.Error(t, err)
}

func TestServiceAddLeaderAlreadyInGroup(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddCustomMatchRequest{
		Leader:   1,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}
	resps, _ := te.add(t, addIn)

	te.join(t, 2, resps[0].GetGroupId(), "")

	in := &pb.AddCustomMatchRequest{
		Leader:   2,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}
	_, err := te.add(t, in)
	assert.Error(t, err)
}

func TestServiceAddStart(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.AddCustomMatchRequest{
		Leader:   1,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}
	var wg sync.WaitGroup
	var resps []*pb.AddCustomMatchResponse
	var err error
	wg.Add(1)
	go func() {
		resps, err = te.add(t, in)
		wg.Done()
	}()
	time.Sleep(time.Second * 2)

	r1 := &pb.AddCustomMatchResponse{
		Users:    []uint64{1},
		MaxUsers: in.GetMaxUsers(),
		Leader:   in.GetLeader(),
		GroupId:  1,
	}

	expectedResps := []*pb.AddCustomMatchResponse{r1}
	for i := 2; i < 11; i++ {
		go te.join(t, uint64(i), 1, "")

		var expectedUsers []uint64
		for j := 1; j < i+1; j++ {
			expectedUsers = append(expectedUsers, uint64(j))
		}
		r := &pb.AddCustomMatchResponse{
			Users:    expectedUsers,
			MaxUsers: in.GetMaxUsers(),
			Leader:   in.GetLeader(),
			GroupId:  1,
		}
		expectedResps = append(expectedResps, r)
		time.Sleep(time.Second * 2)
	}

	// need to call start

	wg.Wait()
	assert.NoError(t, err)
	assert.Equal(t, expectedResps, resps)
}

func TestServiceAdd(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.AddCustomMatchRequest{
		Leader:   1,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}
	resps, err := te.add(t, in)
	assert.NoError(t, err)
	r1 := &pb.AddCustomMatchResponse{
		Users:    []uint64{1},
		MaxUsers: in.GetMaxUsers(),
		Leader:   in.GetLeader(),
		GroupId:  1,
	}

	assert.Equal(t, []*pb.AddCustomMatchResponse{r1}, resps)
}

func TestServiceDeleteDoesNotExist(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.DeleteCustomMatchRequest{
		UserId:  1,
		GroupId: 1,
	}
	_, err := te.c.Delete(context.Background(), in)
	assert.Error(t, err)
}

func TestServiceDeleteNotLeader(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddCustomMatchRequest{
		Leader:   1,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}
	resps, _ := te.add(t, addIn)

	in := &pb.DeleteCustomMatchRequest{
		UserId:  2,
		GroupId: resps[0].GetGroupId(),
	}
	_, err := te.c.Delete(context.Background(), in)
	assert.Error(t, err)
}

func TestServiceDelete(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddCustomMatchRequest{
		Leader:   1,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}
	resps, _ := te.add(t, addIn)

	in := &pb.DeleteCustomMatchRequest{
		UserId:  1,
		GroupId: resps[0].GetGroupId(),
	}
	_, err := te.c.Delete(context.Background(), in)
	assert.NoError(t, err)
}

func TestServiceGetAllNone(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.GetAllCustomMatchRequest{Total: 100, Skip: 0}
	resp, err := te.c.GetAll(context.Background(), in)
	assert.NoError(t, err)
	assert.Len(t, resp.Groups, 0)
}

func TestServiceGetAll(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddCustomMatchRequest{
		Leader:   1,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}
	resps, _ := te.add(t, addIn)
	te.join(t, 3, resps[0].GetGroupId(), "")
	addIn = &pb.AddCustomMatchRequest{
		Leader:   2,
		Name:     "game2",
		Password: "asd",
		MaxUsers: 10,
	}
	resps2, _ := te.add(t, addIn)

	in := &pb.GetAllCustomMatchRequest{Total: 100, Skip: 0}
	resp, err := te.c.GetAll(context.Background(), in)
	assert.NoError(t, err)
	g1 := &pb.CustomMatchGroup{
		GroupId:          resps[0].GetGroupId(),
		Name:             "game",
		Leader:           1,
		PasswordRequired: false,
		MaxUsers:         10,
		TotalUsers:       2,
	}
	g2 := &pb.CustomMatchGroup{
		GroupId:          resps2[0].GetGroupId(),
		Name:             "game2",
		Leader:           2,
		PasswordRequired: true,
		MaxUsers:         10,
		TotalUsers:       1,
	}
	expectedResp := &pb.GetAllCustomMatchResponse{
		Groups: []*pb.CustomMatchGroup{g1, g2},
	}
	assert.Equal(t, expectedResp, resp)
}

func TestServiceJoinDoesNotExist(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	_, err := te.join(t, 1, 1, "")
	assert.Error(t, err)
}

// func TestServiceJoinWrongPassword(t *testing.T) {
// 	te := newTest(t)
// 	defer te.teardown()

// 	addIn := &pb.AddCustomMatchRequest{
// 		Leader:   1,
// 		Name:     "game",
// 		Password: "qwe",
// 		MaxUsers: 10,
// 	}
// 	addResps, _ := te.add(t, addIn)

// 	resps, err := te.join(t, 2, addResps[0].GetGroupId(), "asd")
// 	assert.Error(t, err)
// }

func TestServiceJoinPassword(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddCustomMatchRequest{
		Leader:   1,
		Name:     "game",
		Password: "asd",
		MaxUsers: 10,
	}
	addResps, _ := te.add(t, addIn)

	resps, err := te.join(t, 2, addResps[0].GetGroupId(), "asd")
	assert.NoError(t, err)

	r1 := &pb.JoinCustomMatchResponse{
		Users:    []uint64{1, 2},
		MaxUsers: addIn.MaxUsers,
		Leader:   addIn.Leader,
		GroupId:  resps[0].GetGroupId(),
	}
	assert.Equal(t, []*pb.JoinCustomMatchResponse{r1}, resps)
}

// func TestServiceJoinFull(t *testing.T) {
// 	te := newTest(t)
// 	defer te.teardown()

// 	addIn := &pb.AddCustomMatchRequest{
// 		Leader:   1,
// 		Name:     "game",
// 		Password: "",
// 		MaxUsers: 10,
// 	}
// 	addResps, _ := te.add(t, addIn)

// 	for i := 2; i < 11; i++ {
// 		te.join(t, 2, addResps[0].GetGroupId(), "")
// 	}

// 	resps, err := te.join(t, 11, addResps[0].GetGroupId(), "")
// 	assert.Error(t, err)
// }

func TestServiceJoinStart(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.AddCustomMatchRequest{
		Leader:   1,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}

	addResps, _ := te.add(t, in)
	groupID := addResps[0].GetGroupId()

	var wg sync.WaitGroup
	var resps []*pb.JoinCustomMatchResponse
	var err error
	wg.Add(1)
	go func() {
		resps, err = te.join(t, 2, groupID, "")
		wg.Done()
	}()
	time.Sleep(time.Second * 2)

	r1 := &pb.JoinCustomMatchResponse{
		Users:    []uint64{1, 2},
		MaxUsers: in.MaxUsers,
		Leader:   in.Leader,
		GroupId:  groupID,
	}

	expectedResps := []*pb.JoinCustomMatchResponse{r1}
	for i := 3; i < 11; i++ {
		go te.join(t, uint64(i), groupID, "")

		var expectedUsers []uint64
		for j := 1; j < i+1; j++ {
			expectedUsers = append(expectedUsers, uint64(j))
		}
		r := &pb.JoinCustomMatchResponse{
			Users:    expectedUsers,
			MaxUsers: in.GetMaxUsers(),
			Leader:   in.GetLeader(),
			GroupId:  groupID,
		}
		expectedResps = append(expectedResps, r)
		time.Sleep(time.Second * 2)
	}

	// need to call start

	wg.Wait()
	assert.NoError(t, err)
	assert.Equal(t, expectedResps, resps)
}

func TestServiceJoin(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddCustomMatchRequest{
		Leader:   1,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}
	addResps, _ := te.add(t, addIn)

	resps, err := te.join(t, 2, addResps[0].GetGroupId(), "")
	assert.NoError(t, err)

	r1 := &pb.JoinCustomMatchResponse{
		Users:    []uint64{1, 2},
		MaxUsers: addIn.MaxUsers,
		Leader:   addIn.Leader,
		GroupId:  resps[0].GetGroupId(),
	}
	assert.Equal(t, []*pb.JoinCustomMatchResponse{r1}, resps)
}

func TestServiceLeaveGroupDoesNotExist(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.LeaveCustomMatchRequest{
		UserId:  1,
		GroupId: 1,
	}
	_, err := te.c.Leave(context.Background(), in)
	assert.Error(t, err)
}

func TestServiceLeaveNotInGroup(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddCustomMatchRequest{
		Leader:   1,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}
	addResps, _ := te.add(t, addIn)

	in := &pb.LeaveCustomMatchRequest{
		UserId:  2,
		GroupId: addResps[0].GetGroupId(),
	}
	_, err := te.c.Leave(context.Background(), in)
	assert.Error(t, err)
}

func TestServiceLeave(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddCustomMatchRequest{
		Leader:   1,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}
	addResps, _ := te.add(t, addIn)

	te.join(t, 2, addResps[0].GetGroupId(), "")

	in := &pb.LeaveCustomMatchRequest{
		UserId:  2,
		GroupId: addResps[0].GetGroupId(),
	}
	_, err := te.c.Leave(context.Background(), in)
	assert.NoError(t, err)
}

func TestServiceStartGroupDoesNotExist(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.StartCustomMatchRequest{
		UserId:  1,
		GroupId: 1,
	}
	_, err := te.c.Start(context.Background(), in)
	assert.Error(t, err)
}

func TestServiceStartNotLeader(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddCustomMatchRequest{
		Leader:   1,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}
	addResps, _ := te.add(t, addIn)
	te.join(t, 2, addResps[0].GetGroupId(), "")

	in := &pb.StartCustomMatchRequest{
		UserId:  2,
		GroupId: addResps[0].GetGroupId(),
	}
	_, err := te.c.Start(context.Background(), in)
	assert.Error(t, err)
}


func TestServiceStart(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	addIn := &pb.AddCustomMatchRequest{
		Leader:   1,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}
	addResps, _ := te.add(t, addIn)

	in := &pb.StartCustomMatchRequest{
		UserId:  1,
		GroupId: addResps[0].GetGroupId(),
	}
	resp, err := te.c.Start(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.StartCustomMatchResponse{
		Users: []uint64{1},
		Leader: addIn.Leader,
		GroupId: addResps[0].GetGroupId(),
	}

	assert.Equal(t, expectedResp, resp)
}
