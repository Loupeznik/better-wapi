package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loupeznik/better-wapi/src/api/middleware"
	requests "github.com/loupeznik/better-wapi/src/api/models"
	"github.com/loupeznik/better-wapi/src/models"
	"github.com/loupeznik/better-wapi/src/services"
)

func SetupRoutes(config *models.Config, router *gin.Engine) {
	integrationService := services.NewIntegrationService(config)

	api := router.Group("/api", middleware.Authorize(config))
	{
		api.GET("/test", func(c *gin.Context) {

			c.JSON(http.StatusOK, gin.H{
				"message": "test",
			})
		})

		api.PUT("/domain/:domain/record", func(c *gin.Context) {
			domain := c.Param("domain")
			var request requests.SaveRowRequest
			err := c.ShouldBindJSON(&request)

			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
			}

			updateResult := integrationService.UpdateRecord(domain, request.Subdomain, request.IP)

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

		api.GET("/domain/:domain/:subdomain/info", func(c *gin.Context) {
			domain := c.Param("domain")
			subdomain := c.Param("subdomain")

			getRecordResult := integrationService.GetRecord(domain, subdomain)

			c.JSON(http.StatusOK, gin.H{
				"data": getRecordResult,
			})
		})

		api.POST("/domain/:domain/record", func(c *gin.Context) {
			domain := c.Param("domain")
			var request requests.SaveRowRequest
			err := c.ShouldBindJSON(&request)

			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
			}

			result := integrationService.CreateRecord(domain, request.Subdomain, request.IP)

			c.JSON(http.StatusOK, gin.H{
				"data": result,
			})
		})

		api.DELETE("/domain/:domain/record", func(c *gin.Context) {
			domain := c.Param("domain")
			var request requests.SaveRowRequest

			err := c.ShouldBindJSON(&request)

			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
			}

			result := integrationService.DeleteRecord(domain, request.Subdomain)

			c.JSON(http.StatusOK, gin.H{
				"data": result,
			})
		})

		api.POST("/domain/:domain/commit", func(c *gin.Context) {
			domain := c.Param("domain")

			result := integrationService.CommitChanges(domain)

			c.JSON(http.StatusOK, gin.H{
				"data": result,
			})
		})
	}
}
