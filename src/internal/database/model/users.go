package model

type Users struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}
