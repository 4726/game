package engine

// SimpleEngine implements Engine
type SimpleEngine struct {
	players map[uint64]*Player
	playersLock sync.Mutex
	droppedWeapons []SimpleWeapon
	weapons []weapon.Weapon
	ch chan interface{}
	state State
}

type SimpleWeapon struct {
	Location util.Vector3
	WeaponID int
}

func (e *SimpleEngine) Init(p Player) {
	e.playersLock.Lock()
	defer e.playersLock.Unlock()

	e.players[p.UserID] = &p
}

func (e *SimpleEngine) MoveLeft(userID uint64) {
	e.playersLock.Lock()
	defer e.playersLock.Unlock()

	p, ok := e.players[p.UserID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X--, p.Position.Y, p.Position.Z}
}

func (e *SimpleEngine) MoveLeftUp(userID uint64) {
	e.playersLock.Lock()
	defer e.playersLock.Unlock()

	p, ok := e.players[p.UserID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X--, p.Position.Y++, p.Position.Z}
}

func (e *SimpleEngine) MoveLeftDown(userID uint64) {
	e.playersLock.Lock()
	defer e.playersLock.Unlock()

	p, ok := e.players[p.UserID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X--, p.Position.Y--, p.Position.Z}
}

func (e *SimpleEngine) MoveRight(userID uint64) {
	e.playersLock.Lock()
	defer e.playersLock.Unlock()

	p, ok := e.players[p.UserID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X++, p.Position.Y, p.Position.Z}
}

func (e *SimpleEngine) MoveRightUp(userID uint64) {
	e.playersLock.Lock()
	defer e.playersLock.Unlock()

	p, ok := e.players[p.UserID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X++, p.Position.Y++, p.Position.Z}
}

func (e *SimpleEngine) MoveRightDown(userID uint64) {
	e.playersLock.Lock()
	defer e.playersLock.Unlock()

	p, ok := e.players[p.UserID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X++, p.Position.Y--, p.Position.Z}
}

func (e *SimpleEngine) MoveUp(userID uint64) {
	e.playersLock.Lock()
	defer e.playersLock.Unlock()

	p, ok := e.players[p.UserID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X, p.Position.Y++, p.Position.Z}
}

func (e *SimpleEngine) MoveUpLeft(userID uint64) {
	e.playersLock.Lock()
	defer e.playersLock.Unlock()

	p, ok := e.players[p.UserID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X--, p.Position.Y++, p.Position.Z}
}

func (e *SimpleEngine) MoveUpRight(userID uint64) {
	e.playersLock.Lock()
	defer e.playersLock.Unlock()

	p, ok := e.players[p.UserID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X++, p.Position.Y++, p.Position.Z}
}

func (e *SimpleEngine) MoveDown(userID uint64) {
	e.playersLock.Lock()
	defer e.playersLock.Unlock()

	p, ok := e.players[p.UserID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X, p.Position.Y--, p.Position.Z}
}

func (e *SimpleEngine) MoveDownLeft(userID uint64) {
	e.playersLock.Lock()
	defer e.playersLock.Unlock()

	p, ok := e.players[p.UserID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X--, p.Position.Y--, p.Position.Z}
}

func (e *SimpleEngine) MoveDownRight(userID uint64) {
	e.playersLock.Lock()
	defer e.playersLock.Unlock()

	p, ok := e.players[p.UserID]
	if !ok {
		return
	}

	p.Position = util.Vector3{p.Position.X++, p.Position.Y--, p.Position.Z}
}

