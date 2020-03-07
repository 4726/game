package player

import (
	"time"

	"github.com/4726/game/gameplay/5v5fps/util"
	"github.com/4726/game/gameplay/5v5fps/weapon"
)

type Player struct {
	UserID          uint64
	Team            util.TeamID
	Position        util.Vector3
	Orientation     util.Vector3
	Dead            bool
	LastUpdate      time.Time
	Connected       bool
	PrimaryWeapon   *weapon.Weapon
	SecondaryWeapon *weapon.Weapon
	KnifeWeapon     *weapon.Weapon
	EquippedWeapon  *weapon.Weapon
	UserScore       Score
	Money           int
}

func (p *Player) IsTeammate(other *Player) bool {
	return p.Team == other.Team
}
