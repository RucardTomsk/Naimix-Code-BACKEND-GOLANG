package controller

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/model"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

type VacancyController struct {
	logger         *zap.Logger
	vacancyService *service.VacancyService
}

func NewVacancyController(
	logger *zap.Logger,
	vacancyService *service.VacancyService) *VacancyController {
	return &VacancyController{
		logger:         logger,
		vacancyService: vacancyService,
	}
}

// CreateVacancy
// @Summary      Create Vacancy
// @Description  Create Vacancy
// @Tags         Vacancy
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        company-id path string true "Company id"
// @Param        payload body   model.CreateNewVacancyRequest true "User data"
// @Success      200  {object}  base.ResponseOK "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /vacancy/company/{company-id} [post]
func (a *VacancyController) CreateVacancy(c *gin.Context) {
	companyId, err := uuid.Parse(c.Param("company-id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, api.GeneralParsingError())
		return
	}

	var payload model.CreateNewVacancyRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, base.ResponseFailure{
			Status:  http.StatusText(http.StatusBadRequest),
			Blame:   base.BlameUser,
			Message: "failed to parse json",
		})
		return
	}

	id, serviceErr := a.vacancyService.CreateVacancy(companyId, &payload, c)
	if serviceErr != nil {
		c.JSON(http.StatusBadRequest, serviceErr)
		return
	}

	c.JSON(http.StatusOK, base.ResponseOKWithID{
		ID:     *id,
		Status: http.StatusText(http.StatusOK),
	})
}

// RetrieveVacancy
// @Summary      Retrieve Vacancy
// @Description  Retrieve Vacancy
// @Tags         Vacancy
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        vacancy-id path string true "Vacancy id"
// @Success      200  {object}  model.RetrieveVacancyResponse "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /vacancy/{vacancy-id} [get]
func (a *VacancyController) RetrieveVacancy(c *gin.Context) {
	vacancyId, err := uuid.Parse(c.Param("vacancy-id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, api.GeneralParsingError())
		return
	}

	vacancy, serviceErr := a.vacancyService.RetrieveVacancy(vacancyId, c)
	if serviceErr != nil {
		c.JSON(http.StatusBadRequest, serviceErr)
		return
	}

	c.JSON(http.StatusOK, model.RetrieveVacancyResponse{
		ResponseOK: base.ResponseOK{
			Status: http.StatusText(http.StatusOK),
		},
		Vacancy: *vacancy,
	})
}

// GetVacancy
// @Summary      Get Vacancy
// @Description  Get Vacancy
// @Tags         Vacancy
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200  {object}  model.GetVacancyResponse "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /vacancy [get]
func (a *VacancyController) GetVacancy(c *gin.Context) {
	vacancies, serviceErr := a.vacancyService.GetVacancy(dataProcessing.GetOptions(c), c)
	if serviceErr != nil {
		c.JSON(http.StatusBadRequest, base.ResponseFailure{
			Status:  http.StatusText(http.StatusBadRequest),
			Blame:   serviceErr.Blame,
			Message: serviceErr.Message,
		})
		return
	}

	c.JSON(http.StatusOK, model.GetVacancyResponse{
		ResponseOK: base.ResponseOK{
			Status: http.StatusText(http.StatusOK),
		},
		Vacancies: vacancies,
	})
}
