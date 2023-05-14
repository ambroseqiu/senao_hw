package model

import (
	"context"
	"errors"

	"github.com/ambroseqiu/senao_hw/repository"
	"github.com/ambroseqiu/senao_hw/util"
	"github.com/google/uuid"
)

var (
	ErrCreateUserRequestValidationFailed = errors.New("create user request validation failed")
)

type UsecaseHandler interface {
	CreateUser(ctx context.Context, req AccountRequest) (*AccountResponse, error)
}

type usecaseHandler struct {
	repo repository.UserRepository
}

func NewUsecaseHandler(repo repository.UserRepository) UsecaseHandler {
	return &usecaseHandler{
		repo: repo,
	}
}

func (u *usecaseHandler) CreateUser(ctx context.Context, req AccountRequest) (*AccountResponse, error) {
	rsp := &AccountResponse{
		Success: true,
		Reason:  "",
	}

	if err := req.Validate(); err != nil {
		rsp.Success = false
		rsp.Reason = err.Error()
		return rsp, ErrCreateUserRequestValidationFailed
	}

	uuid := uuid.New()
	hashedPassword, err := util.HashedPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user := &repository.User{
		ID:             uuid,
		Username:       req.Username,
		HashedPassword: hashedPassword,
	}

	if err := u.repo.CreateUser(ctx, user); err != nil {
		rsp.Success = false
		rsp.Reason = err.Error()
		return rsp, err
	}

	return rsp, nil
}
