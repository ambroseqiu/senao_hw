package model

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/ambroseqiu/senao_hw/repository"
	"github.com/ambroseqiu/senao_hw/util"
	"github.com/google/uuid"
)

var (
	maxFailedAttempt      = 5
	timeBlockLoginAttempt = time.Minute
)

var (
	ErrAccountRequestValidationFailed = errors.New("Create user request validation failed")
	ErrAccountIsAlreadyExisted        = errors.New("Account is already existed")
	ErrLoginAccountNotFound           = errors.New("Login account not found")
	ErrLoginWrongPassword             = errors.New("Wrong password")
	ErrLoginAttemptBlocked            = errors.New("too many failed login attempt, please try it later")
)

type UsecaseHandler interface {
	CreateAccount(ctx context.Context, req AccountRequest) (*AccountResponse, error)
	LoginAccount(ctx context.Context, req AccountRequest) (*AccountResponse, error)
}

type usecaseHandler struct {
	mu      sync.RWMutex
	loginAC map[string]LoginAttempt
	repo    repository.AccountRepository
}

func NewUsecaseHandler(repo repository.AccountRepository) UsecaseHandler {
	return &usecaseHandler{
		loginAC: make(map[string]LoginAttempt, 100),
		repo:    repo,
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
		if errors.Is(err, repository.ErrAccountIsDuplicated) {
			rsp.Reason = ErrAccountIsAlreadyExisted.Error()
			return rsp, ErrAccountIsAlreadyExisted
		}
		return nil, err
	}

	return rsp, nil
}

func (u *usecaseHandler) LoginAccount(ctx context.Context, req AccountRequest) (*AccountResponse, error) {
	rsp := &AccountResponse{
		Success: false,
		Reason:  "",
	}
	if err := u.loginValidate(req.Username); err != nil {
		rsp.Reason = err.Error()
		return rsp, err
	}
	account, err := u.repo.GetAccount(ctx, req.Username)
	if err != nil {
		if err == repository.ErrAccountRecordNotFound {
			rsp.Reason = ErrLoginAccountNotFound.Error()
			return rsp, ErrLoginAccountNotFound
		}
		return nil, err
	}
	if err = util.CheckPassword(req.Password, account.HashedPassword); err != nil {
		if err == util.ErrMismatchedPassword {
			u.AddFailedAttempt(account.Username)
			rsp.Reason = ErrLoginWrongPassword.Error()
			return rsp, ErrLoginWrongPassword
		}
		return nil, err
	}
	rsp.Success = true
	return rsp, nil
}

func (u *usecaseHandler) loginValidate(username string) error {
	defer u.mu.RUnlock()
	u.mu.RLock()
	loginAttempt, ok := u.loginAC[username]
	if !ok {
		return nil
	}
	if loginAttempt.FailedAttempt >= maxFailedAttempt {
		elapsedTime := time.Since(loginAttempt.LastTime)
		if elapsedTime < timeBlockLoginAttempt {
			return ErrLoginAttemptBlocked
		}
	}
	return nil
}

func (u *usecaseHandler) AddFailedAttempt(username string) {
	defer u.mu.Unlock()
	u.mu.Lock()
	loginAttempt, ok := u.loginAC[username]
	if !ok {
		u.loginAC[username] = LoginAttempt{
			FailedAttempt: 1,
			LastTime:      time.Now(),
		}
		return
	}
	loginAttempt.FailedAttempt++
	loginAttempt.LastTime = time.Now()
	u.loginAC[username] = LoginAttempt{
		FailedAttempt: loginAttempt.FailedAttempt,
		LastTime:      loginAttempt.LastTime,
	}
}
