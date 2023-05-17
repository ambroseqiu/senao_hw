package model

type DocResponseSuccess struct {
	Success bool   `json:"success" example:"true"`
	Reason  string `json:"reason" example:""`
}

type DocResponseAlreadyExisted struct {
	Success bool   `json:"success" example:"false"`
	Reason  string `json:"reason" example:"Account is already existed"`
}

type DocResponseBadRequest struct {
	Success bool   `json:"success" example:"false"`
	Reason  string `json:"reason" example:"Password is too short"`
}

type DocResponseAccountNotFound struct {
	Success bool   `json:"success" example:"false"`
	Reason  string `json:"reason" example:"Login account not found"`
}

type DocResponseWrongPassword struct {
	Success bool   `json:"success" example:"false"`
	Reason  string `json:"reason" example:"Wrong password"`
}

type DocResponseTooManyRequest struct {
	Success bool   `json:"success" example:"false"`
	Reason  string `json:"reason" example:"too many failed login attempt, please try it later"`
}
