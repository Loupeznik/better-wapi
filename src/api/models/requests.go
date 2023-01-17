package models

type SaveRowRequest struct {
	Subdomain  string `json:"subdomain"`
	IP         string `json:"ip"`
	Autocommit bool   `json:"autocommit"`
}
