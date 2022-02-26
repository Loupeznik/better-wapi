package middleware

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/loupeznik/better-wapi/src/models"
	"github.com/loupeznik/better-wapi/src/services"
	"strings"
)

func Authorize(config *models.Config) gin.HandlerFunc {
	authService := services.NewAuthService(config)

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

func respondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}

	c.JSON(code, resp)
	c.Abort()
}
