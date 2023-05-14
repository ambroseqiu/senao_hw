package util

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var ErrMismatchedPassword = errors.New("Password is not correct")

func HashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Failed to hashed password: %v", err)
	}
	return string(hashedPassword), nil
}

func CheckPassword(password string, hashedPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return ErrMismatchedPassword
		}
		return err
	}
	return nil
}
