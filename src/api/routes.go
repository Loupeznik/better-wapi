package api

import (
	"github.com/gin-gonic/gin"
	"github.com/loupeznik/better-wapi/src/api/handlers"
	"github.com/loupeznik/better-wapi/src/api/middleware"
	"github.com/loupeznik/better-wapi/src/models"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(config *models.Config, router *gin.Engine) {
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/docs/index.html")
	})

	api := router.Group("/api", middleware.Authorize(config))
	{
		api.POST("/domain/:domain/record", handlers.CreateRecord)
		api.PUT("/domain/:domain/record", handlers.UpdateRecord)
		api.GET("/domain/:domain/info", handlers.GetDomainInfo)
		api.GET("/domain/:domain/:subdomain/info", handlers.GetSubdomainInfo)
		api.DELETE("/domain/:domain/record", handlers.DeleteRecord)
		api.POST("/domain/:domain/commit", handlers.CommitChanges)
	}
}
