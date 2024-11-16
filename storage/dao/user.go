package dao

import (
	"context"
	"errors"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (s UserStorage) Create(user *entity.User, ctx context.Context) error {
	return s.db.WithContext(ctx).Create(user).Error
}

func (s UserStorage) Retrieve(id uuid.UUID, ctx context.Context) (*entity.User, error) {
	var user entity.User
	err := s.db.WithContext(ctx).Preload("Sessions").First(&user, id).Error
	return &user, err
}

func (s UserStorage) DeleteUser(id uuid.UUID, ctx context.Context) error {
	return s.db.WithContext(ctx).Unscoped().Delete(&entity.User{}, id).Error
}

func (s UserStorage) Update(user *entity.User, ctx context.Context) error {
	return s.db.WithContext(ctx).Updates(user).Error
}

func (s UserStorage) Get(options *dataProcessing.Options, ctx context.Context) ([]entity.User, int64, error) {
	var users []entity.User
	tx := s.db.WithContext(ctx).Model(&entity.User{})

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

func (s UserStorage) GetUser(email string, ctx context.Context) (*entity.User, error) {
	var user entity.User
	tx := s.db.WithContext(ctx).Preload("Sessions").First(&user, "email = ?", email)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}

func (s UserStorage) GetByIDList(ids []uuid.UUID, ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	err := s.db.WithContext(ctx).Where(ids).Find(&users).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.User{}, nil
		}
		return []entity.User{}, err
	}
	return users, err
}
