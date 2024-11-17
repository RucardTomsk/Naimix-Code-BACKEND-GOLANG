package dao

import (
	"context"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/entity"
	"gorm.io/gorm"
)

type FileStorage struct {
	db *gorm.DB
}

func NewFileStorage(db *gorm.DB) *FileStorage {
	return &FileStorage{db}
}

func (s *FileStorage) Create(file *entity.File, ctx context.Context) error {
	return s.db.WithContext(ctx).Create(file).Error
}
