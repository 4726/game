package app

import (
	"sync"
	"time"

	"github.com/4726/game/gameplay/5v5fps/engine"
	"github.com/4726/game/gameplay/5v5fps/pb"
	"github.com/4726/game/gameplay/5v5fps/player"
	"github.com/4726/game/gameplay/5v5fps/util"
	"github.com/4726/game/gameplay/5v5fps/weapon"
)

type gameServer struct {
	players       map[uint64]*player.Player
	playersLock   sync.Mutex
	weapons       []weapon.Weapon
	en            engine.Engine
	playerStreams map[uint64]chan *pb.GameStateResponse
}

func newGameServer(team1, team2 [5]uint64) *gameServer {
	players := map[uint64]*player.Player{}
	playerStreams := map[uint64]chan *pb.GameStateResponse{}
	for _, v := range team1 {
		p := &player.Player{
			UserID: v,
			Team:   util.Team1,
		}
		players[v] = p
		ch := make(chan *pb.GameStateResponse)
		playerStreams[v] = ch
	}
	for _, v := range team2 {
		p := &player.Player{
			UserID: v,
			Team:   util.Team2,
		}
		players[v] = p
		ch := make(chan *pb.GameStateResponse)
		playerStreams[v] = ch
	}
	return &gameServer{
		players:       players,
		playerStreams: playerStreams,
	}
}

func (s *gameServer) Connect(stream pb.Game_ConnectServer) error {
	var userID uint64
	errChannel := make(chan error)
	for {
		select {
		case err := <-errChannel:
			return err
		default:
		}
		in, err := stream.Recv()
		if err != nil {
			return err
		}

		if initReq := in.GetInit(); initReq != nil {
			s.playersLock.Lock()
			p := s.players[initReq.GetUserId()]
			p.Connected = true
			userID = initReq.GetUserId()

			streamChannel := s.playerStreams[initReq.GetUserId()]
			go func() {
				for msg := range streamChannel {
					out := &pb.GameServerData{
						Data: &pb.GameServerData_State{
							State: msg,
						},
					}
					err := stream.Send(out)
					if err != nil {
						errChannel <- err
						return
					}
				}
			}()

			for _, v := range s.players {
				if !v.Connected {
					s.playersLock.Unlock()
					continue
				}
			}

			for _, v := range s.players {
				var initPosition util.Vector3
				if v.Team == util.Team1 {
					initPosition = util.Vector3{500, 0, 0}
				} else {
					initPosition = util.Vector3{0, 0, 0}
				}
				enginePlayer := engine.Player{
					UserID:          v.UserID,
					Position:        initPosition,
					HP:              100,
					Team:            v.Team,
					PrimaryWeapon:   nil,
					SecondaryWeapon: &weapon.SecondaryOne,
					KnifeWeapon:     nil,
					EquippedWeapon:  &weapon.SecondaryOne,
					Dead:            false,
					UserScore: engine.Score{0, 0, 0},
				}
				s.en.Init(enginePlayer)
			}
			s.en.Start()
			s.playersLock.Unlock()

			continue
		}

		if buyReq := in.GetBuy(); buyReq != nil {
			s.en.Buy(userID, int(buyReq.GetWeaponId()))

			out := &pb.GameServerData{
				Data: &pb.GameServerData_Buy{
					Buy: &pb.GameBuyResponse{
						Success: true,
					},
				}, 
			}
			err := stream.Send(out)
			if err != nil {
				return err
			}
			continue
		}

		if inputReq := in.GetInput(); inputReq != nil {
			s.playersLock.Lock()
			p := s.players[userID]

			if in := inputReq.GetPrimaryWeapon(); in != nil {
				p.EquippedWeapon = p.PrimaryWeapon
			} else if in := inputReq.GetSecondaryWeapon(); in != nil {
				p.EquippedWeapon = p.SecondaryWeapon
			} else if in := inputReq.GetKnifeWeapon(); in != nil {
				p.EquippedWeapon = p.KnifeWeapon
			} else if in := inputReq.GetMoveLeft(); in != nil {
				s.en.MoveLeft(p.UserID)
			} else if in := inputReq.GetMoveLeftUp(); in != nil {
				s.en.MoveLeftUp(p.UserID)
			} else if in := inputReq.GetMoveLeftDown(); in != nil {
				s.en.MoveLeftDown(p.UserID)
			} else if in := inputReq.GetMoveRight(); in != nil {
				s.en.MoveRight(p.UserID)
			} else if in := inputReq.GetMoveRightUp(); in != nil {
				s.en.MoveRightUp(p.UserID)
			} else if in := inputReq.GetMoveRightDown(); in != nil {
				s.en.MoveRightDown(p.UserID)
			} else if in := inputReq.GetMoveUp(); in != nil {
				s.en.MoveUp(p.UserID)
			} else if in := inputReq.GetMoveUpLeft(); in != nil {
				s.en.MoveUpLeft(p.UserID)
			} else if in := inputReq.GetMoveUpRight(); in != nil {
				s.en.MoveUpRight(p.UserID)
			} else if in := inputReq.GetMoveDown(); in != nil {
				s.en.MoveDown(p.UserID)
			} else if in := inputReq.GetMoveDownLeft(); in != nil {
				s.en.MoveDownLeft(p.UserID)
			} else if in := inputReq.GetMoveDownRight(); in != nil {
				s.en.MoveDownRight(p.UserID)
			} else if in := inputReq.GetShoot(); in != nil {
				s.en.Shoot(p.UserID, util.Vector3{
					in.GetTargetX(),
					in.GetTargetY(),
					in.GetTargetZ(),
				})
			} else if in := inputReq.GetReload(); in != nil {
				p.EquippedWeapon.Ammo = p.EquippedWeapon.AmmoMax
			} else if in := inputReq.GetPickupWeapon(); in != nil {
				s.en.PickupWeapon(p.UserID, int(in.GetWeaponId()))
			} else if in := inputReq.GetCrouch(); in != nil {

			} else if in := inputReq.GetJump(); in != nil {

			} else if in := inputReq.GetPing(); in != nil {

			} else if in := inputReq.GetOrientation(); in != nil {
				s.en.SetOrientation(userID, util.Vector3{
					X: in.GetX(),
					Y: in.GetY(),
					Z: in.GetZ(),
				})
			}

			p.LastUpdate = time.Now()
			s.playersLock.Unlock()
			continue
		}
	}
}

