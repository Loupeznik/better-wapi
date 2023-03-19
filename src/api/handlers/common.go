package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apiModels "github.com/loupeznik/better-wapi/src/api/models"
	"github.com/loupeznik/better-wapi/src/helpers"
	"github.com/loupeznik/better-wapi/src/models"
	"github.com/loupeznik/better-wapi/src/services"
)

var authService *services.AuthService

func init() {
	config := helpers.SetupIntegrationConfig()
	authService = services.NewAuthService(config)
}

// GetToken	godoc
// @Summary		Get token
// @Tags		auth
// @Produce		json
// @Accept		json
// @Param 		request	body	models.Login	true	"Request body"
// @Success		200	{object}	models.TokenResponse
// @Failure		400	{object}	apiModels.ErrorResponse
// @Failure		401	{object}	apiModels.ErrorResponse
// @Router		/auth/token [post]
func GetToken(c *gin.Context) {
	var request models.Login
	err := c.ShouldBindJSON(&request)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	token, err := authService.IssueToken(request)

	if err != nil {
		c.JSON(http.StatusUnauthorized, apiModels.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.TokenResponse{
		Token: token,
	})
}
