package model

type Users struct {
	Id   int64   `gorm:"column:id;primaryKey" json:"id"`
	Name *string `gorm:"column:name" json:"name"`
}

func (Users) TableName() string {
	return "users"
}
