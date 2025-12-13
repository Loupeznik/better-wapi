package services

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	apiModels "github.com/loupeznik/better-wapi/src/api/models"
	"github.com/loupeznik/better-wapi/src/models"
)

type MockService struct {
	domains map[string]*DomainData
	mu      sync.RWMutex
	nextID  int
}

type DomainData struct {
	Domain  models.Domain
	Records []models.Record
}

func NewMockService() *MockService {
	service := &MockService{
		domains: make(map[string]*DomainData),
		nextID:  1,
	}

	service.initializeMockData()

	return service
}

func (s *MockService) initializeMockData() {
	domains := []struct {
		name    string
		status  string
		records []struct {
			subdomain string
			ttl       string
			recordType string
			data      string
		}
	}{
		{
			name:   "example.com",
			status: "active",
			records: []struct {
				subdomain string
				ttl       string
				recordType string
				data      string
			}{
				{"@", "3600", "A", "192.0.2.1"},
				{"www", "3600", "A", "192.0.2.1"},
				{"mail", "3600", "A", "192.0.2.10"},
				{"@", "3600", "MX", "10 mail.example.com"},
				{"@", "3600", "TXT", "v=spf1 mx -all"},
			},
		},
		{
			name:   "demo.net",
			status: "active",
			records: []struct {
				subdomain string
				ttl       string
				recordType string
				data      string
			}{
				{"@", "3600", "A", "198.51.100.1"},
				{"www", "3600", "A", "198.51.100.1"},
				{"ftp", "3600", "A", "198.51.100.5"},
				{"@", "7200", "AAAA", "2001:db8::1"},
			},
		},
		{
			name:   "test.org",
			status: "active",
			records: []struct {
				subdomain string
				ttl       string
				recordType string
				data      string
			}{
				{"@", "3600", "A", "203.0.113.1"},
				{"api", "3600", "A", "203.0.113.10"},
				{"cdn", "3600", "CNAME", "cdn.cloudprovider.com"},
			},
		},
	}

	for _, domain := range domains {
		domainData := &DomainData{
			Domain: models.Domain{
				Name:   domain.name,
				Status: domain.status,
			},
			Records: make([]models.Record, 0),
		}

		for _, record := range domain.records {
			domainData.Records = append(domainData.Records, models.Record{
				RecordID:  strconv.Itoa(s.nextID),
				Subdomain: record.subdomain,
				TTL:       record.ttl,
				Type:      record.recordType,
				IP:        record.data,
				UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
				Comment:   "Mock record",
			})
			s.nextID++
		}

		s.domains[domain.name] = domainData
	}
}

func (s *MockService) CreateRecord(domain string, request apiModels.SaveRowRequest) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	domainData, exists := s.domains[domain]
	if !exists {
		return 404, errors.New("domain not found")
	}

	for _, record := range domainData.Records {
		if record.Subdomain == request.Subdomain && record.Type == string(*request.Type) {
			return 409, errors.New("record already exists")
		}
	}

	newRecord := models.Record{
		RecordID:  strconv.Itoa(s.nextID),
		Subdomain: request.Subdomain,
		TTL:       strconv.Itoa(*request.TTL),
		Type:      string(*request.Type),
		IP:        request.Data,
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
		Comment:   "Created via mock API",
	}

	s.nextID++
	domainData.Records = append(domainData.Records, newRecord)

	return 201, nil
}

func (s *MockService) UpdateRecord(domain string, request apiModels.SaveRowRequest) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	domainData, exists := s.domains[domain]
	if !exists {
		return 404, errors.New("domain not found")
	}

	for i, record := range domainData.Records {
		if record.Subdomain == request.Subdomain {
			domainData.Records[i].TTL = strconv.Itoa(*request.TTL)
			domainData.Records[i].Type = string(*request.Type)
			domainData.Records[i].IP = request.Data
			domainData.Records[i].UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
			return 200, nil
		}
	}

	return 404, errors.New("record not found")
}

func (s *MockService) UpdateRecordV2(domain string, rowID int, request apiModels.SaveRowRequestV2) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	domainData, exists := s.domains[domain]
	if !exists {
		return 404, errors.New("domain not found")
	}

	recordIDStr := strconv.Itoa(rowID)
	for i, record := range domainData.Records {
		if record.RecordID == recordIDStr {
			domainData.Records[i].TTL = strconv.Itoa(*request.TTL)
			domainData.Records[i].Type = string(*request.Type)
			domainData.Records[i].IP = request.Data
			domainData.Records[i].UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
			return 204, nil
		}
	}

	return 404, errors.New("record not found")
}

func (s *MockService) DeleteRecord(domain string, subdomain string, commit bool) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	domainData, exists := s.domains[domain]
	if !exists {
		return 404, errors.New("domain not found")
	}

	for i, record := range domainData.Records {
		if record.Subdomain == subdomain {
			domainData.Records = append(domainData.Records[:i], domainData.Records[i+1:]...)
			return 200, nil
		}
	}

	return 404, errors.New("record not found")
}

func (s *MockService) DeleteRecordV2(domain string, rowID int, commit bool) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	domainData, exists := s.domains[domain]
	if !exists {
		return 404, errors.New("domain not found")
	}

	recordIDStr := strconv.Itoa(rowID)
	for i, record := range domainData.Records {
		if record.RecordID == recordIDStr {
			domainData.Records = append(domainData.Records[:i], domainData.Records[i+1:]...)
			return 204, nil
		}
	}

	return 404, errors.New("record not found")
}

func (s *MockService) GetInfo(domainName string) ([]models.Record, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	domainData, exists := s.domains[domainName]
	if !exists {
		return nil, 404, errors.New("domain not found")
	}

	return domainData.Records, 200, nil
}

func (s *MockService) GetRecord(domain string, subdomain string) (models.Record, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	domainData, exists := s.domains[domain]
	if !exists {
		return models.Record{}, 404, errors.New("domain not found")
	}

	for _, record := range domainData.Records {
		if record.Subdomain == subdomain {
			return record, 200, nil
		}
	}

	return models.Record{}, 404, errors.New("record not found")
}

func (s *MockService) CommitChanges(domain string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.domains[domain]; !exists {
		return 404, errors.New("domain not found")
	}

	fmt.Printf("Mock: DNS changes committed for domain %s\n", domain)
	return 200, nil
}

func (s *MockService) ListDomains() ([]models.Domain, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	domains := make([]models.Domain, 0, len(s.domains))
	for _, domainData := range s.domains {
		domains = append(domains, domainData.Domain)
	}

	return domains, 200, nil
}
