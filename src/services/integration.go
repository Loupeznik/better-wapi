package services

import (
	"crypto/sha1"
	"fmt"
	"github.com/loupeznik/better-wapi/src/models"
	"time"
)

type IntegrationService interface {
	CreateRecord()
	UpdateRecord()
	DeleteRecord()
	GetInfo()
}

type service struct {
	config *models.Config
}

func NewIntegrationService(config *models.Config) *service {
	return &service{config: config}
}

func (s *service) UpdateRecord() {
	//token := getApiToken(s.config.WApiUsername, s.config.WApiPassword)
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
