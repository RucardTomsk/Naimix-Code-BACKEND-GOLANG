package service

import (
	"context"
	"errors"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/auth"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/entity"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/model"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/storage/dao"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

type AuthService struct {
	storage        *dao.UserStorage
	storageSession *dao.SessionStorage
	hasher         *auth.Hasher
	jwtManager     *auth.JWTManager
	logger         *zap.Logger
}

func NewAuthService(
	storage *dao.UserStorage,
	storageSession *dao.SessionStorage,
	hasher *auth.Hasher,
	jwtManager *auth.JWTManager,
	logger *zap.Logger) *AuthService {
	return &AuthService{
		storage:        storage,
		storageSession: storageSession,
		hasher:         hasher,
		jwtManager:     jwtManager,
		logger:         logger,
	}
}

func (s *AuthService) Register(request *model.RegisterRequest, ctx context.Context) (*uuid.UUID, *base.ServiceError) {
	hashPassword, err := s.hasher.Hash(request.Password)
	if err != nil {
		return nil, base.NewReadByteError(err)
	}

	user := &entity.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hashPassword,
	}

	if err := s.storage.Create(user, ctx); err != nil {
		return nil, base.NewPostgresWriteError(err)
	}

	return &user.ID, nil
}

func (s *AuthService) Login(request *model.LoginRequest, deviceId string, ctx context.Context) (jwt *string, refToken *uuid.UUID, mainErr *base.ServiceError) {
	user, err := s.storage.GetUser(request.Email, ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Info(request.Email + ": user not found")
			return nil, nil, base.NewLoginError(err)
		}
		return nil, nil, base.NewPostgresReadError(err)
	}
	hashPassword, err := s.hasher.Hash(request.Password)
	if err != nil {
		return nil, nil, base.NewReadByteError(err)
	}

	if hashPassword != user.Password {
		s.logger.Info(request.Email + ": user invalid password")
		return nil, nil, base.NewLoginError(err)
	}

	session := &entity.Session{
		UserID:   user.ID,
		DeviceID: deviceId,
	}

	if err := s.storageSession.Create(session, ctx); err != nil {
		return nil, nil, base.NewPostgresWriteError(err)
	}

	s.logger.Info(user.Email + ": session create")

	token, err := s.jwtManager.NewJWT(user.ID)
	if err != nil {
		return nil, nil, base.NewCreateJWTError(err)
	}

	return &token, &session.ID, nil
}

func (s *AuthService) SignOutAllSession(id uuid.UUID, ctx context.Context) (mainErr *base.ServiceError) {
	sessions, err := s.storageSession.GetByUserID(id, ctx)
	if err != nil {
		return base.NewPostgresReadError(err)
	}

	for _, session := range sessions {
		if serviceErr := s.storageSession.Delete(session.ID, ctx); serviceErr != nil {
			return base.NewPostgresReadError(err)
		}
	}

	return nil
}

func (s *AuthService) Logout(refreshToken uuid.UUID, ctx context.Context) (mainErr *base.ServiceError) {
	session, err := s.storageSession.Retrieve(refreshToken, ctx)
	if err != nil {
		return &base.ServiceError{
			Err:     err,
			Blame:   base.BlameUser,
			Code:    http.StatusUnauthorized,
			Message: "failed get session",
		}
	}
	if err := s.storageSession.Delete(session.ID, ctx); err != nil {
		return base.NewPostgresReadError(err)
	}

	return nil
}

func (s *AuthService) RefreshJWT(id uuid.UUID, ctx context.Context) (*string, *uuid.UUID, *uuid.UUID, *base.ServiceError) {
	session, err := s.storageSession.Retrieve(id, ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, base.NewNotSessionError(err)
		}
		return nil, nil, nil, base.NewPostgresReadError(err)
	}

	newSession := &entity.Session{
		UserID: session.UserID,
	}

	if err := s.storageSession.Delete(session.ID, ctx); err != nil {
		return nil, nil, nil, base.NewPostgresReadError(err)
	}

	if err := s.storageSession.Create(newSession, ctx); err != nil {
		return nil, nil, nil, base.NewPostgresWriteError(err)
	}

	token, err := s.jwtManager.NewJWT(session.UserID)
	if err != nil {
		return nil, nil, nil, base.NewCreateJWTError(err)
	}

	return &token, &newSession.ID, &newSession.UserID, nil
}
