package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/loupeznik/better-wapi/src/api/handlers"
	"github.com/loupeznik/better-wapi/src/api/middleware"
	"github.com/loupeznik/better-wapi/src/models"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ResolveMiddleware(authMode models.AuthMode) gin.HandlerFunc {
	switch authMode {
	case models.AuthModeBasic:
		return middleware.AuthorizeBasic()
	case models.AuthModeBearer:
		return middleware.AuthorizeInternalJwt()
	case models.AuthModeOAuth:
		return middleware.AuthorizeOAuthJwt()
	default:
		return func(c *gin.Context) {
			c.AbortWithError(500, errors.New("unknown auth mode"))
		}
	}
}

func SetupRoutes(router *gin.Engine, authMode models.AuthMode) {
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
		v1 := api.Group("/v1")
		{
			domain := v1.Group("/domain/:domain", ResolveMiddleware(authMode))
			{
				domain.POST("/record", handlers.CreateRecord)
				domain.PUT("/record", handlers.UpdateRecord)
				domain.GET("/info", handlers.GetDomainInfo)
				domain.GET("/:subdomain/info", handlers.GetSubdomainInfo)
				domain.DELETE("/record", handlers.DeleteRecord)
				domain.POST("commit", handlers.CommitChanges)
			}
		}
		v2 := api.Group("/v2")
		{
			domain := v2.Group("/domain/:domain", ResolveMiddleware(authMode))
			{
				domain.PUT("/record/:id", handlers.UpdateRecordById)
				domain.DELETE("/record/:id", handlers.DeleteRecordById)
			}
		}

		api.POST("/auth/token", handlers.GetToken)
	}
}
