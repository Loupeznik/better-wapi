package helpers

import (
	"github.com/joho/godotenv"
	"github.com/loupeznik/better-wapi/src/models"
	"log"
	"os"
)

func SetupIntegrationConfig() *models.Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := models.Config{
		UserLogin:    os.Getenv("user_login"),
		UserSecret:   os.Getenv("user_secret"),
		WApiUsername: os.Getenv("wapi_username"),
		WApiPassword: os.Getenv("wapi_password"),
	}

	return &config
}
