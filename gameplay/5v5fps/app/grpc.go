package app

import (
	"github.com/4726/game/gameplay/player"
	"github.com/4726/game/gameplay/weapon"
)

type gameServer struct {
	players map[uint64]*player.Player
	playersLock sync.Mutex
	weapons []weapon.Weapon
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
			s.playersLock.Unlock()
			continue
		}


		if inputReq := data.GetInput(); inputReq != nil {
			s.playersLock.Lock()
			p := s.players[userID]
			key := inputReq.GetInputKey()

			if key == "primary_weapon" {
				p.EquippedWeapon = p.PrimaryWeapon
			}
			if key == "secondary_weapon" {
				p.EquippedWeapon = p.SecondaryWeapon
			}
			if key == "knife_weapon" {
				p.EquippedWeapon = p.KnifeWeapon
			}
			if key == "slow_walk" {

			}
			if key == "normal_walk" {

			}
			if key == "pickup_weapon" {
				
			}
			if key == "crouch" {
				
			}
			if key == "jump" {
				
			}
			if key == "ping" {
				
			}

			p.LastUpdate = time.now()
			s.playersLock.Unlock()
			continue
		}
}