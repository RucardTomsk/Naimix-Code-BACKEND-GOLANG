package entity

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing/filter"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/enum"
	"github.com/google/uuid"
)

// User represents general system user (student or teacher).
type User struct {
	base.EntityWithIdKey

	Name      string     `json:"name"`
	Email     string     `json:"email" gorm:"uniqueIndex"`
	Password  string     `json:"password"`
	Company   *Company   `json:"company"`
	CompanyID *uuid.UUID `json:"companyID"`
	Sessions  []Session  `json:"sessions,omitempty"`
}

func (User) FilteringRules() map[string]map[string]enum.ValidateType {
	return filter.GetFilterRules(
		base.EntityWithIdKey{},
		"users",
		map[string]map[string]enum.ValidateType{
			"users": {
				"name":  enum.TYPE_STRING,
				"email": enum.TYPE_STRING,
			},
		})
}

type Session struct {
	base.EntityWithIdKey
	UserID   uuid.UUID `json:"user_id"`
	User     *User     `json:"user"`
	DeviceID string    `json:"device_id"`
}
