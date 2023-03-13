package middleware

import (
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"
	requests "github.com/loupeznik/better-wapi/src/api/models"
	"github.com/loupeznik/better-wapi/src/helpers"
	"github.com/loupeznik/better-wapi/src/models"
	"github.com/loupeznik/better-wapi/src/services"
)

var authService *services.AuthService

func init() {
	config := helpers.SetupIntegrationConfig()
	authService = services.NewAuthService(config)
}

// Deprecated
func AuthorizeBasic(config *models.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 || auth[0] != "Basic" {
			respondWithError(401, "Unauthorized", c)
			return
		}
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		isAuthenticated := authService.ValidateCredentials(models.Login{
			Login:  pair[0],
			Secret: pair[1],
		})

		if len(pair) != 2 || !isAuthenticated {
			respondWithError(401, "Unauthorized", c)
			return
		}

		c.Next()
	}
}

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 || auth[0] != "Bearer" {
			respondWithError(401, "Unauthorized", c)
			return
		}

		_, err := authService.ValidateToken(auth[1])
		if err != nil {
			respondWithError(401, "Unauthorized", c)
			return
		}

		c.Next()
	}
}

func respondWithError(code int, message string, c *gin.Context) {
	resp := requests.ErrorResponse{
		Error: message,
	}

	c.JSON(code, resp)
	c.Abort()
}
