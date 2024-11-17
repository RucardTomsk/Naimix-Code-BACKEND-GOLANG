package entity

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing/filter"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/enum"
	"github.com/google/uuid"
)

type Candidate struct {
	base.EntityWithIdKey
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	SystemID  string    `json:"system_id"`
	Vacancy   Vacancy   `json:"vacancy"`
	VacancyID uuid.UUID `json:"vacancy_id"`
}

func (Candidate) FilteringRules() map[string]map[string]enum.ValidateType {
	return filter.GetFilterRules(
		base.EntityWithIdKey{},
		"candidates",
		map[string]map[string]enum.ValidateType{
			"candidates": {
				"name": enum.TYPE_STRING,
			},
		})
}
