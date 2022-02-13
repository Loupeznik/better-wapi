package api

import (
	"github.com/gin-gonic/gin"
	requests "github.com/loupeznik/better-wapi/src/api/models"
	"github.com/loupeznik/better-wapi/src/models"
	"github.com/loupeznik/better-wapi/src/services"
	"net/http"
)

func SetupRoutes(config *models.Config, router *gin.Engine) {
	integrationService := services.NewIntegrationService(config)

	api := router.Group("/api")
	{
		api.GET("/test", func(c *gin.Context) {

			c.JSON(http.StatusOK, gin.H{
				"message": "test",
			})
		})

		api.PUT("/domain", func(c *gin.Context) {
			var request requests.UpdateRequest
			err := c.ShouldBindJSON(&request)

			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
			}

			updateResult := integrationService.UpdateRecord(request.Domain, request.IP)

			c.JSON(http.StatusOK, gin.H{
				"data": updateResult,
			})
		})

		api.GET("/domain/:domain/info", func(c *gin.Context) {
			domain := c.Param("domain")

			getInfoResult := integrationService.GetInfo(domain)

			c.JSON(http.StatusOK, gin.H{
				"data": getInfoResult,
			})
		})
	}
}
