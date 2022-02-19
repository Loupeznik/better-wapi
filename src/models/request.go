package models

type Request struct {
	Body RequestBody `json:"request"`
}

type RequestBody struct {
	Login   string      `json:"user"`
	Secret  string      `json:"auth"`
	Command string      `json:"command"`
	Data    RequestData `json:"data"`
}

type RequestData struct {
	Domain    string `json:"domain"`
	Subdomain string `json:"name,omitempty"`
	TTL       int    `json:"ttl,omitempty"`
	Type      string `json:"type,omitempty"`
	RowID     int    `json:"row_id,omitempty"`
	IP        string `json:"rdata,omitempty"`
}
