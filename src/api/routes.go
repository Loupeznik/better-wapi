package api

import (
	"github.com/gin-gonic/gin"
	"github.com/loupeznik/better-wapi/src/api/handlers"
	"github.com/loupeznik/better-wapi/src/api/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/docs/index.html")
	})

	router.POST("/token", handlers.GetToken)

	api := router.Group("/api", middleware.Authorize())
	{
		api.POST("/domain/:domain/record", handlers.CreateRecord)
		api.PUT("/domain/:domain/record", handlers.UpdateRecord)
		api.GET("/domain/:domain/info", handlers.GetDomainInfo)
		api.GET("/domain/:domain/:subdomain/info", handlers.GetSubdomainInfo)
		api.DELETE("/domain/:domain/record", handlers.DeleteRecord)
		api.POST("/domain/:domain/commit", handlers.CommitChanges)
	}
}
