package helpers

import (
	"github.com/joho/godotenv"
	"github.com/loupeznik/better-wapi/src/models"
	"os"
)

func SetupIntegrationConfig() *models.Config {
	_ = godotenv.Load()

	config := models.Config{
		UserLogin:    os.Getenv("BW_USER_LOGIN"),
		UserSecret:   os.Getenv("BW_USER_SECRET"),
		WApiUsername: os.Getenv("BW_WAPI_USERNAME"),
		WApiPassword: os.Getenv("BW_WAPI_PASSWORD"),
	}

	return &config
}
