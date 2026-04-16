package container

import (
	"owner-api-proxy/internal/repository"
	"owner-api-proxy/internal/service"

	"gorm.io/gorm"
)

type UserContainer struct {
	UserRepository *repository.User
	UserService    *service.UserService
}

func NewUserContainer(db *gorm.DB) *UserContainer {
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)

	return &UserContainer{
		UserRepository: userRepository,
		UserService:    userService,
	}
}
