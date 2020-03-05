package engine

import (
	"github.com/4726/game/gameplay/player"
	"github.com/4726/game/gameplay/util"
)

type Player struct {
	UserID uint64
	Position util.Vector3
	HP int
	Team util.Team
	PrimaryWeapon *weapon.Weapon
	SecondaryWeapon *weapon.Weapon
	KnifeWeapon *weapon.Weapon
	EquippedWeapon *weapon.Weapon
	Dead bool
	Private bool
	Kills int
	Deaths int
	Assists int
}

type State struct {
	Scores map[util.Team]int
}

type Engine interface {
	Init(Player)
	MoveLeft(userID uint64)
	MoveLeftUp(userID uint64)
	MoveLeftDown(userID uint64)
	MoveRight(userID uint64)
	MoveRightUp(userID uint64)
	MoveRightDown(userID uint64)
	MoveUp(userID uint64)
	MoveUpLeft(userID uint64)
	MoveUpRight(userID uint64)
	MoveDown(userID uint64)
	MoveDownLeft(userID uint64)
	MoveDownRight(userID uint64)
	Shoot(userID uint64, target player.Vector3)
	All(userID uint64) []Player
	PickupWeapon(userID uint64, weaponID int)
	Channel() chan interface{}
	State() State
	DropWeapon(userID uint64)
	SwitchPrimaryWeapon(userID uint64)
	SwitchSecondaryWeapon(userID uint64)
	SwitchKnifeWeapon(userID uint64)
}

type RoundEnd struct {
	WinningTeam util.Team
	LosingTeam util.Team	
}

type MatchEnd struct {
	WinningTeam util.Team
	LosingTeam util.Team
}