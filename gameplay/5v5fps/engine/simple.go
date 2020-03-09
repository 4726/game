package engine

import (
	"sync"

	"github.com/4726/game/gameplay/5v5fps/util"
	"github.com/4726/game/gameplay/5v5fps/weapon"
)

// SimpleEngine implements Engine
type SimpleEngine struct {
	players map[uint64]*Player
	sync.Mutex
	droppedWeapons []SimpleWeapon
	weapons        []weapon.Weapon
	ch             chan interface{}
	state          State
}

type SimpleWeapon struct {
	Location util.Vector3
	WeaponID int
}

func NewSimpleEngine() *SimpleEngine {
	return &SimpleEngine{
		players: map[uint64]*Player{},
		ch:      make(chan interface{}),
		state: State{
			Scores: map[util.TeamID]int{},
		},
	}
}

func (e *SimpleEngine) Init(p Player) {
	e.Lock()
	defer e.Unlock()

	e.players[p.UserID] = &p
}

func (e *SimpleEngine) MoveLeft(userID uint64) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X - 1, p.Position.Y, p.Position.Z}
}

func (e *SimpleEngine) MoveLeftUp(userID uint64) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X - 1, p.Position.Y + 1, p.Position.Z}
}

func (e *SimpleEngine) MoveLeftDown(userID uint64) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X - 1, p.Position.Y - 1, p.Position.Z}
}

func (e *SimpleEngine) MoveRight(userID uint64) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X + 1, p.Position.Y, p.Position.Z}
}

func (e *SimpleEngine) MoveRightUp(userID uint64) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X + 1, p.Position.Y + 1, p.Position.Z}
}

func (e *SimpleEngine) MoveRightDown(userID uint64) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X + 1, p.Position.Y - 1, p.Position.Z}
}

func (e *SimpleEngine) MoveUp(userID uint64) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X, p.Position.Y + 1, p.Position.Z}
}

func (e *SimpleEngine) MoveUpLeft(userID uint64) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X - 1, p.Position.Y + 1, p.Position.Z}
}

func (e *SimpleEngine) MoveUpRight(userID uint64) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X + 1, p.Position.Y + 1, p.Position.Z}
}

func (e *SimpleEngine) MoveDown(userID uint64) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X, p.Position.Y - 1, p.Position.Z}
}

func (e *SimpleEngine) MoveDownLeft(userID uint64) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X - 1, p.Position.Y - 1, p.Position.Z}
}

func (e *SimpleEngine) MoveDownRight(userID uint64) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X + 1, p.Position.Y - 1, p.Position.Z}
}

func (e *SimpleEngine) Shoot(userID uint64, target util.Vector3) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}

	if p.EquippedWeapon.Ammo = 0 {
		return
	}

	p.EquippedWeapon.Ammo--

	var haveDeath bool
	for _, v := range e.players {
		if v.UserID == userID {
			continue
		}

		if v.Position != target {
			continue
		}

		v.HP -= int(p.EquippedWeapon.Damage)
		if v.HP <= 0 {
			v.Dead = true
			haveDeath = true
			p.UserScore.Kills++
		}
	}

	if !haveDeath {
		return
	}

	teamDeathCount := map[util.TeamID]int{}
	for _, v := range e.players {
		if v.Dead {
			count, ok := teamDeathCount[v.Team]
			if !ok {
				teamDeathCount[v.Team] = 1
			} else {
				teamDeathCount[v.Team] = count + 1
			}
		}
	}

	for k, v := range teamDeathCount {
		var winner util.TeamID
		var loser util.TeamID
		if v == 5 {
			loser = k
			if loser == util.Team1 {
				winner = util.Team2
			} else {
				winner = util.Team1
			}
			e.ch <- RoundEnd{
				WinningTeam: winner,
				LosingTeam:  loser,
			}

			score := e.state.Scores[winner]
			e.state.Scores[winner] = score + 1

			if e.state.Scores[winner] == 16 {
				e.ch <- MatchEnd{
					WinningTeam: winner,
					LosingTeam:  loser,
				}
			} else {
				e.ch <- RoundStart{}
			}
			break
		}
	}
}

