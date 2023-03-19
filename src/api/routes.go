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

	router.POST("/token", func(c *gin.Context) {
		c.Request.URL.Path = "/api/auth/token"
		router.HandleContext(c)
	})

	api := router.Group("/api")
	{
		domain := api.Group("/domain/:domain", middleware.Authorize())
		{
			domain.POST("/record", handlers.CreateRecord)
			domain.PUT("/record", handlers.UpdateRecord)
			domain.GET("/info", handlers.GetDomainInfo)
			domain.GET("/:subdomain/info", handlers.GetSubdomainInfo)
			domain.DELETE("/record", handlers.DeleteRecord)
			domain.POST("commit", handlers.CommitChanges)
		}

		api.POST("/auth/token", handlers.GetToken)
	}
}
