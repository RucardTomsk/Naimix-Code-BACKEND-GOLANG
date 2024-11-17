package service

import (
	"context"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/entity"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/model"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/storage/dao"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type VacancyService struct {
	logger           *zap.Logger
	companyStorage   *dao.CompanyStorage
	vacancyStorage   *dao.VacancyStorage
	candidateService *CandidateService
}

func NewVacancyService(logger *zap.Logger, vacancyStorage *dao.VacancyStorage, companyStorage *dao.CompanyStorage, candidateService *CandidateService) *VacancyService {
	return &VacancyService{
		logger:           logger,
		vacancyStorage:   vacancyStorage,
		companyStorage:   companyStorage,
		candidateService: candidateService,
	}
}

func (s *VacancyService) CreateVacancy(companyID uuid.UUID, request *model.CreateNewVacancyRequest, ctx context.Context) (*uuid.UUID, *base.ServiceError) {
	company, err := s.companyStorage.Retrieve(companyID, ctx)
	if err != nil {
		return nil, base.NewPostgresReadError(err)
	}

	newVacancy := &entity.Vacancy{
		Name:        request.Name,
		Salary:      request.Salary,
		City:        request.City,
		Description: request.Description,
		Company:     *company,
		CompanyID:   company.ID,
	}

	if err := s.vacancyStorage.Create(newVacancy, ctx); err != nil {
		return nil, base.NewPostgresWriteError(err)
	}

	return &newVacancy.ID, nil
}

func (s *VacancyService) RetrieveVacancy(vacancyId uuid.UUID, ctx context.Context) (*model.VacancyObject, *base.ServiceError) {
	vacancy, err := s.vacancyStorage.Retrieve(vacancyId, ctx)
	if err != nil {
		return nil, base.NewPostgresReadError(err)
	}

	co := make([]model.CandidateObject, 0, len(vacancy.Candidates))

	for _, v := range vacancy.Candidates {
		q, serviceErr := s.candidateService.RetrieveCandidate(v.ID, ctx)
		if serviceErr != nil {
			return nil, serviceErr
		}
		co = append(co, *q)
	}

	return &model.VacancyObject{
		ID:          vacancy.ID,
		CreatedAt:   vacancy.CreatedAt,
		UpdatedAt:   vacancy.UpdatedAt,
		Name:        vacancy.Name,
		Salary:      vacancy.Salary,
		City:        vacancy.City,
		Description: vacancy.Description,
		Candidates:  co,
	}, nil
}

func (s *VacancyService) GetVacancy(options *dataProcessing.Options, ctx context.Context) ([]model.VacancyObject, *base.ServiceError) {
	vacancy, _, err := s.vacancyStorage.Get(options, ctx)
	if err != nil {
		return nil, base.NewPostgresReadError(err)
	}

	result := make([]model.VacancyObject, 0, len(vacancy))

	for _, v := range vacancy {
		result = append(result, model.VacancyObject{
			ID:          v.ID,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
			Name:        v.Name,
			Salary:      v.Salary,
			City:        v.City,
			Description: v.Description,
			Candidates:  nil,
		})
	}

	return result, nil
}