func (e *SimpleEngine) All(userID uint64) []Player {
	e.Lock()
	defer e.Unlock()

	var players []Player
	p, ok := e.players[userID]
	if !ok {
		return players
	}
	for _, v := range e.players {
		playerData := *v

		//doesn't return information that should be private
		if e.inLineOfSight(*p, playerData) && v.Team != p.Team {
			playerData.Private = false
			playerData.Position = v.Position
			playerData.EquippedWeapon = v.EquippedWeapon
			playerData.Orientation = util.Vector3{}
		} else if v.Team != p.Team {
			playerData.Private = true
			playerData.Position = util.Vector3{}
			playerData.HP = 0
			playerData.PrimaryWeapon = nil
			playerData.SecondaryWeapon = nil
			playerData.KnifeWeapon = nil
			playerData.EquippedWeapon = nil
			playerData.Money = 0
			playerData.Orientation = util.Vector3{}
		}

		players = append(players, playerData)
	}

	return players
}

func (e *SimpleEngine) PickupWeapon(userID uint64, weaponID int) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}
	for _, v := range e.droppedWeapons {
		if v.WeaponID != weaponID {
			return
		}

		if v.Location == p.Position {
			for _, w := range e.weapons {
				if w.ID == weaponID {
					switch w.WT {
					case weapon.Primary:
						p.PrimaryWeapon = &w
					case weapon.Secondary:
						p.SecondaryWeapon = &w
					case weapon.Knife:
						p.KnifeWeapon = &w
					}
				}
			}
		}
	}
}

func (e *SimpleEngine) Channel() <-chan interface{} {
	return e.ch
}

func (e *SimpleEngine) State() State {
	return e.state
}

func (e *SimpleEngine) DropWeapon(userID uint64) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}
	if p.EquippedWeapon == nil {
		return
	}

	e.droppedWeapons = append(e.droppedWeapons, SimpleWeapon{p.Position, p.EquippedWeapon.ID})
	if p.EquippedWeapon == p.PrimaryWeapon {
		p.PrimaryWeapon = nil
	} else if p.EquippedWeapon == p.SecondaryWeapon {
		p.SecondaryWeapon = nil
	} else if p.EquippedWeapon == p.KnifeWeapon {
		p.KnifeWeapon = nil
	}
	p.EquippedWeapon = nil
}

func (e *SimpleEngine) SwitchPrimaryWeapon(userID uint64) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}
	p.EquippedWeapon = p.PrimaryWeapon
}

func (e *SimpleEngine) SwitchSecondaryWeapon(userID uint64) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}
	p.EquippedWeapon = p.SecondaryWeapon
}

func (e *SimpleEngine) SwitchKnifeWeapon(userID uint64) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}
	p.EquippedWeapon = p.KnifeWeapon
}

func (e *SimpleEngine) Start() {
	e.Lock()
	defer e.Unlock()

	e.ch <- MatchStart{}
	e.ch <- RoundStart{}
}

func (e *SimpleEngine) Buy(userID uint64, weaponID int) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}

	var wp weapon.Weapon
	for _, v := range e.weapons {
		if v.ID == weaponID {
			wp = v
			break
		}
	}
	if wp.ID == 0 {
		return
	}

	if p.Money < int(wp.Price) {
		return
	} else {
		p.Money -= int(wp.Price)
		switch wp.WT {
		case weapon.Primary:
			p.PrimaryWeapon = &wp
		case weapon.Secondary:
			p.SecondaryWeapon = &wp
		case weapon.Knife:
			p.KnifeWeapon = &wp
		}
	}
}

func (e *SimpleEngine) SetOrientation(userID uint64, orientation util.Vector3) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}

	p.Orientation = orientation
}

func (e *SimpleEngine) Reload(userID uint64) {
	e.Lock()
	defer e.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}

	p.EquippedWeapon.Ammo = p.EquippedWeapon.AmmoMax
}

func (e *SimpleEngine) inLineOfSight(player, other Player) bool {
	return true
}
