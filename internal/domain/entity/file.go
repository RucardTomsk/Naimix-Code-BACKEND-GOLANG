package entity

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/google/uuid"
)

// File is a general object stored in s3.
type File struct {
	base.EntityWithIdKey
	Key    uuid.UUID `json:"key" example:"00000000-0000-0000-0000-000000000000"`
	Bucket string    `json:"bucket"`
}
