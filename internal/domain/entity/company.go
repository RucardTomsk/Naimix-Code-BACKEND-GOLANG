package entity

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing/filter"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/enum"
	"github.com/google/uuid"
)

type Company struct {
	base.EntityWithIdKey

	Name        string     `json:"name"`
	Description string     `json:"description"`
	Owner       uuid.UUID  `json:"owner"`
	FileID      *uuid.UUID `json:"avatar_id"`
	File        *File
	Users       []User    `json:"users" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Vacancies   []Vacancy `json:"vacancies" gorm:"constraint:OnUpdate:CASCADE;"`
}

func (Company) FilteringRules() map[string]map[string]enum.ValidateType {
	return filter.GetFilterRules(
		base.EntityWithIdKey{},
		"companies",
		map[string]map[string]enum.ValidateType{
			"companies": {
				"name": enum.TYPE_STRING,
			},
		})
}
