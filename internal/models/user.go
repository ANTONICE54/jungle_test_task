package models

type User struct {
	ID           uint   `gorm:"type:int;autoincrement;primary_key"`
	Username     string `gorm:"type:varchar(255);unique"`
	Email        string `gorm:"type:varchar(255);unique"`
	PasswordHash string `gorm:"type:varchar(255)"`
}
