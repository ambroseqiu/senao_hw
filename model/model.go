package model

import "time"

type AccountRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AccountResponse struct {
	Success bool   `json:"success" binding:"required"`
	Reason  string `json:"reason" binding:"required"`
}

type LoginAttempt struct {
	FailedAttempt int
	LastTime      time.Time
}
