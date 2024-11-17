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

type CandidateController struct {
	logger           *zap.Logger
	candidateService *service.CandidateService
}

func NewCandidateController(logger *zap.Logger, candidateService *service.CandidateService) *CandidateController {
	return &CandidateController{
		logger:           logger,
		candidateService: candidateService,
	}
}

// CreateCandidate
// @Summary      Create Candidate
// @Description  Create Candidate
// @Tags         Candidate
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        vacancy-id path string true "Vacancy id"
// @Param        payload body   model.AddNewCandidateRequest true "User data"
// @Success      200  {object}  base.ResponseOK "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /candidate/vacancy/{vacancy-id} [post]
func (a *CandidateController) CreateCandidate(c *gin.Context) {
	vacancyId, err := uuid.Parse(c.Param("vacancy-id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, api.GeneralParsingError())
		return
	}

	var payload model.AddNewCandidateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, base.ResponseFailure{
			Status:  http.StatusText(http.StatusBadRequest),
			Blame:   base.BlameUser,
			Message: "failed to parse json",
		})
		return
	}

	if serviceErr := a.candidateService.AddNewCandidate(vacancyId, &payload, c); serviceErr != nil {
		c.JSON(http.StatusBadRequest, serviceErr)
		return
	}

	c.JSON(http.StatusOK, base.ResponseOK{
		Status: http.StatusText(http.StatusOK),
	})
}

// RetrieveCandidate
// @Summary      Retrieve Candidate
// @Description  Retrieve Candidate
// @Tags         Candidate
// @Accept       json
// @Produce      json
// @Param        candidate-id path string true "Candidate id"
// @Success      200  {object}  model.RetrieveCandidateResponse "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /candidate/{candidate-id} [get]
func (a *CandidateController) RetrieveCandidate(c *gin.Context) {
	candidateId, err := uuid.Parse(c.Param("candidate-id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, api.GeneralParsingError())
		return
	}

	candidate, serviceErr := a.candidateService.RetrieveCandidate(candidateId, c)
	if serviceErr != nil {
		c.JSON(http.StatusBadRequest, serviceErr)
		return
	}

	c.JSON(http.StatusOK, model.RetrieveCandidateResponse{
		ResponseOK: base.ResponseOK{
			Status: http.StatusText(http.StatusOK),
		},
		Candidate: *candidate,
	})
}

// GetCandidates
// @Summary      Get Candidates
// @Description  Get Candidates
// @Tags         Candidate
// @Accept       json
// @Produce      json
// @Success      200  {object}   model.GetCandidatesResponse "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /candidate [get]
func (a *CandidateController) GetCandidates(c *gin.Context) {
	candidates, serviceErr := a.candidateService.GetCandidate(dataProcessing.GetOptions(c), c)
	if serviceErr != nil {
		c.JSON(http.StatusBadRequest, serviceErr)
		return
	}

	c.JSON(http.StatusOK, model.GetCandidatesResponse{
		ResponseOK: base.ResponseOK{
			Status: http.StatusText(http.StatusOK),
		},
		Candidates: candidates,
	})
}
