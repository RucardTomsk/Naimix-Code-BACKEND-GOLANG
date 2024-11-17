package controller

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/middleware"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/model"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

type CompanyController struct {
	logger         *zap.Logger
	companyService *service.CompanyService
}

func NewCompanyController(logger *zap.Logger, companyService *service.CompanyService) *CompanyController {
	return &CompanyController{
		logger:         logger,
		companyService: companyService,
	}
}

// CreateCompany
// @Summary      Create Company
// @Description  Create Company
// @Tags         Company
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        payload body   model.CreateCompanyRequest true "User data"
// @Success      200  {object}  base.ResponseOK "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /company [post]
func (a *CompanyController) CreateCompany(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)

	var payload model.CreateCompanyRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, base.ResponseFailure{
			Status:  http.StatusText(http.StatusBadRequest),
			Blame:   base.BlameUser,
			Message: "failed to parse json",
		})
		return
	}

	id, serviceErr := a.companyService.CreateNewCompany(userID.(uuid.UUID), &payload, c)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr))
		return
	}

	c.JSON(http.StatusOK, base.ResponseOKWithID{
		ID:     *id,
		Status: http.StatusText(http.StatusOK),
	})
}

// UploadLogo
// @Summary      Upload Logo
// @Description  Upload Logo
// @Tags         Company
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        company-id path string true "Company id"
// @Param        file formData file true "file"
// @Success      200  {object}  base.ResponseOK "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /company/{company-id}/logo [post]
func (a *CompanyController) UploadLogo(c *gin.Context) {
	companyId, err := uuid.Parse(c.Param("company-id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, api.GeneralParsingError())
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, struct {
			Status string
			MSG    string
		}{Status: "error",
			MSG: err.Error()})

		return
	}

	if serviceErr := a.companyService.UploadLogo(companyId, file, c); serviceErr != nil {
		c.JSON(http.StatusBadRequest, base.ResponseFailure{
			Status:  http.StatusText(http.StatusBadRequest),
			Blame:   serviceErr.Blame,
			Message: serviceErr.Message,
		})
		return
	}

	c.JSON(http.StatusOK, base.ResponseOK{
		Status: http.StatusText(http.StatusOK),
	})
}

// RetrieveCompany
// @Summary      Retrieve Company
// @Description  Retrieve Company
// @Tags         Company
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        company-id path string true "Company id"
// @Success      200  {object}  model.RetrieveCompanyResponse "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /company/{company-id} [get]
func (a *CompanyController) RetrieveCompany(c *gin.Context) {
	companyId, err := uuid.Parse(c.Param("company-id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, api.GeneralParsingError())
		return
	}

	company, serviceErr := a.companyService.RetrieveCompany(companyId, c)
	if serviceErr != nil {
		c.JSON(http.StatusBadRequest, base.ResponseFailure{
			Status:  http.StatusText(http.StatusBadRequest),
			Blame:   serviceErr.Blame,
			Message: serviceErr.Message,
		})
		return
	}

	c.JSON(http.StatusOK, model.RetrieveCompanyResponse{
		ResponseOK: base.ResponseOK{
			Status: http.StatusText(http.StatusOK),
		},
		Company: *company,
	})
}

// GetCompany
// @Summary      Get Company
// @Description  Get Company
// @Tags         Company
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200  {object}  model.GetCompanyResponse "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /company [get]
func (a *CompanyController) GetCompany(c *gin.Context) {
	companies, serviceErr := a.companyService.GetCompany(dataProcessing.GetOptions(c), c)
	if serviceErr != nil {
		c.JSON(http.StatusBadRequest, base.ResponseFailure{
			Status:  http.StatusText(http.StatusBadRequest),
			Blame:   serviceErr.Blame,
			Message: serviceErr.Message,
		})
		return
	}

	c.JSON(http.StatusOK, model.GetCompanyResponse{
		ResponseOK: base.ResponseOK{
			Status: http.StatusText(http.StatusOK),
		},
		Companies: companies,
	})
}
