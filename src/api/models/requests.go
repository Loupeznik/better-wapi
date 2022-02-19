package models

type UpdateRequest struct {
	Domain string `json:"domain"`
	IP     string `json:"ip"`
}

type CreateRequest struct {
	Domain    string `json:"domain"`
	Subdomain string `json:"subdomain"`
	IP        string `json:"ip"`
}
