package repository

import (
	"context"

	"owner-api-proxy/internal/database/model"

	"gorm.io/gorm"
)

type User struct {
	client *gorm.DB
}

func NewUserRepository(client *gorm.DB) *User {
	return &User{client: client}
}

func (r *User) GetUsers(ctx context.Context) ([]model.Users, error) {
	var users []model.Users

	err := r.client.WithContext(ctx).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
