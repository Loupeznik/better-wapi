package services

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/loupeznik/better-wapi/src/helpers"
	"github.com/loupeznik/better-wapi/src/models"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type IntegrationService interface {
	CreateRecord()
	UpdateRecord()
	DeleteRecord()
	GetInfo()
	GetRecord()
}

type integrationService struct {
	config  *models.Config
	baseUrl string
}

func NewIntegrationService(config *models.Config) *integrationService {
	wapiBaseUrl := "https://api.wedos.com/wapi/json"
	return &integrationService{config: config, baseUrl: wapiBaseUrl}
}

func (s *integrationService) CreateRecord(domain string, subdomain string, ip string) string {
	token := getApiToken(s.config.WApiUsername, s.config.WApiPassword)
	client := &http.Client{Timeout: time.Duration(60) * time.Second}
	request := &models.Request{Body: models.RequestBody{
		Login:   s.config.WApiUsername,
		Secret:  token,
		Command: "dns-row-add",
		Data: models.RequestData{
			Domain:    domain,
			Subdomain: subdomain,
			TTL:       1800,
			Type:      "A",
			IP:        ip},
	}}

	response, err := client.Do(helpers.BuildRequest(s.baseUrl, request))

	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		return string(bodyBytes)
	}

	return response.Status
}

func (s *integrationService) UpdateRecord(domain string, subdomain string, newIp string) string {
	token := getApiToken(s.config.WApiUsername, s.config.WApiPassword)
	client := &http.Client{Timeout: time.Duration(60) * time.Second}

	rowID, _ := strconv.Atoi(s.GetRecord(domain, subdomain).RecordID)

	request := &models.Request{Body: models.RequestBody{
		Login:   s.config.WApiUsername,
		Secret:  token,
		Command: "dns-row-update",
		Data: models.RequestData{
			Domain: domain,
			RowID:  rowID,
			TTL:    1800,
			Type:   "A",
			IP:     newIp},
	}}

	response, err := client.Do(helpers.BuildRequest(s.baseUrl, request))

	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		return string(bodyBytes)
	}

	return response.Status
}

func (s *integrationService) DeleteRecord(domain string, subdomain string) string {
	token := getApiToken(s.config.WApiUsername, s.config.WApiPassword)
	client := &http.Client{Timeout: time.Duration(60) * time.Second}

	rowID, _ := strconv.Atoi(s.GetRecord(domain, subdomain).RecordID)

	request := &models.Request{Body: models.RequestBody{
		Login:   s.config.WApiUsername,
		Secret:  token,
		Command: "dns-row-delete",
		Data: models.RequestData{
			Domain: domain,
			RowID:  rowID},
	}}

	response, err := client.Do(helpers.BuildRequest(s.baseUrl, request))

	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		return string(bodyBytes)
	}

	return response.Status
}

func (s *integrationService) GetInfo(domainName string) models.WApiResponse {
	token := getApiToken(s.config.WApiUsername, s.config.WApiPassword)
	client := &http.Client{Timeout: time.Duration(60) * time.Second}
	request := &models.Request{Body: models.RequestBody{
		Login:   s.config.WApiUsername,
		Secret:  token,
		Command: "dns-rows-list",
		Data: models.RequestData{
			Domain: domainName},
	}}

	response, err := client.Do(helpers.BuildRequest(s.baseUrl, request))

	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	var result models.WApiResponse

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(bodyBytes, &result)
		if err != nil {
			return models.WApiResponse{}
		}
	}

	return result
}

func (s *integrationService) GetRecord(domain string, subdomain string) models.Record {
	records := s.GetInfo(domain)
	var record models.Record

	for _, row := range records.Body.Data.Records {
		if row.Subdomain == subdomain {
			record = row
		}
	}

	return record
}

func getApiToken(username string, password string) string {
	passwordHash := sha1.New()
	passwordHash.Write([]byte(password))
	location, _ := time.LoadLocation("Europe/Prague")
	hour := formatHour(time.Now().In(location).Hour())
	log.Print(hour)
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
