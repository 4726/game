package app

type Group struct {
	ID       int64 `gorm:"AUTO_INCREMENT;NOT NULL;PRIMARY_KEY"`
	Leader   uint64
	Name     string
	Password string
	MaxUsers uint32
	Users    []User `gorm:"foreignkey:GroupID"`
}

type User struct {
	UserID  uint64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT:false"`
	GroupID int64
}
