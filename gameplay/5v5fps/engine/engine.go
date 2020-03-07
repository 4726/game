package engine

import (
	"github.com/4726/game/gameplay/5v5fps/util"
	"github.com/4726/game/gameplay/5v5fps/weapon"
)

type Player struct {
	UserID          uint64
	Position        util.Vector3
	HP              int
	Team            util.TeamID
	PrimaryWeapon   *weapon.Weapon
	SecondaryWeapon *weapon.Weapon
	KnifeWeapon     *weapon.Weapon
	EquippedWeapon  *weapon.Weapon
	Dead            bool
	Private         bool
	Money           int
	Orientation     util.Vector3
	UserScore       Score
}

type Score struct {
	Kills, Deaths, Assists int
}

type State struct {
	Scores map[util.TeamID]int
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
	Shoot(userID uint64, target util.Vector3)
	All(userID uint64) []Player
	PickupWeapon(userID uint64, weaponID int)
	Channel() chan interface{}
	State() State
	DropWeapon(userID uint64)
	SwitchPrimaryWeapon(userID uint64)
	SwitchSecondaryWeapon(userID uint64)
	SwitchKnifeWeapon(userID uint64)
	Start()
	Buy(userID uint64, weaponID int)
	SetOrientation(userID uint64, orientation util.Vector3)
}

type RoundEnd struct {
	WinningTeam util.TeamID
	LosingTeam  util.TeamID
}

type MatchEnd struct {
	WinningTeam util.TeamID
	LosingTeam  util.TeamID
}

type RoundStart struct {
}

type MatchStart struct {
}
