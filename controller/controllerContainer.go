package controller

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/service"
	"go.uber.org/zap"
)

type Container struct {
	AuthController *AuthController
	UserController *UserController
}

func NewControllerContainer(
	logger *zap.Logger,
	authService *service.AuthService,
	userService *service.UserService,
) *Container {
	return &Container{
		AuthController: NewAuthController(logger, authService),
		UserController: NewUserController(logger, userService),
	}
}
