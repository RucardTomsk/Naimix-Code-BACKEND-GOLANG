package dao

import (
	"context"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CompanyStorage struct {
	db *gorm.DB
}

func NewCompanyStorage(db *gorm.DB) *CompanyStorage {
	return &CompanyStorage{db}
}

func (s CompanyStorage) Create(company *entity.Company, ctx context.Context) error {
	return s.db.WithContext(ctx).Create(company).Error
}

func (s CompanyStorage) Retrieve(id uuid.UUID, ctx context.Context) (*entity.Company, error) {
	var company entity.Company
	err := s.db.WithContext(ctx).Preload("File").Preload("Users").Preload("Vacancies").First(&company, id).Error
	return &company, err
}

func (s CompanyStorage) Update(user *entity.Company, ctx context.Context) error {
	return s.db.WithContext(ctx).Updates(user).Error
}

func (s CompanyStorage) Get(options *dataProcessing.Options, ctx context.Context) ([]entity.Company, int64, error) {
	var users []entity.Company
	tx := s.db.WithContext(ctx).Model(&entity.Company{}).Preload("File").Preload("Users").Preload("Vacancies")

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
