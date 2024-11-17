package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/entity"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/helpers"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/model"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/storage/dao"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

type CandidateService struct {
	logger                 *zap.Logger
	vacancyStorage         *dao.VacancyStorage
	candidateStorage       *dao.CandidateStorage
	cameoMetricsHttpClient *helpers.HttpClient
}

func NewCandidateService(logger *zap.Logger, vacancyStorage *dao.VacancyStorage, candidateStorage *dao.CandidateStorage, cameoMetricsHttpClient *helpers.HttpClient) *CandidateService {
	return &CandidateService{
		logger:                 logger,
		candidateStorage:       candidateStorage,
		vacancyStorage:         vacancyStorage,
		cameoMetricsHttpClient: cameoMetricsHttpClient,
	}
}

func (s *CandidateService) AddNewCandidate(vacancyId uuid.UUID, request *model.AddNewCandidateRequest, ctx context.Context) *base.ServiceError {
	vacancy, err := s.vacancyStorage.Retrieve(vacancyId, ctx)
	if err != nil {
		return base.NewPostgresReadError(err)
	}

	type responseModel struct {
		SystemID string `json:"id"`
	}

	r := &responseModel{}

	type requestModel struct {
		Name     string `json:"name"`
		Position string `json:"position"`
		Status   string `json:"status"`
	}

	requestBody, err := json.Marshal(requestModel{
		Name:     request.Name,
		Position: "test",
		Status:   "test",
	})

	if err != nil {
		return base.NewJsonMarshalError(err)
	}

	oRequest, serviceErr := s.cameoMetricsHttpClient.HttpRequest(
		http.MethodPost,
		"workers",
		bytes.NewBuffer(requestBody),
		ctx,
	)

	if serviceErr != nil {
		return &base.ServiceError{
			Message: "failure create user",
			Blame:   base.BlameServer,
			Code:    serviceErr.Code,
			Err:     serviceErr.Err,
		}
	}

	if err := json.Unmarshal(oRequest, &r); err != nil {
		return base.NewJsonUnmarshalError(err)
	}

	newCandidate := &entity.Candidate{
		Name:      request.Name,
		Email:     request.Email,
		SystemID:  r.SystemID,
		Vacancy:   *vacancy,
		VacancyID: vacancy.ID,
	}

	if err := s.candidateStorage.Create(newCandidate, ctx); err != nil {
		return base.NewPostgresWriteError(err)
	}

	return nil
}

func (s *CandidateService) RetrieveCandidate(candidateID uuid.UUID, ctx context.Context) (*model.CandidateObject, *base.ServiceError) {
	candidate, err := s.candidateStorage.Retrieve(candidateID, ctx)
	if err != nil {
		return nil, base.NewPostgresReadError(err)
	}

	return &model.CandidateObject{
		ID:        candidate.ID,
		CreatedAt: candidate.CreatedAt,
		UpdatedAt: candidate.UpdatedAt,
		Name:      candidate.Name,
		Email:     candidate.Email,
		SystemID:  candidate.SystemID,
	}, nil
}

func (s *CandidateService) GetCandidate(option *dataProcessing.Options, ctx context.Context) ([]model.CandidateObject, *base.ServiceError) {
	candidates, _, err := s.candidateStorage.Get(option, ctx)
	if err != nil {
		return nil, base.NewPostgresReadError(err)
	}

	result := make([]model.CandidateObject, 0, len(candidates))

	for _, candidate := range candidates {
		result = append(result, model.CandidateObject{
			ID:        candidate.ID,
			CreatedAt: candidate.CreatedAt,
			UpdatedAt: candidate.UpdatedAt,
			Name:      candidate.Name,
			Email:     candidate.Email,
			SystemID:  candidate.SystemID,
		})
	}

	return result, nil
}
