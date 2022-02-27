package models

type WApiResponse struct {
	Body ResponseBody `json:"response"`
}

type ResponseBody struct {
	StatusCode   int         `json:"code"`
	ResultStatus string      `json:"result"`
	Timestamp    int64       `json:"timestamp,omitempty"`
	RequestID    string      `json:"svTRID,omitempty"`
	Command      string      `json:"command,omitempty"`
	Data         ResponseRow `json:"data,omitempty"`
}

type ResponseRow struct {
	Records []Record `json:"row,omitempty"`
}

type Record struct {
	RecordID  string `json:"ID,omitempty"`
	Subdomain string `json:"name,omitempty"`
	TTL       string `json:"ttl,omitempty"`
	Type      string `json:"rdtype,omitempty"`
	IP        string `json:"rdata,omitempty"`
	UpdatedAt string `json:"changed_date,omitempty"`
	Comment   string `json:"author_comment,omitempty"`
}
