package model

import (
	"context"
	"errors"

	"github.com/ambroseqiu/senao_hw/repository"
	"github.com/ambroseqiu/senao_hw/util"
	"github.com/google/uuid"
)

var (
	ErrAccountRequestValidationFailed = errors.New("Create user request validation failed")
	ErrLoginAccountNotFound           = errors.New("Login account not found")
	ErrLoginAccountNotAllowed         = errors.New("Login access not allowed")
)

type UsecaseHandler interface {
	CreateAccount(ctx context.Context, req AccountRequest) (*AccountResponse, error)
	LoginAccount(ctx context.Context, req AccountRequest) (*AccountResponse, error)
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
		return rsp, ErrAccountRequestValidationFailed
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

func (u *usecaseHandler) LoginAccount(ctx context.Context, req AccountRequest) (*AccountResponse, error) {
	rsp := &AccountResponse{
		Success: true,
		Reason:  "",
	}

	if err := req.Validate(); err != nil {
		rsp.Success = false
		rsp.Reason = err.Error()
		return rsp, ErrAccountRequestValidationFailed
	}

	account, err := u.repo.GetAccount(ctx, req.Username)
	if err != nil {
		rsp.Success = false
		if err == repository.ErrAccountRecordNotFound {
			rsp.Reason = err.Error()
			return rsp, ErrLoginAccountNotFound
		}
		return nil, err
	}
	if err = util.CheckPassword(req.Password, account.HashedPassword); err != nil {
		rsp.Success = false
		if err == util.ErrMismatchedPassword {
			rsp.Reason = err.Error()
			return rsp, ErrLoginAccountNotAllowed
		}
		return nil, err
	}
	return rsp, nil
}
