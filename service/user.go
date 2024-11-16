package service

import (
	"context"
	"errors"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/auth"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/model"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/storage/dao"
	"github.com/google/uuid"
	"net/http"
)

type UserService struct {
	userStorage *dao.UserStorage
	authService *AuthService
	hasher      *auth.Hasher
	adminID     uuid.UUID
}

func NewUserService(
	userStorage *dao.UserStorage,
	authService *AuthService,
	hasher *auth.Hasher,
	adminID uuid.UUID) *UserService {
	return &UserService{
		userStorage: userStorage,
		authService: authService,
		hasher:      hasher,
		adminID:     adminID,
	}
}

func (s *UserService) Get(option *dataProcessing.Options, ctx context.Context) ([]model.UserObject, int64, *base.ServiceError) {

	users, total, err := s.userStorage.Get(option, ctx)
	if err != nil {
		return nil, total, base.NewPostgresReadError(err)
	}

	result := make([]model.UserObject, 0, len(users))

	for _, user := range users {
		result = append(result, model.UserObject{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Name:      user.Name,
			Email:     user.Email,
		})
	}

	return result, total, nil
}

func (s *UserService) RetrieveUser(id uuid.UUID, ctx context.Context) (*model.UserObject, *base.ServiceError) {

	user, err := s.userStorage.Retrieve(id, ctx)
	if err != nil {
		return nil, base.NewPostgresReadError(err)
	}

	var isAdmin bool

	if id == s.adminID {
		isAdmin = true
	} else {
		isAdmin = false
	}

	return &model.UserObject{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		IsAdmin:   isAdmin,
		Email:     user.Email,
	}, nil
}

func (s *UserService) Update(id uuid.UUID, request model.UpdateUserAllFieldRequest, ctx context.Context) (mainErr *base.ServiceError) {
	user, err := s.userStorage.Retrieve(id, ctx)
	if err != nil {
		return base.NewPostgresReadError(err)
	}

	user.Email = request.Email
	if request.Password != nil {
		user.Password, err = s.hasher.Hash(*request.Password)
	}

	if err != nil {
		return base.NewUnauthorizedError(err)
	}

	user.Name = request.FIO

	if err := s.userStorage.Update(user, ctx); err != nil {
		return base.NewPostgresWriteError(err)
	}

	return nil
}

func (s *UserService) UpdateAuthorizationFields(id uuid.UUID, request model.UpdateUserAuthorizationFieldsRequest, ctx context.Context) (mainErr *base.ServiceError) {
	hashOldPassword, err := s.hasher.Hash(*request.OldPassword)
	if err != nil {
		return base.NewUnauthorizedError(err)
	}

	user, err := s.userStorage.Retrieve(id, ctx)
	if err != nil {
		return base.NewPostgresReadError(err)
	}

	if user.Password != hashOldPassword {
		return &base.ServiceError{
			Err:     errors.New("authentication failed. Please provide valid credentials"),
			Message: "authentication failed. Please provide valid credentials",
			Blame:   base.BlameUser,
			Code:    http.StatusUnauthorized,
		}
	}

	hashNewPassword, err := s.hasher.Hash(request.NewPassword)
	if err != nil {
		return base.NewUnauthorizedError(err)
	}

	user.Email = request.Email
	user.Password = hashNewPassword

	if err := s.userStorage.Update(user, ctx); err != nil {
		return base.NewPostgresWriteError(err)
	}

	if serviceErr := s.authService.SignOutAllSession(user.ID, ctx); serviceErr != nil {
		return serviceErr
	}

	return nil
}

func (s *UserService) GetUsersById(ids []uuid.UUID, ctx context.Context) ([]model.UserObject, *base.ServiceError) {
	users, err := s.userStorage.GetByIDList(ids, ctx)
	if err != nil {
		return nil, base.NewPostgresReadError(err)
	}
	var result []model.UserObject
	for _, user := range users {
		result = append(result, model.UserObject{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Name:      user.Name,
			Email:     user.Email,
		})
	}
	return result, nil
}

func (s *UserService) DeleteUser(id uuid.UUID, ctx context.Context) *base.ServiceError {
	if err := s.userStorage.DeleteUser(id, ctx); err != nil {
		return base.NewPostgresReadError(err)
	}

	return nil
}
