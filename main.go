package main

import (
	"context"
	"fmt"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/config"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/controller"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/docs"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/auth"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/common"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/helpers"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/s3"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/telemetry/log"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/router"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/server"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/service"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/storage/dao"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/storage/migration"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	logger := log.NewLogger()

	appCli := common.InitAppCli()
	if err := appCli.Run(os.Args); err != nil {
		logger.Fatal(err.Error())
	}

	// read config
	var cfg config.Config
	if err := viper.MergeInConfig(); err != nil {
		logger.Fatal(fmt.Sprintf("error reading config file: %v", err))
	}

	err := viper.Unmarshal(&cfg)
	if err != nil {
		logger.Fatal(fmt.Sprintf("unable to decode into struct: %v", err))
	}

	// configure swagger
	swaggerConfig := common.NewSwaggerConfig("User API", "TBD", "unreleased")

	docs.SwaggerInfo.Title = swaggerConfig.Title
	docs.SwaggerInfo.Description = swaggerConfig.Description
	docs.SwaggerInfo.Version = swaggerConfig.Version
	docs.SwaggerInfo.BasePath = swaggerConfig.BasePath
	docs.SwaggerInfo.Schemes = swaggerConfig.Schemes

	// init connections
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Port)

	gormLogger := gormlogger.New(
		log.NewZapWriter(logger), // io writer
		gormlogger.Config{
			SlowThreshold:             time.Second,       // slow SQL threshold
			LogLevel:                  gormlogger.Silent, // log level
			IgnoreRecordNotFoundError: true,              // ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,              // don't include params in the SQL log
			Colorful:                  false,             // disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		logger.Fatal(fmt.Sprintf("can't connect to database: %v", err))
	}

	logger.Info(fmt.Sprintf("successfully connected to database %s on %s:%d as %s",
		cfg.DB.Name, cfg.DB.Host, cfg.DB.Port, cfg.DB.User))

	hasher := auth.NewHasher(cfg.Auth.Salt)

	adminID, err := uuid.Parse(cfg.AdminMigration.AdminID)
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed parse uuid admin: %v", err))
	}

	adminPassword, err := hasher.Hash(cfg.AdminMigration.AdminPassword)
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed hash admin password: %v", err))
	}

	if err := migration.Migrate(
		db,
		adminID,
		cfg.AdminMigration.AdminUserName,
		cfg.AdminMigration.AdminEmail,
		adminPassword); err != nil {
		logger.Fatal(fmt.Sprintf("failed to migrate database: %v", err))
	}

	logger.Info("database migrated successfully")

	jwtManager, err := auth.NewJWTManager(cfg.Auth.SigningKey, cfg.Auth.TimeToLive)
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to create JWTManager: %v", err))
	}

	//init minio
	// init minio connection
	var minioService s3.ObjectStoreService
	if cfg.Minio.UseMocks {
		logger.Warn("using minio mock instead of real service")
	} else {
		minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.Minio.AccessKey, cfg.Minio.SecretKey, cfg.Minio.Token),
			Secure: cfg.Minio.UseSSL,
		})
		if err != nil {
			logger.Fatal(fmt.Sprintf("failed to init minio client: %v", err))
		}

		logger.Info(fmt.Sprintf("connected to minio on %s", minioClient.EndpointURL().String()))
		minioService = s3.NewMinioService(minioClient)
	}

	// init mail service
	//mailService := mail.NewMailService(cfg.SmtpConfig, logger)

	//init http client
	cameoMetricsHttpClient, err := helpers.NewHttpClient(common.NewHttpClientConfig(cfg.CameoMetricsHttpClient.URL, cfg.CameoMetricsHttpClient.RateLimiting))
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed initialisation httpClient: %v", err))
	}

	// init storage
	userStorage := dao.NewUserStorage(db)
	sessionStorage := dao.NewSessionStorage(db)
	fileStorage := dao.NewFileStorage(db)
	companyStorage := dao.NewCompanyStorage(db)
	vacancyStorage := dao.NewVacancyStorage(db)
	candidateStorage := dao.NewCandidateStorage(db)

	// init service
	authService := service.NewAuthService(
		userStorage,
		sessionStorage,
		hasher,
		jwtManager,
		logger)

	userService := service.NewUserService(
		userStorage,
		authService,
		hasher,
		uuid.MustParse(cfg.AdminMigration.AdminID))

	candidateService := service.NewCandidateService(
		logger,
		vacancyStorage,
		candidateStorage,
		cameoMetricsHttpClient)

	vacancyService := service.NewVacancyService(logger, vacancyStorage, companyStorage, candidateService)

	companyService := service.NewCompanyService(logger, companyStorage, userService, vacancyService, fileStorage, minioService)
	// init controller
	controllers := controller.NewControllerContainer(
		logger,
		authService,
		userService,
		companyService,
		vacancyService,
		candidateService,
	)

	newDataProcessing := dataProcessing.NewDataProcessing("created_at", "ASC", 10)
	// init server
	handler := router.NewRouter(cfg)

	publicServer := new(server.Server)

	go func() {
		if err := publicServer.Run(cfg.Server.Host, cfg.Server.Port, handler.InitRoutes(
			logger,
			controllers,
			newDataProcessing,
			jwtManager,
			uuid.MustParse(cfg.AdminMigration.AdminID))); err != nil {
			logger.Error(fmt.Sprintf("error accured while running http server: %s", err.Error()))
		}
	}()

	logger.Info(fmt.Sprintf("listening public on %s:%s", cfg.Server.Host, cfg.Server.Port))

	// handle signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info("shutting down gracefully...")
	defer func() { logger.Info("shutdown complete") }()

	// perform shutdown
	if err := publicServer.Shutdown(context.Background()); err != nil {
		logger.Error(fmt.Sprintf("error occured on public server shutting down: %s", err.Error()))
	}
}
