package models

type UpdateRequest struct {
	Domain string `json:"domain"`
	IP     string `json:"ip"`
}
