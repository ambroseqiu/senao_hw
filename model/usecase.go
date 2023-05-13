package model

import (
	"context"

	"github.com/ambroseqiu/senao_hw/repository"
	"github.com/google/uuid"
)

type UsecaseHandler interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error)
}

type usecaseHandler struct {
	repo repository.UserRepository
}

func NewUsecaseHandler(repo repository.UserRepository) UsecaseHandler {
	return &usecaseHandler{
		repo: repo,
	}
}

func (u *usecaseHandler) CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
	uuid := uuid.New()
	user := &repository.User{
		ID:             uuid,
		Username:       req.Username,
		HashedPassword: "",
	}
	if err := u.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	rsp := &CreateUserResponse{
		Success: true,
		Reason:  "",
	}

	return rsp, nil
}
