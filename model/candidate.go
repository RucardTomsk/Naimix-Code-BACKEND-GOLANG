package model

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/google/uuid"
	"time"
)

type CandidateObject struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	SystemID  string    `json:"system_id"`
}

type (
	AddNewCandidateRequest struct {
		Name  string `json:"name"`
		Email string `json:"email" gorm:"uniqueIndex"`
	}

	RetrieveCandidateResponse struct {
		base.ResponseOK
		Candidate CandidateObject `json:"candidate"`
	}

	GetCandidatesResponse struct {
		base.ResponseOK
		Candidates []CandidateObject `json:"candidates"`
	}
)
