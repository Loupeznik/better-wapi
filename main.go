package main

import (
	"github.com/gin-gonic/gin"
	"github.com/loupeznik/better-wapi/src/api"
	"github.com/loupeznik/better-wapi/src/helpers"
)

func main() {
	router := gin.Default()
	config := helpers.SetupIntegrationConfig()

	api.SetupRoutes(config, router)

	err := router.Run("localhost:8000")
	if err != nil {
		return
	}
}
