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
	All(userID uint64) []Player //returns all teamates and players in line of sight
	PickupWeapon(userID uint64, weaponID int)
}