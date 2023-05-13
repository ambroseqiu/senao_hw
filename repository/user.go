package repository

import (
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}
