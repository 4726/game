package app

type Request struct {
	ID       int64 `gorm:"AUTO_INCREMENT;NOT NULL;PRIMARY_KEY"`
	From     uint64
	To       uint64
	Accepted bool
}
