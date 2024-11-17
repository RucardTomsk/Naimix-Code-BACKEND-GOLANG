package config

import "github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/common"

type Config struct {
	DB                     common.DatabaseConfig
	Server                 common.ServerConfig
	Auth                   common.AuthConfig
	AdminMigration         common.AdminMigrationConfig
	Minio                  common.MinioConfig
	Mail                   common.MainMail
	SmtpConfig             common.SmtpConfig
	CameoMetricsHttpClient common.HttpClientConfig
}
