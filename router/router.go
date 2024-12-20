package router

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/config"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/controller"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/middleware"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/auth"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/common"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/entity"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type Router struct {
	config config.Config
}

func NewRouter(config config.Config) *Router {
	return &Router{
		config: config,
	}
}

func (h *Router) InitRoutes(
	logger *zap.Logger,
	controllerContainer *controller.Container,
	dataProcessing *dataProcessing.DataProcessing,
	JWTManager *auth.JWTManager,
	adminID uuid.UUID,
) *gin.Engine {
	gin.SetMode(h.config.Server.GinMode)

	router := gin.Default()

	router.Use(middleware.SetRecoveryHandler(*logger))
	router.Use(cors.New(common.DefaultCorsConfig()))

	router.GET("api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	baseRouter := router.Group("/api")

	candidate := baseRouter.Group("/candidate")
	{
		candidate.POST("vacancy/:vacancy-id", controllerContainer.CandidateController.CreateCandidate)
		candidate.GET(":candidate-id", controllerContainer.CandidateController.RetrieveCandidate)
		candidate.GET("", dataProcessing.ApplyMiddleware(*logger, entity.Candidate{}.FilteringRules(), nil), controllerContainer.CandidateController.GetCandidates)
	}

	vacancy := baseRouter.Group("/vacancy")
	{
		vacancy.POST("company/:company-id", controllerContainer.VacancyController.CreateVacancy)
		vacancy.GET(":vacancy-id", controllerContainer.VacancyController.RetrieveVacancy)
		vacancy.GET("", dataProcessing.ApplyMiddleware(*logger, entity.Vacancy{}.FilteringRules(), nil), controllerContainer.VacancyController.GetVacancy)
	}

	company := baseRouter.Group("/company")
	{
		company.POST("", middleware.SetAuthorizationCheck(JWTManager, *logger), controllerContainer.CompanyController.CreateCompany)
		company.POST(":company-id/logo", middleware.SetAuthorizationCheck(JWTManager, *logger), controllerContainer.CompanyController.UploadLogo)
		company.GET(":company-id", controllerContainer.CompanyController.RetrieveCompany)
		company.GET("", dataProcessing.ApplyMiddleware(*logger, entity.Company{}.FilteringRules(), nil), controllerContainer.CompanyController.GetCompany)
	}

	user := baseRouter.Group("user")
	{
		user.POST("register",
			middleware.SetAuthorizationAdminCheck(JWTManager, adminID, *logger),
			controllerContainer.AuthController.Register)
		user.POST("login", controllerContainer.AuthController.Login)
		user.POST("refresh", controllerContainer.AuthController.RecreateJWT)
		user.POST(
			"logout",
			middleware.SetAuthorizationCheck(JWTManager, *logger),
			controllerContainer.AuthController.Logout)
		user.DELETE(":user-id/delete", middleware.SetAuthorizationAdminCheck(JWTManager, adminID, *logger), controllerContainer.UserController.DeleteUser)
		user.POST("field/update", middleware.SetAuthorizationCheck(JWTManager, *logger), controllerContainer.UserController.Update)
		user.GET(
			"get",
			middleware.SetAuthorizationCheck(JWTManager, *logger),
			dataProcessing.ApplyMiddleware(*logger, entity.User{}.FilteringRules(), nil),
			controllerContainer.UserController.Get)
		user.GET("retrieve", middleware.SetAuthorizationCheck(JWTManager, *logger), controllerContainer.UserController.RetrieveUser)
	}

	return router
}
