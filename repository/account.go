package repository

import (
	"context"

	"github.com/pkg/errors"

	"gorm.io/gorm"
)

var (
	ErrAccountIsDuplicated   = errors.New("Account is duplicated")
	ErrAccountRecordNotFound = errors.New("Account is not found")
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, account *Account) error
	GetAccount(ctx context.Context, username string) (*Account, error)
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{
		db: db,
	}
}

func (r *accountRepository) CreateAccount(ctx context.Context, account *Account) error {
	if err := r.db.Create(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrAccountIsDuplicated
		}
		return err
	}
	return nil
}

func (r *accountRepository) GetAccount(ctx context.Context, username string) (*Account, error) {
	account := &Account{}
	if err := r.db.Where("username = ?", username).First(account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAccountRecordNotFound
		}
		return nil, err
	}
	return account, nil
}
