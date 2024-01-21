package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	apiModels "github.com/loupeznik/better-wapi/src/api/models"
	"github.com/loupeznik/better-wapi/src/helpers"
	"github.com/loupeznik/better-wapi/src/models"
)

type MockIntegrationService struct {
	config *models.Config
}

func NewMockIntegrationService(config *models.Config) *MockIntegrationService {
	return &MockIntegrationService{config: config}
}

func (s *MockIntegrationService) CreateRecord(domain string, request apiModels.SaveRowRequest) (int, error) {
	data := models.RequestData{
		Domain:    domain,
		Subdomain: request.Subdomain,
		TTL:       *request.TTL,
		Type:      string(*request.Type),
		IP:        request.Data,
	}

	status, _, err := s.makeRequest("dns-row-add", data, *request.Autocommit)

	return status, err
}

func (s *MockIntegrationService) UpdateRecord(domain string, request apiModels.SaveRowRequest) (int, error) {
	record, _, _ := s.GetRecord(domain, request.Subdomain)
	rowID, _ := strconv.Atoi(record.RecordID)

	data := models.RequestData{
		Domain: domain,
		RowID:  rowID,
		TTL:    *request.TTL,
		Type:   string(*request.Type),
		IP:     request.Data,
	}

	status, _, err := s.makeRequest("dns-row-update", data, *request.Autocommit)

	return status, err
}

func (s *MockIntegrationService) DeleteRecord(domain string, subdomain string, commit bool) (int, error) {
	record, _, _ := s.GetRecord(domain, subdomain)
	rowID, _ := strconv.Atoi(record.RecordID)

	data := models.RequestData{
		Domain: domain,
		RowID:  rowID}

	status, _, err := s.makeRequest("dns-row-delete", data, commit)

	return status, err
}

func (s *MockIntegrationService) GetInfo(domainName string) ([]models.Record, int, error) {
	data := models.RequestData{
		Domain: domainName}

	status, result, err := s.makeRequest("dns-rows-list", data, false)

	return result.Body.Data.Records, status, err
}

func (s *MockIntegrationService) GetRecord(domain string, subdomain string) (models.Record, int, error) {
	records, status, err := s.GetInfo(domain)
	var record models.Record

	if err != nil {
		return record, status, err
	}

	for _, row := range records {
		if row.Subdomain == subdomain {
			record = row
		}
	}

	if record.IP == "" {
		return record, 404, errors.New("not found")
	}

	return record, status, nil
}

func (s *MockIntegrationService) CommitChanges(domain string) (int, error) {
	data := models.RequestData{
		Subdomain: domain}

	status, _, err := s.makeRequest("dns-domain-commit", data, false)

	return status, err
}

func (s *MockIntegrationService) makeRequest(command string, data models.RequestData, commit bool) (int, models.WApiResponse, error) {
	token := helpers.GetApiToken(s.config.WApiUsername, s.config.WApiPassword)
	client := &http.Client{Timeout: time.Duration(60) * time.Second}

	request := &models.Request{Body: models.RequestBody{
		Login:   s.config.WApiUsername,
		Secret:  token,
		Command: command,
		Data:    data,
	}}

	var result models.WApiResponse

	response, err := client.Do(helpers.BuildRequest(s.baseUrl, request))
	if err != nil {
		helpers.Log("error", fmt.Sprintf("Request to WAPI failed: %s", err.Error()))
		return 500, result, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			helpers.Log("error", fmt.Sprintf("Failed to read response body: %s", err.Error()))
		}
	}(response.Body)

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return 500, result, err
	}

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return 500, result, err
	}

	if commit {
		s.CommitChanges(data.Domain)
	}

	status, err := helpers.ResolveStatusCode(result.Body.StatusCode)

	if err == nil && status >= 400 {
		err = errors.New(result.Body.ResultStatus)
	}

	return status, result, err
}
