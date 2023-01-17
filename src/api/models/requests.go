package models

type SaveRowRequest struct {
	Subdomain  string `json:"subdomain" validate:"required"`
	IP         string `json:"ip" validate:"optional"`
	Autocommit bool   `json:"autocommit" validate:"optional" default:"false"`
}
