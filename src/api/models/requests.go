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
	Subdomain  string      `json:"subdomain" default:""`
	Data       string      `json:"data" binding:"required"`
	Type       *RecordType `json:"type,omitempty" default:"A"`
	TTL        *int        `json:"ttl,omitempty" default:"3600"`
	Autocommit *bool       `json:"autocommit,omitempty" default:"false"`
}

type SaveRowRequestV2 struct {
	Data       string      `json:"data" binding:"required"`
	Type       *RecordType `json:"type,omitempty" default:"A"`
	TTL        *int        `json:"ttl,omitempty" default:"3600"`
	Autocommit *bool       `json:"autocommit,omitempty" default:"false"`
}

type DeleteRowRequest struct {
	Subdomain  string `json:"subdomain" binding:"required"`
	Autocommit *bool  `json:"autocommit,omitempty" default:"false"`
}

type DeleteRowRequestV2 struct {
	Autocommit *bool `json:"autocommit,omitempty" default:"false"`
}
