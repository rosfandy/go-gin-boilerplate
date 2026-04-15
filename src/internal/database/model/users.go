package model

type Users struct {
	Id    int64   `gorm:"column:id;primaryKey" json:"id"`
	Name  *string `gorm:"column:name" json:"name"`
	Email *string `gorm:"column:email" json:"email"`
}

func (Users) TableName() string {
	return "users"
}
