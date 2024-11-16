package common

import (
	"github.com/gin-contrib/cors"
	"net/http"
	"time"
)

// ServerConfig configures gin server.
type ServerConfig struct {
	Host             string
	Port             string
	UseAuthorization bool

	GinMode string
}

// DatabaseConfig stores DB credentials.
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type MainMail struct {
	Mail string
}

// MinioConfig is used to connect to minio (s3).
type MinioConfig struct {
	UseMocks  bool
	Endpoint  string
	AccessKey string
	SecretKey string
	Token     string
	UseSSL    bool
}

const (
	defaultHost     = "localhost"
	defaultBasePath = "/api"
)

var defaultSchemes = []string{"http", "https"}

// SwaggerConfig configures swaggo/swag.
type SwaggerConfig struct {
	Title       string
	Description string
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
}

// NewSwaggerConfig returns *SwaggerConfig with preconfigured fields.
func NewSwaggerConfig(title, description, version string) *SwaggerConfig {
	return &SwaggerConfig{
		Title:       title,
		Description: description,
		Version:     version,
		Host:        defaultHost,
		BasePath:    defaultBasePath,
		Schemes:     defaultSchemes,
	}
}

// DefaultCorsConfig returns cors.Config with very permissive policy.
func DefaultCorsConfig() cors.Config {
	return cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPatch, http.MethodPut, http.MethodOptions},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"*"},
		MaxAge:        12 * time.Hour,
	}
}

// DataProcessingConfig configures default sort, order and pagination parameters.
type DataProcessingConfig struct {
	DefaultSortField string
	DefaultSortOrder string
	DefaultLimit     int
}

// NewDataProcessingConfig returns *DataProcessingConfig with preconfigured fields.
func NewDataProcessingConfig(
	defaultSortField string,
	defaultSortOrder string,
	defaultLimit int,
) *DataProcessingConfig {
	return &DataProcessingConfig{
		DefaultSortField: defaultSortField,
		DefaultSortOrder: defaultSortOrder,
		DefaultLimit:     defaultLimit,
	}
}

type HttpClientConfig struct {
	URL          string
	RateLimiting int
}

func NewHttpClientConfig(
	url string,
	rateLimiting int) *HttpClientConfig {
	return &HttpClientConfig{URL: url, RateLimiting: rateLimiting}
}

type AuthConfig struct {
	Salt       string
	SigningKey string
	TimeToLive time.Duration
}

type AdminMigrationConfig struct {
	AdminID       string
	AdminUserName string
	AdminEmail    string
	AdminPassword string
}

type SmtpConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}
