package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apiModels "github.com/loupeznik/better-wapi/src/api/models"
	"github.com/loupeznik/better-wapi/src/helpers"
	_ "github.com/loupeznik/better-wapi/src/models"
	"github.com/loupeznik/better-wapi/src/services"
)

var integrationService *services.IntegrationService

func init() {
	config := helpers.SetupIntegrationConfig()
	integrationService = services.NewIntegrationService(config)
}

// CreateRecord	godoc
// @Summary		Create a new A record
// @Tags		domain
// @Produce		json
// @Accept		json
// @Param 		request	body	apiModels.SaveRowRequest	true	"Request body"
// @Param		domain	path	string	true	"Domain"
// @Success		200
// @Failure		400	{object}	apiModels.ErrorResponse
// @Failure		401	{object}	apiModels.ErrorResponse
// @Failure		404	{object}	apiModels.ErrorResponse
// @Failure		409	{object}	apiModels.ErrorResponse
// @Failure		429	{object}	apiModels.ErrorResponse
// @Failure		500	{object}	apiModels.ErrorResponse
// @Router		/domain/{domain}/record [post]
// @Security	BasicAuth
func CreateRecord(c *gin.Context) {
	domain := c.Param("domain")
	var request apiModels.SaveRowRequest
	err := c.ShouldBindJSON(&request)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	status, err := integrationService.CreateRecord(domain, request.Subdomain, request.IP, request.Autocommit)

	if err != nil {
		c.JSON(status, apiModels.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(status)
}

// UpdateRecord	godoc
// @Summary		Update an existing A record
// @Tags		domain
// @Produce		json
// @Accept		json
// @Param 		request	body	apiModels.SaveRowRequest	true	"Request body"
// @Param		domain	path	string	true	"Domain"
// @Success		200
// @Failure		400	{object}	apiModels.ErrorResponse
// @Failure		401	{object}	apiModels.ErrorResponse
// @Failure		404	{object}	apiModels.ErrorResponse
// @Failure		409	{object}	apiModels.ErrorResponse
// @Failure		429	{object}	apiModels.ErrorResponse
// @Failure		500	{object}	apiModels.ErrorResponse
// @Router		/domain/{domain}/record [put]
// @Security	BasicAuth
func UpdateRecord(c *gin.Context) {
	domain := c.Param("domain")
	var request apiModels.SaveRowRequest
	err := c.ShouldBindJSON(&request)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	status, err := integrationService.UpdateRecord(domain, request.Subdomain, request.IP, request.Autocommit)

	if err != nil {
		c.JSON(status, apiModels.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(status)
}

// DeleteRecord	godoc
// @Summary		Delete an existing A record
// @Tags		domain
// @Produce		json
// @Accept		json
// @Param 		request	body	apiModels.SaveRowRequest	true	"Request body"
// @Param		domain	path	string	true	"Domain"
// @Success		200
// @Failure		400	{object}	apiModels.ErrorResponse
// @Failure		401	{object}	apiModels.ErrorResponse
// @Failure		404	{object}	apiModels.ErrorResponse
// @Failure		409	{object}	apiModels.ErrorResponse
// @Failure		429	{object}	apiModels.ErrorResponse
// @Failure		500	{object}	apiModels.ErrorResponse
// @Router		/domain/{domain}/record [delete]
// @Security	BasicAuth
func DeleteRecord(c *gin.Context) {
	domain := c.Param("domain")
	var request apiModels.SaveRowRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	status, err := integrationService.DeleteRecord(domain, request.Subdomain, request.Autocommit)

	if err != nil {
		c.JSON(status, apiModels.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(status)
}

// GetDomainInfo	godoc
// @Summary		Get all DNS records for a domain
// @Tags		domain
// @Produce		json
// @Param		domain	path	string	true	"Domain"
// @Success		200	{object}	[]models.Record
// @Failure		400	{object}	apiModels.ErrorResponse
// @Failure		401	{object}	apiModels.ErrorResponse
// @Failure		404	{object}	apiModels.ErrorResponse
// @Failure		409	{object}	apiModels.ErrorResponse
// @Failure		429	{object}	apiModels.ErrorResponse
// @Failure		500	{object}	apiModels.ErrorResponse
// @Router		/domain/{domain}/info [get]
// @Security	BasicAuth
func GetDomainInfo(c *gin.Context) {
	domain := c.Param("domain")

	result, status, err := integrationService.GetInfo(domain)

	if err != nil {
		c.JSON(status, apiModels.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(status, result)
}

// GetSubdomainInfo	godoc
// @Summary		Get DNS record for a specific subdomain
// @Tags		domain
// @Produce		json
// @Param		domain	path	string	true	"Domain"
// @Param		subdomain	path	string	true	"Subdomain"
// @Success		200 {object}	models.Record
// @Failure		400	{object}	apiModels.ErrorResponse
// @Failure		401	{object}	apiModels.ErrorResponse
// @Failure		404	{object}	apiModels.ErrorResponse
// @Failure		409	{object}	apiModels.ErrorResponse
// @Failure		429	{object}	apiModels.ErrorResponse
// @Failure		500	{object}	apiModels.ErrorResponse
// @Router		/domain/{domain}/{subdomain}/info [get]
// @Security	BasicAuth
func GetSubdomainInfo(c *gin.Context) {
	domain := c.Param("domain")
	subdomain := c.Param("subdomain")

	result, status, err := integrationService.GetRecord(domain, subdomain)

	if status >= 400 {
		c.JSON(status, apiModels.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(status, result)
}

// CommitChanges	godoc
// @Summary		Commit DNS changes
// @Tags		domain
// @Produce		json
// @Accept		json
// @Param		domain	path	string	true	"Domain"
// @Success		200
// @Failure		400	{object}	apiModels.ErrorResponse
// @Failure		401	{object}	apiModels.ErrorResponse
// @Failure		404	{object}	apiModels.ErrorResponse
// @Failure		409	{object}	apiModels.ErrorResponse
// @Failure		429	{object}	apiModels.ErrorResponse
// @Failure		500	{object}	apiModels.ErrorResponse
// @Security	BasicAuth
func CommitChanges(c *gin.Context) {
	domain := c.Param("domain")

	status, err := integrationService.CommitChanges(domain)

	if err != nil {
		c.JSON(status, apiModels.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(status)
}
