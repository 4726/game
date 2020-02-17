package app

type Group struct {
	ID int64
	Leader uint64 `pg:unique`
	Name string
	Password string
	MaxUsers uint32
	Users []uint64 `pg:,array`
}