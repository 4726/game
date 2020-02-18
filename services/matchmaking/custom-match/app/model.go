package app

type Group struct {
	ID         int64
	Leader     uint64 `pg:unique`
	Name       string
	Password   string
	MaxUsers   uint32
	TotalUsers uint32
}

type User struct {
	ID      uint64
	GroupID int64
}