func (e *SimpleEngine) Shoot(userID uint64, target player.Vector3) {
	e.playersLock.Lock()
	defer e.playersLock.Unlock()
	
	p, ok := e.players[userID]
	if !ok {
		return
	}

	var haveDeath int
	for _, v := range e.players {
		if v.UserID == p.UserID {
			continue
		}

		if v.Position != target {
			continue
		}

		v.HP -= p.EquippedWeapon.Damage
		if v.HP <= 0 {
			v.Dead = true
			haveDeath = true
		}
	}

	if !haveDeath {
		return
	}

	teamDeathCount := map[util.Team]int{}
	for _, v := range e.players {
		if v.Dead {
			count, ok := totalDeathCount[v.Team]
			if !ok {
				totalDeathCount[v.Team] = 1
			} else {
				totalDeathCount[v.Team] = count++
			}
		}
	}

	for k, v := range teamDeathCount {
		var winner util.Team
		var loser util.Team
		if v == 5 {
			loser = k
			if loser == util.Team1 {
				winner = util.Team2
			} else {
				winner = util.Team1
			}
			e.ch <- RoundEnd{
				WinningTeam: winner,
				LosingTeam: loser,
			}

			score := e.state.Score[winner]
			e.state.Score[winner] = score++

			if e.state.Score[winner] == 16 {
				e.ch <- MatchEnd{
					WinningTeam: winner,
					LosingTeam: loser,
				}
			}
			break
		}
	}
}

func (e *SimpleEngine) All(userID uint64) []Player {
	e.playersLock.Lock()
	defer e.playersLock.Unlock()

	var players []Player
	p, ok := e.players
	if !ok {
		return players
	}
	for _, v := range e.players {
		playerData := *v

		//doesn't return information that should be private
		if v.Team != p.Team {
			playerData.Private = true
			playerData.Position = util.Vector3{}
			playerData.HP = 0
			playerData.PrimaryWeapon = nil
			playerData.SecondaryWeapon = nil
			playerData.KnifeWeapon = nil
			playerData.EquippedWeapon = nil
		}

		if e.inLineOfSight(p, v) {
			playerData.Private = false
			playerData.Position = v.Position
			playerData.EquippedWeapon = v.EquippedWeapon
		}

		players = append(players, playerData)
	}

	return players
}

func (e *SimpleEngine) PickupWeapon(userID uint64, weaponID int) {
	e.playersLock.Lock()
	defer e.playersLock.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}
	for _, v := range e.droppedWeapons {
		if v.WeaponID != weaponID {
			return
		}
		
		if v.Location == p.Position {
			for _, w := range weapons {
				if w.ID == weaponID {
					switch w.WT {
					case weapon.Primary:
						p.PrimaryWeapon = w
					case weapon.Secondary:
						p.SecondWeapon = w
					case weapon.Knife:
						p.KnifeWeapon = w
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
	e.playersLock.Lock()
	defer e.playersLock.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}
	if p.EquippedWeapon == nil {
		return 
	}

	s.DroppedWeapons = append(s.DroppedWeapons, SimpleWeapon{p.Location, p.EquippedWeapon})
	if s.EquippedWeapon == s.PrimaryWeapon {
		s.PrimaryWeapon = nil
	} else if s.EquippedWeapon == s.SecondaryWeapon {
		s.SecondaryWeapon = nil
	} else if s.EquippedWeapon == s.KnifeWeeapon {
		s.KnifeWeapon = nil
	}
	s.EquippedWeapon = nil
}

func (e *SimpleEngine) SwitchPrimaryWeapon(userID uint64) {
	e.playersLock.Lock()
	defer e.playersLock.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}
	p.EquippedWeapon = p.PrimaryWeapon
}

func (e *SimpleEngine) SwitchSecondaryWeapon(userID uint64) {
	e.playersLock.Lock()
	defer e.playersLock.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}
	p.EquippedWeapon = p.SecondaryWeapon
}

func (e *SimpleEngine) SwitchKnifeWeapon(userID uint64) {
	e.playersLock.Lock()
	defer e.playersLock.Unlock()

	p, ok := e.players[userID]
	if !ok {
		return
	}
	p.EquippedWeapon = p.KnifeWeapon
}

func (e *SimpleEngine) inLineOfSight(player, other Player) bool {
	return true
}