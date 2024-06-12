package models

type ImageURLs struct {
	ID       uint   `gorm:"type:int;autoincrement;primary_key"`
	UserID   uint   `gorm:"type:int"`
	ImageURL string `gorm:"type:varchar(255)"`
}
