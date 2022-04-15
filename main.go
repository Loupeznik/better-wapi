package main

import (
	"github.com/gin-gonic/gin"
	"github.com/loupeznik/better-wapi/src/api"
	"github.com/loupeznik/better-wapi/src/helpers"
	"log"
)

func main() {
	router := gin.Default()
	config := helpers.SetupIntegrationConfig()

	api.SetupRoutes(config, router)

	err := router.Run("0.0.0.0:8000")
	if err != nil {
		log.Fatal("Could not start router")
	}
}
