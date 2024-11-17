package model

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/google/uuid"
	"time"
)

type VacancyObject struct {
	ID          uuid.UUID         `json:"id"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Name        string            `json:"name"`
	Salary      int               `json:"salary"`
	City        string            `json:"city"`
	Description string            `json:"description"`
	Candidates  []CandidateObject `json:"candidates" `
}

type (
	CreateNewVacancyRequest struct {
		Name        string `json:"name"`
		Salary      int    `json:"salary"`
		City        string `json:"city"`
		Description string `json:"description"`
	}

	RetrieveVacancyResponse struct {
		base.ResponseOK
		Vacancy VacancyObject `json:"vacancy"`
	}

	GetVacancyResponse struct {
		base.ResponseOK
		Vacancies []VacancyObject `json:"vacancies"`
	}
)
