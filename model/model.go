package model

import (
	"errors"
	"regexp"
)

const (
	maxUsernameLength = 32
	minUsernameLength = 3
	maxPasswordLength = 32
	minPasswordLength = 8
)

var (
	isValidateFormat            = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
	ErrInvalidUsernameFormat    = errors.New("Invalid username format, only alphabets and numbers are allowed")
	ErrUsernameIsTooLarge       = errors.New("Username is too large")
	ErrUsernameIsTooShort       = errors.New("Username is too short")
	ErrInvalidPasswordFormat    = errors.New("Invalid password format, only alphabets and numbers are allowed")
	ErrPasswordIsTooLarge       = errors.New("Password is too large")
	ErrPasswordIsTooShort       = errors.New("Password is too short")
	ErrPasswordValidationFailed = errors.New("Invalid password format. It should contain at least 1 uppercase letter, 1 lowercase letter, and 1 number")
)

type AccountRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AccountResponse struct {
	Success bool   `json:"success" binding:"required"`
	Reason  string `json:"reason" binding:"required"`
}

func (req *AccountRequest) Validate() error {
	if !isValidateFormat(req.Username) {
		return ErrInvalidUsernameFormat
	}
	if !isValidateFormat(req.Password) {
		return ErrInvalidPasswordFormat
	}
	if len(req.Username) > maxUsernameLength {
		return ErrUsernameIsTooLarge
	}
	if len(req.Username) < minUsernameLength {
		return ErrUsernameIsTooShort
	}
	if len(req.Password) > maxPasswordLength {
		return ErrPasswordIsTooLarge
	}
	if len(req.Password) < minPasswordLength {
		return ErrPasswordIsTooShort
	}
	return validPassword(req.Password)
}

func validPassword(password string) error {
	hasUppercase := false
	hasLowercase := false
	hasNumber := false

	for _, ch := range password {
		if ch >= 'A' && ch <= 'Z' {
			hasUppercase = true
		}
		if ch >= 'a' && ch <= 'z' {
			hasLowercase = true
		}
		if ch >= '0' && ch <= '9' {
			hasNumber = true
		}
	}
	if !hasLowercase || !hasUppercase || !hasNumber {
		return ErrPasswordValidationFailed
	}
	return nil
}
