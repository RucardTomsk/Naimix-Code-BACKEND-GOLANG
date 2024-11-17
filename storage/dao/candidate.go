package dao

import (
	"context"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CandidateStorage struct {
	db *gorm.DB
}

func NewCandidateStorage(db *gorm.DB) *CandidateStorage {
	return &CandidateStorage{db}
}

func (s CandidateStorage) Create(company *entity.Candidate, ctx context.Context) error {
	return s.db.WithContext(ctx).Create(company).Error
}

func (s CandidateStorage) Retrieve(id uuid.UUID, ctx context.Context) (*entity.Candidate, error) {
	var company entity.Candidate
	err := s.db.WithContext(ctx).First(&company, id).Error
	return &company, err
}

func (s CandidateStorage) Update(user *entity.Candidate, ctx context.Context) error {
	return s.db.WithContext(ctx).Updates(user).Error
}

func (s CandidateStorage) Get(options *dataProcessing.Options, ctx context.Context) ([]entity.Candidate, int64, error) {
	var users []entity.Candidate
	tx := s.db.WithContext(ctx).Model(&entity.Candidate{})

	tx, total, err := options.UseProcessing(tx)
	if err != nil {
		return nil, total, err
	}

	tx.Find(&users)
	if tx.Error != nil {
		return nil, total, tx.Error
	}

	return users, total, nil
}
