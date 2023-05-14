package repository

import (
	"context"

	"github.com/pkg/errors"

	"gorm.io/gorm"
)

var (
	ErrUserIsAlreadyExisted = errors.New("Username already exists")
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
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrUserIsAlreadyExisted
		}
		return err
	}
	return nil
}
