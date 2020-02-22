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
	doneChs []chan struct{}
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

	return &test{c, service, []chan struct{}{}}
}

func (te *test) teardown() {
	for _, v := range te.doneChs {
		v <- struct{}{}
	}
	te.service.Close()
}

func (te *test) fillData(t testing.TB) []Group {
	in := &pb.AddCustomMatchRequest{
		Leader:   1,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}
	_, err := te.add(t, in)
	assert.NoError(t, err)

	in2 := &pb.AddCustomMatchRequest{
		Leader:   2,
		Name:     "game2",
		Password: "qwe",
		MaxUsers: 10,
	}
	resp2, err := te.add(t, in2)
	assert.NoError(t, err)
	go te.join(t, 3, resp2[0].GetGroupId(), "qwe")

	in3 := &pb.AddCustomMatchRequest{
		Leader:   4,
		Name:     "game3",
		Password: "",
		MaxUsers: 10,
	}
	resp3, err := te.add(t, in3)
	assert.NoError(t, err)
	for i := 5; i < 14; i++ {
		go te.join(t, uint64(i), resp3[0].GetGroupId(), "")
	}
	time.Sleep(time.Second * 5)

	var groups []Group
	err = te.service.s.db.Preload("Users").Find(&groups).Error
	assert.NoError(t, err)
	return groups
}

