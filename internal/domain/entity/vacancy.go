package entity

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing/filter"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/enum"
	"github.com/google/uuid"
)

type Vacancy struct {
	base.EntityWithIdKey
	Name        string      `json:"name"`
	Salary      int         `json:"salary"`
	City        string      `json:"city"`
	Description string      `json:"description"`
	Candidates  []Candidate `json:"candidates" `
	Company     Company     `json:"company"`
	CompanyID   uuid.UUID   `json:"company_id"`
}

func (Vacancy) FilteringRules() map[string]map[string]enum.ValidateType {
	return filter.GetFilterRules(
		base.EntityWithIdKey{},
		"vacancies",
		map[string]map[string]enum.ValidateType{
			"vacancies": {
				"name": enum.TYPE_STRING,
			},
		})
}
