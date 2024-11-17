package controller

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/service"
	"go.uber.org/zap"
)

type Container struct {
	AuthController      *AuthController
	UserController      *UserController
	CompanyController   *CompanyController
	VacancyController   *VacancyController
	CandidateController *CandidateController
}

func NewControllerContainer(
	logger *zap.Logger,
	authService *service.AuthService,
	userService *service.UserService,
	companyService *service.CompanyService,
	vacancyService *service.VacancyService,
	candidateService *service.CandidateService,
) *Container {
	return &Container{
		AuthController:      NewAuthController(logger, authService),
		UserController:      NewUserController(logger, userService),
		CompanyController:   NewCompanyController(logger, companyService),
		VacancyController:   NewVacancyController(logger, vacancyService),
		CandidateController: NewCandidateController(logger, candidateService),
	}
}
