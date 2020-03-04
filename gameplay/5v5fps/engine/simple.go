package engine

// SimpleEngine implements Engine
type SimpleEngine struct {
	players map[uint64]*Player
	playersLock sync.Mutex
	droppedWeapons []SimpleWeapon
	weapons []weapon.Weapon
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
		if v.UserID == p.UserID {
			players = append(players, v)
		}
		if v.Team == p.Team {
			players = append(players, v)
		}
		if e.inLineOfSight(p, v) {
			players = append(players, v)
		}
	}

	return players
}

func (e *SimpleEngine) PickupWeapon(userID uint64, weaponID int) {
	e.playersLock.Lock()
	defer e.playersLock.Unlock()

	p, ok := e.players
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

func (e *SimpleEngine) inLineOfSight(player, other Player) bool {
	return true
}