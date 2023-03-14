package helpers

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/loupeznik/better-wapi/src/models"
)

func SetupIntegrationConfig() *models.Config {
	_ = godotenv.Load()

	useLogFile, err := strconv.ParseBool(os.Getenv("BW_USE_LOGFILE"))

	if err != nil {
		useLogFile = false
	}

	config := models.Config{
		UserLogin:    os.Getenv("BW_USER_LOGIN"),
		UserSecret:   os.Getenv("BW_USER_SECRET"),
		WApiUsername: os.Getenv("BW_WAPI_USERNAME"),
		WApiPassword: os.Getenv("BW_WAPI_PASSWORD"),
		BaseUrl:      os.Getenv("BW_BASE_URL"),
		UseLogFile:   useLogFile,
		JsonWebKey:   os.Getenv("BW_JSON_WEB_KEY"),
	}

	return &config
}
