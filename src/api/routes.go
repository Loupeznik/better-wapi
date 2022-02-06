package api

import (
	"github.com/gin-gonic/gin"
	"github.com/loupeznik/better-wapi/src/models"
	"net/http"
)

func SetupRoutes(config *models.Config, router *gin.Engine) {
	//integrationService := services.NewIntegrationService(config)

	api := router.Group("/api")
	{
		api.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "test",
			})
		})
	}
}
