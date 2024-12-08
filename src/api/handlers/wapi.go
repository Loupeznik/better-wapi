package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	apiModels "github.com/loupeznik/better-wapi/src/api/models"
	"github.com/loupeznik/better-wapi/src/helpers"

	"github.com/creasty/defaults"
	"github.com/loupeznik/better-wapi/src/services"
)

var integrationService *services.IntegrationService

func init() {
	config := helpers.SetupIntegrationConfig()
	integrationService = services.NewIntegrationService(config)
}

// CreateRecord	godoc
// @Summary		Create a new record
// @Tags		domain
// @Produce		json
// @Accept		json
// @Param 		request	body	apiModels.SaveRowRequest	true	"Request body"
// @Param		domain	path	string	true	"Domain"
// @Success		201
// @Failure		400	{object}	apiModels.ErrorResponse
// @Failure		401	{object}	apiModels.ErrorResponse
// @Failure		404	{object}	apiModels.ErrorResponse
// @Failure		409	{object}	apiModels.ErrorResponse
// @Failure		429	{object}	apiModels.ErrorResponse
// @Failure		500	{object}	apiModels.ErrorResponse
// @Router		/v1/domain/{domain}/record [post]
func CreateRecord(c *gin.Context) {
	domain := c.Param("domain")
	var request apiModels.SaveRowRequest
	err := c.BindJSON(&request)

	if err != nil {
		returnValidationError(c, http.StatusBadRequest, nil)
	}

	if err := defaults.Set(&request); err != nil {
		returnValidationError(c, http.StatusInternalServerError, err)
	}

	status, err := integrationService.CreateRecord(domain, request)

	if err != nil {
		c.JSON(status, apiModels.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if status == 200 {
		status = http.StatusCreated
	}

	c.Status(status)
}

// UpdateRecord	godoc
// @Summary		Update an existing record
// @Tags		domain
// @Produce		json
// @Accept		json
// @Param 		request	body	apiModels.SaveRowRequest	true	"Request body"
// @Param		domain	path	string	true	"Domain"
// @Success		204
// @Failure		400	{object}	apiModels.ErrorResponse
// @Failure		401	{object}	apiModels.ErrorResponse
// @Failure		404	{object}	apiModels.ErrorResponse
// @Failure		409	{object}	apiModels.ErrorResponse
// @Failure		429	{object}	apiModels.ErrorResponse
// @Failure		500	{object}	apiModels.ErrorResponse
// @Router		/v1/domain/{domain}/record [put]
func UpdateRecord(c *gin.Context) {
	domain := c.Param("domain")
	var request apiModels.SaveRowRequest
	err := c.ShouldBindJSON(&request)

	if err != nil {
		returnValidationError(c, http.StatusBadRequest, nil)
	}

	if err := defaults.Set(&request); err != nil {
		returnValidationError(c, http.StatusInternalServerError, err)
	}

	status, err := integrationService.UpdateRecord(domain, request)

	if err != nil {
		c.JSON(status, apiModels.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if status == 200 {
		status = http.StatusNoContent
	}

	c.Status(status)
}

// UpdateRecord	godoc
// @Summary		Update an existing record
// @Tags		domain
// @Produce		json
// @Accept		json
// @Param 		request	body	apiModels.SaveRowRequestV2	true	"Request body"
// @Param		domain	path	string	true	"Domain"
// @Param		id		path	int		true	"Record ID"
// @Success		204
// @Failure		400	{object}	apiModels.ErrorResponse
// @Failure		401	{object}	apiModels.ErrorResponse
// @Failure		404	{object}	apiModels.ErrorResponse
// @Failure		409	{object}	apiModels.ErrorResponse
// @Failure		429	{object}	apiModels.ErrorResponse
// @Failure		500	{object}	apiModels.ErrorResponse
// @Router		/v2/domain/{domain}/record/{id} [put]
func UpdateRecordById(c *gin.Context) {
	domain := c.Param("domain")
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		returnValidationError(c, http.StatusBadRequest, errors.New("id not provided"))
	}

	var request apiModels.SaveRowRequestV2
	err = c.ShouldBindJSON(&request)

	if err != nil {
		returnValidationError(c, http.StatusBadRequest, nil)
	}

	if err := defaults.Set(&request); err != nil {
		returnValidationError(c, http.StatusInternalServerError, err)
	}

	status, err := integrationService.UpdateRecordV2(domain, id, request)

	if err != nil {
		c.JSON(status, apiModels.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(status)
}

// DeleteRecord	godoc
// @Summary		Delete an existing record
// @Tags		domain
// @Produce		json
// @Accept		json
// @Param 		request	body	apiModels.DeleteRowRequest	true	"Request body"
// @Param		domain	path	string	true	"Domain"
// @Success		200
// @Failure		400	{object}	apiModels.ErrorResponse
// @Failure		401	{object}	apiModels.ErrorResponse
// @Failure		404	{object}	apiModels.ErrorResponse
// @Failure		409	{object}	apiModels.ErrorResponse
// @Failure		429	{object}	apiModels.ErrorResponse
// @Failure		500	{object}	apiModels.ErrorResponse
// @Router		/v1/domain/{domain}/record [delete]
func DeleteRecord(c *gin.Context) {
	domain := c.Param("domain")
	var request apiModels.DeleteRowRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		returnValidationError(c, http.StatusBadRequest, nil)
	}

	status, err := integrationService.DeleteRecord(domain, request.Subdomain, *request.Autocommit)

	if err != nil {
		c.JSON(status, apiModels.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(status)
}

// DeleteRecord	godoc
// @Summary		Delete an existing record
// @Tags		domain
// @Produce		json
// @Accept		json
// @Param 		request	body	apiModels.DeleteRowRequestV2	true	"Request body"
// @Param		domain	path	string	true	"Domain"
// @Param		id		path	int		true	"Record ID"
// @Success		204
// @Failure		400	{object}	apiModels.ErrorResponse
// @Failure		401	{object}	apiModels.ErrorResponse
// @Failure		404	{object}	apiModels.ErrorResponse
// @Failure		500	{object}	apiModels.ErrorResponse
// @Router		/v2/domain/{domain}/record/{id} [delete]
func DeleteRecordById(c *gin.Context) {
	domain := c.Param("domain")
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		returnValidationError(c, http.StatusBadRequest, errors.New("id not provided"))
	}

	var request apiModels.DeleteRowRequestV2

	err = c.ShouldBindJSON(&request)

	if err != nil {
		returnValidationError(c, http.StatusBadRequest, err)
		return
	}

	status, err := integrationService.DeleteRecordV2(domain, id, *request.Autocommit)

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
// @Router		/v1/domain/{domain}/info [get]
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
// @Router		/v1/domain/{domain}/{subdomain}/info [get]
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
// @Router		/v1/domain/{domain}/commit [post]
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

func returnValidationError(c *gin.Context, status int, err error) {
	c.Header("Content-Type", "application/problem+json")

	if err != nil {
		c.AbortWithStatusJSON(status, apiModels.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.AbortWithStatus(status)
}
