package model

import (
	"context"
	"errors"

	"github.com/ambroseqiu/senao_hw/repository"
	"github.com/ambroseqiu/senao_hw/util"
	"github.com/google/uuid"
)

var (
	ErrCreateAccountRequestValidationFailed = errors.New("create user request validation failed")
)

type UsecaseHandler interface {
	CreateAccount(ctx context.Context, req AccountRequest) (*AccountResponse, error)
}

type usecaseHandler struct {
	repo repository.AccountRepository
}

func NewUsecaseHandler(repo repository.AccountRepository) UsecaseHandler {
	return &usecaseHandler{
		repo: repo,
	}
}

func (u *usecaseHandler) CreateAccount(ctx context.Context, req AccountRequest) (*AccountResponse, error) {
	rsp := &AccountResponse{
		Success: true,
		Reason:  "",
	}

	if err := req.Validate(); err != nil {
		rsp.Success = false
		rsp.Reason = err.Error()
		return rsp, ErrCreateAccountRequestValidationFailed
	}

	uuid := uuid.New()
	hashedPassword, err := util.HashedPassword(req.Password)
	if err != nil {
		return nil, err
	}
	account := &repository.Account{
		ID:             uuid,
		Username:       req.Username,
		HashedPassword: hashedPassword,
	}

	if err := u.repo.CreateAccount(ctx, account); err != nil {
		rsp.Success = false
		rsp.Reason = err.Error()
		return rsp, err
	}

	return rsp, nil
}

// func (u *usecaseHandler) LoginUser(req AccountRequest){}
