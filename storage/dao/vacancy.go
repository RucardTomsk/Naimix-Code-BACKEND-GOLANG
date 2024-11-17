package dao

import (
	"context"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VacancyStorage struct {
	db *gorm.DB
}

func NewVacancyStorage(db *gorm.DB) *VacancyStorage {
	return &VacancyStorage{db}
}

func (s VacancyStorage) Create(company *entity.Vacancy, ctx context.Context) error {
	return s.db.WithContext(ctx).Create(company).Error
}

func (s VacancyStorage) Retrieve(id uuid.UUID, ctx context.Context) (*entity.Vacancy, error) {
	var company entity.Vacancy
	err := s.db.WithContext(ctx).Preload("Candidates").First(&company, id).Error
	return &company, err
}

func (s VacancyStorage) Update(user *entity.Vacancy, ctx context.Context) error {
	return s.db.WithContext(ctx).Updates(user).Error
}

func (s VacancyStorage) Get(options *dataProcessing.Options, ctx context.Context) ([]entity.Vacancy, int64, error) {
	var users []entity.Vacancy
	tx := s.db.WithContext(ctx).Model(&entity.Vacancy{}).Preload("Candidates")

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
