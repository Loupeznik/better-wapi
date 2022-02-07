package services

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/loupeznik/better-wapi/src/models"
	"io"
	"log"
	"net/http"
	"time"
)

type IntegrationService interface {
	CreateRecord()
	UpdateRecord()
	DeleteRecord()
	GetInfo()
}

type service struct {
	config  *models.Config
	baseUrl string
}

func NewIntegrationService(config *models.Config) *service {
	wapiBaseUrl := "https://api.wedos.com/wapi/json"
	return &service{config: config, baseUrl: wapiBaseUrl}
}

func (s *service) UpdateRecord(domainName string, newIp string) int {
	token := getApiToken(s.config.WApiUsername, s.config.WApiPassword)
	client := http.Client{Timeout: time.Duration(60) * time.Second}
	request := models.Request{Body: models.RequestBody{
		Login:   s.config.WApiUsername,
		Secret:  token,
		Command: "dns-row-update",
		Data: models.RequestData{
			Domain: domainName,
			RowID:  20,
			TTL:    20,
			Type:   "A",
			IP:     newIp},
	}}

	jsonBody, err := json.Marshal(request)

	if err != nil {
		panic(err)
	}

	result, err := client.Post(s.baseUrl, "application/x-www-form-urlencoded", bytes.NewBuffer(jsonBody))

	if err != nil {
		panic(err)
	}

	return result.StatusCode
}

func (s *service) GetInfo(domainName string) string {
	token := getApiToken(s.config.WApiUsername, s.config.WApiPassword)
	client := http.Client{Timeout: time.Duration(60) * time.Second}
	request := models.Request{Body: models.RequestBody{
		Login:   s.config.WApiUsername,
		Secret:  token,
		Command: "dns-rows-list",
		Data: models.RequestData{
			Domain: domainName},
	}}

	jsonBody, err := json.Marshal(request)

	if err != nil {
		panic(err)
	}

	result, err := client.Post(fmt.Sprintf("%s?request=%s", s.baseUrl, jsonBody), "application/x-www-form-urlencoded", nil)

	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(result.Body)

	if result.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(result.Body)
		if err != nil {
			log.Fatal(err)
		}

		return string(bodyBytes)
	}

	return ""
}

func getApiToken(username string, password string) string {
	passwordHash := sha1.New()
	passwordHash.Write([]byte(password))
	location, _ := time.LoadLocation("Europe/Prague")
	hour := time.Now().In(location).Hour()
	passwordHashString := fmt.Sprintf("%x", passwordHash.Sum(nil))

	token := fmt.Sprintf("%s%s%d", username, passwordHashString, hour)

	tokenHash := sha1.New()
	tokenHash.Write([]byte(token))

	return fmt.Sprintf("%x", tokenHash.Sum(nil))
}
