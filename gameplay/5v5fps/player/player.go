package player

import (
	"github.com/4726/game/gameplay/weapon"
)

type Vector3 struct {
	X, Y, Z float64
}

type TeamID bool

const (
	Team1 TeamID = true
	Team2 TeamID = false
)

type Player struct {
	UserID uint64
	Team TeamID
	Position Vector3
	Orientation Vecotr3
	Dead bool
	LastUpdate time.Time
	Connected bool
	PrimaryWeapon weapon.Weapon
	SecondaryWeapon weapon.Weapon
	KnifeWeapon weapon.Weapon
	EquippedWeapon weapon.Weapon
	LOS LineOfSight
	UserScore Score
	Money int
}

func (p *Player) IsTeammate(other *Player) bool {
	return p.Team == other.Team 
}

func (p *Player) InLineOfSight(other *Player) bool {
	return p.LOS.In(p, other)
}