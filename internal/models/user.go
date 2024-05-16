package models

type User struct {
	Id    string `gorm:"primaryKey"`
	Name  string `gorm:"type:varchar(255);not null"`
	Email string `gorm:"type:varchar(50);unique;not null"`
}
