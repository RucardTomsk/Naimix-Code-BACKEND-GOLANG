package base

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// EntityWithIdKey is a base DB entity with uuid.UUID as a primary key.
type EntityWithIdKey struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v1();primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

// EntityWithIdKeyUniqueIndex is a base DB entity with uuid.UUID as a unique index.
type EntityWithIdKeyUniqueIndex struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v1();uniqueIndex"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

// ArchivableEntityWithIdKey is an EntityWithGuidKey struct with extra ArchivedAt field.
type ArchivableEntityWithIdKey struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v1();primaryKey"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	ArchivedAt gorm.DeletedAt `json:"-"`
}

// EntityWithIntegerKey is a base DB entity with uint as a primary key.
type EntityWithIntegerKey struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
