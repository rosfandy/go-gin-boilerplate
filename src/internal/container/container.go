package container

import "gorm.io/gorm"

type Container struct {
	User *UserContainer
}

func NewContainer(db *gorm.DB) *Container {
	return &Container{
		User: NewUserContainer(db),
	}
}
