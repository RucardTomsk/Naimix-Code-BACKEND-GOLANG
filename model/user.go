package model

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/google/uuid"
	"time"
)

type (
	UserObject struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		IsAdmin   bool      `json:"is_admin"`
	}
)

type (
	RegisterRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	UpdateUserAllFieldRequest struct {
		Email    string  `json:"email"`
		Password *string `json:"password"`
		FIO      string  `json:"fio"`
	}

	UpdateUserAuthorizationFieldsRequest struct {
		Email       string  `json:"email"`
		OldPassword *string `json:"old_password"`
		NewPassword string  `json:"new_password"`
	}

	UsersByIdListRequest struct {
		Ids []uuid.UUID `json:"ids"`
	}

	RecreateJWTRequest struct {
		RefreshToken uuid.UUID `json:"refresh_token"`
	}
)

type (
	GetUserResponse struct {
		base.ResponseOK
		User UserObject `json:"user"`
	}

	GetUsersResponse struct {
		base.ResponseOK
		Total int64        `json:"total"`
		Users []UserObject `json:"users"`
	}

	LoginResponse struct {
		base.ResponseOK
		JWT          string    `json:"token"`
		RefreshToken uuid.UUID `json:"refresh_token"`
	}
)
