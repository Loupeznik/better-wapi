package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/loupeznik/better-wapi/docs"

	"github.com/gin-gonic/gin"
	"github.com/loupeznik/better-wapi/src/api"
	"github.com/loupeznik/better-wapi/src/helpers"
)

// @title	Better WAPI
// @version	1.0
// @description	A wrapper around the Wedos API (WAPI)

// @license.name	GNU General Public License v3.0
// @license.url	https://github.com/Loupeznik/better-wapi/blob/master/LICENSE

// @BasePath	/api
// @host	http://localhost:8000

// @securityDefinitions.basic	BasicAuth
func main() {
	router := gin.Default()
	config := helpers.SetupIntegrationConfig()

	if config.BaseUrl != "" {
		docs.SwaggerInfo.Host = config.BaseUrl
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AddAllowMethods("OPTIONS")

	router.Use(cors.New(corsConfig))

	api.SetupRoutes(config, router)

	err := router.Run("0.0.0.0:8000")
	if err != nil {
		log.Fatal("Could not start router")
	}
}
