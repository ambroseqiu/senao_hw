package model

import (
	"context"

	"github.com/ambroseqiu/senao_hw/repository"
)

type UsecaseHandler interface {
	GetApi(ctx context.Context) error
}

type usecaseHandler struct {
	repo repository.UserRepository
}

func NewUsecaseHandler(repo repository.UserRepository) UsecaseHandler {
	return &usecaseHandler{
		repo: repo,
	}
}

func (u *usecaseHandler) GetApi(ctx context.Context) error {
	return u.repo.GetApi(ctx)
}
