package util

type Vector3 struct {
	X, Y, Z float64
}

type TeamID bool

const (
	Team1 TeamID = true
	Team2 TeamID = false
)
