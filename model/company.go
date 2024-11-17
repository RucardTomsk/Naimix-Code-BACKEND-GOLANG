package model

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/google/uuid"
	"time"
)

type CompanyObject struct {
	ID          uuid.UUID       `json:"id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	LogoURL     string          `json:"logoUrl"`
	Users       []UserObject    `json:"users" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Vacancies   []VacancyObject `json:"vacancies" gorm:"constraint:OnUpdate:CASCADE;"`
}

type (
	CreateCompanyRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	RetrieveCompanyResponse struct {
		base.ResponseOK
		Company CompanyObject `json:"company"`
	}

	GetCompanyResponse struct {
		base.ResponseOK
		Companies []CompanyObject `json:"companies"`
	}
)
