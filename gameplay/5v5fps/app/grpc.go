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
}

func newGameServer(team1, team2 [5]uint64) *gameServer {
	players := map[uint64]*player.Player{}
	for _, v := range team1 {
		p := &player.Player{
			UserID: v,
			Team: player.Team1,
		}
		players[v] = p
	}
	for _, v := range team2 {
		p := &player.Player{
			UserID: v,
			Team: player.Team2,
		}
		players[v] = p
	}
	return &gameServer{
		players: players,
	}
}

func (s *gameServer) Connect(stream pb.GameConnectServer) error {
	var userID uint64
	for {
		in, err := stream.Recv()
		if err != nil {
			return err
		}

		data := in.Data()

		if initReq := data.GetInit(); initReq != nil {
			s.playersLock.Lock()
			p := s.players[in.GetUserId()]
			p.Connected = true
			s.playersLock.Unlock()
			userID = initReq.GetUserId()
			continue
		}

		if buyReq := data.GetBuy(); buyReq != nil {
			var weapon weapon.Weapon
			for _, v := range s.weapons {
				if v.ID == buyReq,GetWeaponId() {
					weapon = v
					break
				}
			}
			s.playersLock.Lock()
			p := s.players[userID]
			if p.Money < weapon.Price {
				out := &pb.GameServerData{
					Data: &pb.GameBuyResponse{
						Success: false,
						Error: "not enough money",
					}
				}
				err := stream.Send(out)
				if err != nil {
					s.playersLock.Unlock()
					return err
				}
			} else {
				p.Money -= weapon.Price
				switch weapon.WT {
				case weapon.Primary:
					p.PrimaryWeapon = weapon
				case weapon.Secondary:
					p.SecondaryWeapon = weapon
				case weapon.Knife:
					p.KnifeWeapon = weapon
				}
				out := &pb.GameServerData{
					Data: &pb.GameBuyResponse{
						Success: true,
					}
				}
				err := stream.Send(out)
				if err != nil {
					s.playersLock.Unlock()
					return err
				}
			}
			p.LastUpdate = time.now()
			s.playersLock.Unlock()
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

			}

			p.LastUpdate = time.now()
			s.playersLock.Unlock()
			continue
		}
}