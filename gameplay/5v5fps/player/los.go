package player

type LineOfSight interface {
	In(p, other *Player) bool
}