package services

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/loupeznik/better-wapi/src/helpers"
	"github.com/loupeznik/better-wapi/src/models"
)

type IntegrationService struct {
	config  *models.Config
	baseUrl string
}

func NewIntegrationService(config *models.Config) *IntegrationService {
	wapiBaseUrl := "https://api.wedos.com/wapi/json"
	return &IntegrationService{config: config, baseUrl: wapiBaseUrl}
}

func (s *IntegrationService) CreateRecord(domain string, subdomain string, ip string, commit bool) (int, error) {
	data := models.RequestData{
		Domain:    domain,
		Subdomain: subdomain,
		TTL:       1800,
		Type:      "A",
		IP:        ip,
	}

	status, _, err := s.makeRequest("dns-row-add", data, commit)

	return status, err
}

func (s *IntegrationService) UpdateRecord(domain string, subdomain string, newIp string, commit bool) (int, error) {
	record, _, _ := s.GetRecord(domain, subdomain)
	rowID, _ := strconv.Atoi(record.RecordID)

	data := models.RequestData{
		Domain: domain,
		RowID:  rowID,
		TTL:    1800,
		Type:   "A",
		IP:     newIp,
	}

	status, _, err := s.makeRequest("dns-row-update", data, commit)

	return status, err
}

func (s *IntegrationService) DeleteRecord(domain string, subdomain string, commit bool) (int, error) {
	record, _, _ := s.GetRecord(domain, subdomain)
	rowID, _ := strconv.Atoi(record.RecordID)

	data := models.RequestData{
		Domain: domain,
		RowID:  rowID}

	status, _, err := s.makeRequest("dns-row-delete", data, commit)

	return status, err
}

func (s *IntegrationService) GetInfo(domainName string) ([]models.Record, int, error) {
	data := models.RequestData{
		Domain: domainName}

	status, result, err := s.makeRequest("dns-rows-list", data, false)

	return result.Body.Data.Records, status, err
}

func (s *IntegrationService) GetRecord(domain string, subdomain string) (models.Record, int, error) {
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

func (s *IntegrationService) CommitChanges(domain string) (int, error) {
	data := models.RequestData{
		Subdomain: domain}

	status, _, err := s.makeRequest("dns-domain-commit", data, false)

	return status, err
}

func getApiToken(username string, password string) string {
	passwordHash := sha1.New()
	passwordHash.Write([]byte(password))
	location, _ := time.LoadLocation("Europe/Prague")
	hour := formatHour(time.Now().In(location).Hour())

	passwordHashString := fmt.Sprintf("%x", passwordHash.Sum(nil))

	token := fmt.Sprintf("%s%s%s", username, passwordHashString, hour)

	tokenHash := sha1.New()
	tokenHash.Write([]byte(token))

	return fmt.Sprintf("%x", tokenHash.Sum(nil))
}

func formatHour(hour int) string {
	var formattedHour string

	if hour < 10 {
		formattedHour = fmt.Sprintf("0%d", hour)
	} else {
		formattedHour = fmt.Sprintf("%d", hour)
	}

	return formattedHour
}

func (s *IntegrationService) makeRequest(command string, data models.RequestData, commit bool) (int, models.WApiResponse, error) {
	token := getApiToken(s.config.WApiUsername, s.config.WApiPassword)
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
