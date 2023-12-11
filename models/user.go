package models

type User struct {
	ID       int64
	Name     string
	Email    string
	Password string
	Username string
}

type GormUser struct {
	ID       int64 `gorm:"primary_key"`
	Name     string
	Email    string `gorm:"unique"`
	Password string
	Username string `gorm:"unique"`
}

func (User) TableName() string {
	return "gorm_users"
}