func (s *gameServer) runGameLoop() {
	for {
		s.playersLock.Lock()
		for k := range s.players {
			all := s.en.All(k)

			streamChannel := s.playerStreams[k]
			streamChannel <- &pb.GameStateResponse{
				Players: enginePlayersToPBPlayers(all),
			}
		}
		s.playersLock.Unlock()
		time.Sleep(time.Millisecond * 250)
	}
}

func enginePlayersToPBPlayers(eps []engine.Player) []*pb.GamePlayerData {
	var pbPlayers []*pb.GamePlayerData
	for _, v := range eps {
		pbPlayer := &pb.GamePlayerData{
			UserId:         v.UserID,
			Teammate:       true,
			KnownPos:       v.Private,
			PosX:           v.Position.X,
			PosY:           v.Position.Y,
			PosZ:           v.Position.Z,
			OrientationX:   v.Orientation.X,
			OrientationY:   v.Orientation.Y,
			OrientationZ:   v.Orientation.Z,
			Dead:           v.Dead,
			Connected:      true,
			EquippedWeapon: int32(v.EquippedWeapon.ID),
			ScoreKills:     int32(v.UserScore.Kills),
			ScoreDeaths:    int32(v.UserScore.Deaths),
			ScoreAssists:   int32(v.UserScore.Assists),
			Money:          int32(v.Money),
		}
		pbPlayers = append(pbPlayers, pbPlayer)
	}

	return pbPlayers
}
