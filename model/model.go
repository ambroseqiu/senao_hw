package model

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8,max=32,regex=^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9]).*$"`
}

type CreateUserResponse struct {
	Success bool   `json:"success"`
	Reason  string `json:"reason"`
}
