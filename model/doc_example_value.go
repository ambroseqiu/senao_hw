package model

type DocResponseSuccess struct {
	Success bool   `json:"success" example:"true"`
	Reason  string `json:"reason" example:""`
}

type DocResponseBadRequest struct {
	Success bool   `json:"success" example:"false"`
	Reason  string `json:"reason" example:"Password is too short"`
}
