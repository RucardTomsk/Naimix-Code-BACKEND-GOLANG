package service

import (
	"context"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/entity"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/enum"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/s3"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/model"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/storage/dao"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"os"
)

type CompanyService struct {
	logger         *zap.Logger
	companyStorage *dao.CompanyStorage
	userService    *UserService
	vacancyService *VacancyService
	fileStorage    *dao.FileStorage
	minioService   s3.ObjectStoreService
}

func NewCompanyService(
	logger *zap.Logger,
	companyStorage *dao.CompanyStorage,
	userService *UserService,
	vacancyService *VacancyService,
	fileStorage *dao.FileStorage,
	minioService s3.ObjectStoreService) *CompanyService {
	return &CompanyService{
		logger:         logger,
		companyStorage: companyStorage,
		userService:    userService,
		vacancyService: vacancyService,
		fileStorage:    fileStorage,
		minioService:   minioService,
	}
}

func (s *CompanyService) CreateNewCompany(ownerID uuid.UUID, request *model.CreateCompanyRequest, ctx context.Context) (*uuid.UUID, *base.ServiceError) {
	newCompany := &entity.Company{
		Name:        request.Name,
		Description: request.Description,
		Owner:       ownerID,
	}

	if err := s.companyStorage.Create(newCompany, ctx); err != nil {
		return nil, base.NewPostgresWriteError(err)
	}

	file, err := os.Open("./static/default_avatar.jpg")
	if err != nil {
		return nil, base.NewPostgresReadError(err)
	}
	defer file.Close()

	if serviceErr := s.UploadLogo(newCompany.ID, file, ctx); serviceErr != nil {
		return nil, serviceErr
	}

	return &newCompany.ID, nil

}

func (s *CompanyService) UploadLogo(companyID uuid.UUID, file io.Reader, ctx context.Context) *base.ServiceError {
	company, err := s.companyStorage.Retrieve(companyID, ctx)
	if err != nil {
		return base.NewPostgresReadError(err)
	}

	fileId, serviceErr := s.minioService.UploadAsWebP(ctx, enum.CompanyLogo, file)
	if serviceErr != nil {
		return serviceErr
	}

	newFile := &entity.File{
		Key:    *fileId,
		Bucket: string(enum.CompanyLogo),
	}

	if err := s.fileStorage.Create(newFile, ctx); err != nil {
		return base.NewPostgresWriteError(err)
	}

	if company.File != nil {
		if serviceErr := s.minioService.DeleteWebPFile(ctx, enum.CompanyLogo, company.File.Key); serviceErr != nil {
			return serviceErr
		}
	}

	company.FileID = fileId
	company.File = newFile

	if err := s.companyStorage.Update(company, ctx); err != nil {
		return base.NewPostgresWriteError(err)
	}

	return nil
}

func (s *CompanyService) getLogoURL(company *entity.Company, ctx context.Context) (*string, *base.ServiceError) {
	var pictureURL *string

	if company.File != nil {
		url, serviceErr := s.minioService.GetWebPFileURL(ctx, enum.CompanyLogo, company.File.Key)
		if serviceErr != nil {
			return nil, serviceErr
		}

		pictureURL = &url
	}

	return pictureURL, nil
}

func (s *CompanyService) RetrieveCompany(companyID uuid.UUID, ctx context.Context) (*model.CompanyObject, *base.ServiceError) {
	company, err := s.companyStorage.Retrieve(companyID, ctx)
	if err != nil {
		return nil, base.NewPostgresReadError(err)
	}

	logoURL, serviceErr := s.getLogoURL(company, ctx)
	if serviceErr != nil {
		return nil, serviceErr
	}

	users := make([]model.UserObject, 0, len(company.Users))

	for _, user := range company.Users {
		um, serviceErr := s.userService.RetrieveUser(user.ID, ctx)
		if serviceErr != nil {
			return nil, serviceErr
		}
		users = append(users, *um)
	}

	vacances := make([]model.VacancyObject, 0, len(company.Vacancies))

	for _, vacancy := range company.Vacancies {
		vm, serviceErr := s.vacancyService.RetrieveVacancy(vacancy.ID, ctx)
		if serviceErr != nil {
			return nil, serviceErr
		}

		vacances = append(vacances, *vm)
	}

	if logoURL != nil {
		return &model.CompanyObject{
			ID:          company.ID,
			CreatedAt:   company.CreatedAt,
			UpdatedAt:   company.UpdatedAt,
			Name:        company.Name,
			Description: company.Description,
			LogoURL:     *logoURL,
			Users:       users,
			Vacancies:   vacances,
		}, nil
	} else {
		return &model.CompanyObject{
			ID:          company.ID,
			CreatedAt:   company.CreatedAt,
			UpdatedAt:   company.UpdatedAt,
			Name:        company.Name,
			Description: company.Description,
			LogoURL:     "",
			Users:       users,
			Vacancies:   vacances,
		}, nil
	}
}

func (s *CompanyService) GetCompany(options *dataProcessing.Options, ctx context.Context) ([]model.CompanyObject, *base.ServiceError) {
	companies, _, err := s.companyStorage.Get(options, ctx)
	if err != nil {
		return nil, base.NewPostgresReadError(err)
	}

	result := make([]model.CompanyObject, 0, len(companies))

	for _, company := range companies {
		logoURL, serviceErr := s.getLogoURL(&company, ctx)
		if serviceErr != nil {
			return nil, serviceErr
		}

		result = append(result, model.CompanyObject{
			ID:          company.ID,
			CreatedAt:   company.CreatedAt,
			UpdatedAt:   company.UpdatedAt,
			Name:        company.Name,
			Description: company.Description,
			LogoURL:     *logoURL,
			Users:       nil,
			Vacancies:   nil,
		})
	}

	return result, nil
}
