package models

type User struct {
	ID    string `gorm:"primarykey"`
	Name  string `gorm:"type:varchar(255);not null"`
	Email string `gorm:"type:varchar(50);unique;not null"`
}
