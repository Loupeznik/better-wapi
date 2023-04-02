package models

type RecordType string

const (
	A     RecordType = "A"
	CNAME RecordType = "CNAME"
	MX    RecordType = "MX"
	TXT   RecordType = "TXT"
	NS    RecordType = "NS"
	SRV   RecordType = "SRV"
	AAAA  RecordType = "AAAA"
	CAA   RecordType = "CAA"
	NAPTR RecordType = "NAPTR"
	TLSA  RecordType = "TLSA"
	SSHFP RecordType = "SSHFP"
)

type SaveRowRequest struct {
	Subdomain  string      `json:"subdomain" binding:"required"`
	Data       string      `json:"data"`
	Type       *RecordType `json:"type,omitempty" default:"A"`
	TTL        *int        `json:"ttl,omitempty" default:"3600"`
	Autocommit *bool       `json:"autocommit,omitempty" default:"false"`
}
