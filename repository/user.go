package repository

import (
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetApi(ctx context.Context) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetApi(ctx context.Context) error {
	return nil
}