func (te *test) queryGroups(t testing.TB) []Group {
	var groups []Group
	err := te.service.s.db.Preload("Users").Find(&groups).Error
	assert.NoError(t, err)
	return groups
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

	done := make(chan struct{}, 1)
	te.doneChs = append(te.doneChs, done)

	for {
		select {
		case <- done:
			return resps, nil
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

	done := make(chan struct{}, 1)
	te.doneChs = append(te.doneChs, done)

	for {
		select {
		case <- done:
			return resps, nil
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

	groups := te.fillData(t)

	in := &pb.AddCustomMatchRequest{
		Leader:   groups[0].Leader,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}
	_, err := te.add(t, in)
	assert.Error(t, err)
	groupsAfter := te.queryGroups(t)
	assert.Equal(t, groups, groupsAfter)
}

func TestServiceAddLeaderAlreadyInGroup(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	groups := te.fillData(t)

	in := &pb.AddCustomMatchRequest{
		Leader:   groups[2].Users[1].UserID,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}
	_, err := te.add(t, in)
	assert.Error(t, err)
	groupsAfter := te.queryGroups(t)
	assert.Equal(t, groups, groupsAfter)
}

func TestServiceAddStart(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	groups := te.fillData(t)

	in := &pb.AddCustomMatchRequest{
		Leader:   101,
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

	groupID := groups[len(groups)-1].ID + 1

	r1 := &pb.AddCustomMatchResponse{
		Users:    []uint64{in.GetLeader()},
		MaxUsers: in.GetMaxUsers(),
		Leader:   in.GetLeader(),
		GroupId:  groupID,
	}

	expectedResps := []*pb.AddCustomMatchResponse{r1}
	for i := 102; i < 111; i++ {
		go te.join(t, uint64(i), groupID, "")

		var expectedUsers []uint64
		for j := 101; j < i+1; j++ {
			expectedUsers = append(expectedUsers, uint64(j))
		}
		r := &pb.AddCustomMatchResponse{
			Users:    expectedUsers,
			MaxUsers: in.GetMaxUsers(),
			Leader:   in.GetLeader(),
			GroupId:  groupID,
		}
		expectedResps = append(expectedResps, r)
		time.Sleep(time.Second * 2)
	}

	startIn := &pb.StartCustomMatchRequest{
		UserId:  in.GetLeader(),
		GroupId: groupID,
	}
	te.c.Start(context.Background(), startIn)

	r2 := &pb.AddCustomMatchResponse{
		Users:    []uint64{101, 102, 103, 104, 105, 106, 107, 108, 109, 110},
		MaxUsers: in.GetMaxUsers(),
		Leader:   in.GetLeader(),
		GroupId:  groupID,
		Started:  true,
	}
	expectedResps = append(expectedResps, r2)

	wg.Wait()
	assert.NoError(t, err)
	assert.Equal(t, expectedResps, resps)
	groupsAfter := te.queryGroups(t)
	assert.Equal(t, groups, groupsAfter)
}

func TestServiceAdd(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	groups := te.fillData(t)

	in := &pb.AddCustomMatchRequest{
		Leader:   100,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}
	resps, err := te.add(t, in)
	assert.NoError(t, err)

	expectedGroupID := groups[len(groups)-1].ID + 1
	r1 := &pb.AddCustomMatchResponse{
		Users:    []uint64{100},
		MaxUsers: in.GetMaxUsers(),
		Leader:   in.GetLeader(),
		GroupId:  expectedGroupID,
	}

	assert.Equal(t, []*pb.AddCustomMatchResponse{r1}, resps)
	groupsAfter := te.queryGroups(t)

	expectedGroup := Group{
		ID:       expectedGroupID,
		Leader:   in.GetLeader(),
		Name:     in.GetName(),
		Password: in.GetPassword(),
		MaxUsers: in.GetMaxUsers(),
		Users:    []User{User{100, expectedGroupID}},
	}
	groups = append(groups, expectedGroup)
	assert.Equal(t, groups, groupsAfter)
}

func TestServiceDeleteDoesNotExist(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	groups := te.fillData(t)

	in := &pb.DeleteCustomMatchRequest{
		UserId:  100,
		GroupId: 20,
	}
	_, err := te.c.Delete(context.Background(), in)
	assert.Error(t, err)
	groupsAfter := te.queryGroups(t)
	assert.Equal(t, groups, groupsAfter)
}

func TestServiceDeleteNotLeader(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	groups := te.fillData(t)

	in := &pb.DeleteCustomMatchRequest{
		UserId:  groups[0].Leader + 100,
		GroupId: groups[0].ID,
	}
	_, err := te.c.Delete(context.Background(), in)
	assert.Error(t, err)
	groupsAfter := te.queryGroups(t)
	assert.Equal(t, groups, groupsAfter)
}

func TestServiceDelete(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	groups := te.fillData(t)

	addIn := &pb.AddCustomMatchRequest{
		Leader:   100,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}
	resps, _ := te.add(t, addIn)

	in := &pb.DeleteCustomMatchRequest{
		UserId:  100,
		GroupId: resps[0].GetGroupId(),
	}
	_, err := te.c.Delete(context.Background(), in)
	assert.NoError(t, err)
	groupsAfter := te.queryGroups(t)
	assert.Equal(t, groups, groupsAfter)
}

func TestServiceGetAllNone(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	in := &pb.GetAllCustomMatchRequest{Total: 100, Skip: 0}
	resp, err := te.c.GetAll(context.Background(), in)
	assert.NoError(t, err)
	assert.Len(t, resp.Groups, 0)
	groups := te.queryGroups(t)
	assert.Len(t, groups, 0)
}

func TestServiceGetAll(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	groups := te.fillData(t)

	in := &pb.GetAllCustomMatchRequest{Total: 100, Skip: 0}
	resp, err := te.c.GetAll(context.Background(), in)
	assert.NoError(t, err)
	g1 := &pb.CustomMatchGroup{
		GroupId:          groups[0].ID,
		Name:             groups[0].Name,
		Leader:           groups[0].Leader,
		PasswordRequired: false,
		MaxUsers:         groups[0].MaxUsers,
		TotalUsers:       uint32(len(groups[0].Users)),
	}
	g2 := &pb.CustomMatchGroup{
		GroupId:          groups[1].ID,
		Name:             groups[1].Name,
		Leader:           groups[1].Leader,
		PasswordRequired: true,
		MaxUsers:         groups[1].MaxUsers,
		TotalUsers:       uint32(len(groups[1].Users)),
	}
	g3 := &pb.CustomMatchGroup{
		GroupId:          groups[2].ID,
		Name:             groups[2].Name,
		Leader:           groups[2].Leader,
		PasswordRequired: false,
		MaxUsers:         groups[2].MaxUsers,
		TotalUsers:       uint32(len(groups[2].Users)),
	}
	expectedResp := &pb.GetAllCustomMatchResponse{
		Groups: []*pb.CustomMatchGroup{g1, g2, g3},
	}
	assert.Equal(t, expectedResp, resp)
	groupsAfter := te.queryGroups(t)
	assert.Equal(t, groups, groupsAfter)
}

func TestServiceJoinDoesNotExist(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	groups := te.fillData(t)

	_, err := te.join(t, 100, 20, "")
	assert.Error(t, err)
	groupsAfter := te.queryGroups(t)
	assert.Equal(t, groups, groupsAfter)
}

func TestServiceJoinWrongPassword(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	groups := te.fillData(t)

	_, err := te.join(t, 100, groups[1].ID, "qqq")
	assert.Error(t, err)
	groupsAfter := te.queryGroups(t)
	assert.Equal(t, groups, groupsAfter)
}

func TestServiceJoinPassword(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	groups := te.fillData(t)

	resps, err := te.join(t, 100, groups[1].ID, groups[1].Password)
	assert.NoError(t, err)

	r1 := &pb.JoinCustomMatchResponse{
		Users:    []uint64{groups[1].Users[0].UserID, groups[1].Users[1].UserID, 100},
		MaxUsers: groups[1].MaxUsers,
		Leader:   groups[1].Leader,
		GroupId:  groups[1].ID,
	}
	assert.Equal(t, []*pb.JoinCustomMatchResponse{r1}, resps)
	groupsAfter := te.queryGroups(t)
	updatedGroup := groups[1]
	updatedGroup.Users = append(updatedGroup.Users, User{100, updatedGroup.ID})
	groups[1] = updatedGroup
	assert.Equal(t, groups, groupsAfter)
}

func TestServiceJoinFull(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	groups := te.fillData(t)

	_, err := te.join(t, 100, groups[2].ID, "")
	assert.Error(t, err)
	groupsAfter := te.queryGroups(t)
	assert.Equal(t, groups, groupsAfter)
}

func TestServiceJoinStart(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	groups := te.fillData(t)

	in := &pb.AddCustomMatchRequest{
		Leader:   101,
		Name:     "game",
		Password: "",
		MaxUsers: 10,
	}

	te.add(t, in)
	groupID := groups[len(groups)-1].ID + 1

	var wg sync.WaitGroup
	var resps []*pb.JoinCustomMatchResponse
	var err error
	wg.Add(1)
	go func() {
		resps, err = te.join(t, 102, groupID, "")
		wg.Done()
	}()
	time.Sleep(time.Second * 2)

	r1 := &pb.JoinCustomMatchResponse{
		Users:    []uint64{101, 102},
		MaxUsers: in.MaxUsers,
		Leader:   in.Leader,
		GroupId:  groupID,
	}

	expectedResps := []*pb.JoinCustomMatchResponse{r1}
	for i := 103; i < 111; i++ {
		go te.join(t, uint64(i), groupID, "")

		var expectedUsers []uint64
		for j := 101; j < i+1; j++ {
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

	startIn := &pb.StartCustomMatchRequest{
		UserId:  in.GetLeader(),
		GroupId: groupID,
	}
	te.c.Start(context.Background(), startIn)

	r2 := &pb.JoinCustomMatchResponse{
		Users:    []uint64{101, 102, 103, 104, 105, 106, 107, 108, 109, 110},
		MaxUsers: in.GetMaxUsers(),
		Leader:   in.GetLeader(),
		GroupId:  groupID,
		Started:  true,
	}
	expectedResps = append(expectedResps, r2)

	wg.Wait()
	assert.NoError(t, err)
	assert.Equal(t, expectedResps, resps)
	groupsAfter := te.queryGroups(t)
 	assert.Equal(t, groups, groupsAfter)
}

func TestServiceJoin(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	groups := te.fillData(t)

	resps, err := te.join(t, 100, groups[0].ID, "")
	assert.NoError(t, err)

	r1 := &pb.JoinCustomMatchResponse{
		Users:    []uint64{groups[0].Leader, 100},
		MaxUsers: groups[0].MaxUsers,
		Leader:   groups[0].Leader,
		GroupId:  groups[0].ID,
	}
	assert.Equal(t, []*pb.JoinCustomMatchResponse{r1}, resps)
	groupsAfter := te.queryGroups(t)
	updatedGroup := groups[0]
	updatedGroup.Users = append(updatedGroup.Users, User{100, updatedGroup.ID})
	groups[0] = updatedGroup
	assert.Equal(t, groups, groupsAfter)
}

func TestServiceLeaveGroupDoesNotExist(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	groups := te.fillData(t)

	in := &pb.LeaveCustomMatchRequest{
		UserId:  100,
		GroupId: 20,
	}
	_, err := te.c.Leave(context.Background(), in)
	assert.Error(t, err)
	groupsAfter := te.queryGroups(t)
	assert.Equal(t, groups, groupsAfter)
}

func TestServiceLeaveNotInGroup(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	groups := te.fillData(t)

	in := &pb.LeaveCustomMatchRequest{
		UserId:  100,
		GroupId: groups[0].ID,
	}
	_, err := te.c.Leave(context.Background(), in)
	assert.Error(t, err)
	groupsAfter := te.queryGroups(t)
	assert.Equal(t, groups, groupsAfter)
}

func TestServiceLeave(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	groups := te.fillData(t)

	_, err := te.join(t, 100, groups[0].ID, "")
	assert.NoError(t, err)

	in := &pb.LeaveCustomMatchRequest{
		UserId:  100,
		GroupId: groups[0].ID,
	}
	_, err = te.c.Leave(context.Background(), in)
	assert.NoError(t, err)
	groupsAfter := te.queryGroups(t)
	assert.Equal(t, groups, groupsAfter)
}

func TestServiceStartGroupDoesNotExist(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	groups := te.fillData(t)

	in := &pb.StartCustomMatchRequest{
		UserId:  100,
		GroupId: 20,
	}
	_, err := te.c.Start(context.Background(), in)
	assert.Error(t, err)

	groupsAfter := te.queryGroups(t)
	assert.Equal(t, groups, groupsAfter)
}

func TestServiceStartNotLeader(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	groups := te.fillData(t)

	in := &pb.StartCustomMatchRequest{
		UserId:  100,
		GroupId: groups[0].ID,
	}
	_, err := te.c.Start(context.Background(), in)
	assert.Error(t, err)
	groupsAfter := te.queryGroups(t)
	assert.Equal(t, groups, groupsAfter)
}

func TestServiceStart(t *testing.T) {
	te := newTest(t)
	defer te.teardown()

	groups := te.fillData(t)

	in := &pb.StartCustomMatchRequest{
		UserId:  groups[0].Leader,
		GroupId: groups[0].ID,
	}
	resp, err := te.c.Start(context.Background(), in)
	assert.NoError(t, err)
	expectedResp := &pb.StartCustomMatchResponse{
		Users:   []uint64{groups[0].Leader},
		Leader:  groups[0].Leader,
		GroupId: groups[0].ID,
	}

	assert.Equal(t, expectedResp, resp)
	groupsAfter := te.queryGroups(t)
	assert.Equal(t, groups[1:], groupsAfter)
}
