package models

type Response_DTO struct {
	Success bool `json:"success"`
	Error   string `json:"error"`
}
