package services

import (
	apiModels "github.com/loupeznik/better-wapi/src/api/models"
	"github.com/loupeznik/better-wapi/src/models"
)

type DNSService interface {
	CreateRecord(domain string, request apiModels.SaveRowRequest) (int, error)
	UpdateRecord(domain string, request apiModels.SaveRowRequest) (int, error)
	UpdateRecordV2(domain string, rowID int, request apiModels.SaveRowRequestV2) (int, error)
	DeleteRecord(domain string, subdomain string, commit bool) (int, error)
	DeleteRecordV2(domain string, rowID int, commit bool) (int, error)
	GetInfo(domainName string) ([]models.Record, int, error)
	GetRecord(domain string, subdomain string) (models.Record, int, error)
	CommitChanges(domain string) (int, error)
	ListDomains() ([]models.Domain, int, error)
}
