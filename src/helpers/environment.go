package helpers

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/loupeznik/better-wapi/src/models"
)

func SetupIntegrationConfig() *models.Config {
	_ = godotenv.Load()

	config := models.Config{
		UserLogin:    os.Getenv("BW_USER_LOGIN"),
		UserSecret:   os.Getenv("BW_USER_SECRET"),
		WApiUsername: os.Getenv("BW_WAPI_USERNAME"),
		WApiPassword: os.Getenv("BW_WAPI_PASSWORD"),
		BaseUrl:      os.Getenv("BW_BASE_URL"),
	}

	return &config
}
