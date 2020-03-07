package app

import (
	"github.com/4726/game/gameplay/player"
	"github.com/4726/game/gameplay/weapon"
	"github.com/4726/game/gameplay/engine"
	"github.com/4726/game/gameplay/util"
)

type gameServer struct {
	players map[uint64]*player.Player
	playersLock sync.Mutex
	weapons []weapon.Weapon
	en engine.Engine
	playerStreams map[uint64]chan *pb.GameStateResponse
}

func newGameServer(team1, team2 [5]uint64) *gameServer {
	players := map[uint64]*player.Player{}
	playerStreams := map[uint64]chan *pb.GameStateResponse{}
	for _, v := range team1 {
		p := &player.Player{
			UserID: v,
			Team: player.Team1,
		}
		players[v] = p
		ch := make(chan *pb.GameStateResponse)
		playerStreams[v] = ch
	}
	for _, v := range team2 {
		p := &player.Player{
			UserID: v,
			Team: player.Team2,
		}
		players[v] = p
		ch := make(chan *pb.GameStateResponse)
		playerStreams[v] = ch
	}
	return &gameServer{
		players: players,
		playerStreams: playerStreams
	}
}

func (s *gameServer) Connect(stream pb.GameConnectServer) error {
	var userID uint64
	errChannel := make(chan error)
	for {
		select {
		case err := <- errChannel:
			return err
		default:
		}
		in, err := stream.Recv()
		if err != nil {
			return err
		}

		data := in.Data()

		if initReq := data.GetInit(); initReq != nil {
			s.playersLock.Lock()
			p := s.players[in.GetUserId()]
			p.Connected = true
			userID = initReq.GetUserId()
			
			streamChannel := s.players[in.GetUserId()]
			go func() {
				for msg := range streamChannel {
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
					UserID: v.UserID,
					Position: initPosition,
					HP: 100,
					Team: v.Team,
					PrimaryWeapon: nil,
					SecondaryWeapon: weapon.SecondaryOne,
					KnifeWeapon: nil,
					EquippedWeapon: weapon.SecondaryOne,
					Dead: false,
					Kills: 0,
					Deaths: 0,
					Assists: 0,
				}
				s.en.Init(enginePlayer)
			}
			s.en.Start()
			s.playersLock.Unlock()

			continue
		}

		if buyReq := data.GetBuy(); buyReq != nil {
			s.en.Buy(userID, buyReq.GetWeaponId())
			out := &pb.GameServerData{
				Data: &pb.GameBuyResponse{
					Success: true,
				}
			}
			err := stream.Send(out)
			if err != nil {
				return err
			}
			continue
		}


		if inputReq := data.GetInput(); inputReq != nil {
			s.playersLock.Lock()
			p := s.players[userID]
			key := inputReq.GetInputKey()

			if in := inputReq.GetPrimaryWeapon(); in != nil {
				p.EquippedWeapon = p.PrimaryWeapon
			} else if in := inputReq.GetSecondaryWeapon(); in != nil {
				p.EquippedWeapon = p.SecondaryWeapon
			} else if in := inputReq.GetKnifeWeapon(); in != nil {
				p.EquippedWeapon = p.KnifeWeapon
			} else if in := inputReq.GetMoveLeft(); in != nil {
				s.engine.MoveLeft(p.UserID)
			} else if in := inputReq.GetMoveLeftUp(); in != nil {
				s.engine.MoveLeftUp(p.UserID)
			} else if in := inputReq.GetMoveLeftDown(); in != nil {
				s.engine.MoveLeftDown(p.UserID)
			} else if in := inputReq.GetMoveRight(); in != nil {
				s.engine.MoveRight(p.UserID)
			} else if in := inputReq.GetMoveRightUp(); in != nil {
				s.engine.MoveRightUp(p.UserID)
			} else if in := inputReq.GetMoveRightDown(); in != nil {
				s.engine.MoveRightDown(p.UserID)
			} else if in := inputReq.GetMoveUp(); in != nil {
				s.engine.MoveUp(p.UserID)
			} else if in := inputReq.GetMoveUpLeft(); in != nil {
				s.engine.MoveUpLeft(p.UserID)
			} else if in := inputReq.GetMoveUpRight(); in != nil {
				s.engine.MoveUpRight(p.UserID)
			} else if in := inputReq.GetMoveDown(); in != nil {
				s.engine.MoveDown(p.UserID)
			} else if in := inputReq.GetMoveDownLeft(); in != nil {
				s.engine.MoveDownLeft(p.UserID)
			} else if in := inputReq.GetMoveDownRight(); in != nil {
				s.engine.MoveDownRight(p.UserID)
			} else if in := inputReq.GetShoot(); in != nil {
				s.engine.Shoot(p.UserID, util.Vector3{
					in.GetTargetX(),
					in.GetTargetY(),
					in.GetTargetZ(),
				})
			} else if in := inputReq.GetReload(); in != nil {
				p.EquippedWeapon.Ammo = p.EquippedWeapon.AmmoMax
			} else if in := inputReq.GetPickupWeapon(); in != nil {
				s.engine.PickupWeapon(p.UserID, in.GetWeaponId())
			} else if in := inputReq.GetCrouch(); in != nil {

			} else if in := inputReq.GetJump(); in != nil {

			} else if in := inputReq.GetPing(); in != nil {

			} else if in := inputReq.GetOrientation(); in != nil {
				s.engine.SetOrientation(userID, in.GetOrientation)
			}

			p.LastUpdate = time.now()
			s.playersLock.Unlock()
			continue
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

func enginePlayersToPBPlayers(eps []engine.Player) *pb.GamePlayerData {
	var pbPlayers []*pb.GamePlayerData
	for _, v := range eps {
		pbPlayer := &pb.GamePlayerData{
			UserId: v.UserID,
			Teammate: true,
			KnownPos: v.Private,
			PosX: v.Position.X,
			PosY: v.Position.Y,
			PosZ: v.Position.Z,
			OrientationX: v.Orientation.X,
			OrientationY: v.Orientation.Y,
			OrientationZ: v.Orientation.Z,
			Dead: v.Dead,
			Connected: true,
			EquippedWeapon: v.EquippedWeapon.ID,
			ScoreKills: v.UserScore.Kills,
			ScoreDeaths: v.UserScore.Deaths,
			ScoreAssists: v.UserScore.Assists,
			Money: v.Money,
		}
	}

	return pbPlayers
}