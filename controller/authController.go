package controller

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/helpers"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/model"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type AuthController struct {
	logger  *zap.Logger
	service *service.AuthService
}

func NewAuthController(
	logger *zap.Logger,
	service *service.AuthService) *AuthController {
	return &AuthController{
		logger:  logger,
		service: service,
	}
}

// Register User registration user-api
// @Summary      User registration
// @Description  User registration
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param payload body model.RegisterRequest true "User request"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200  {object}  base.ResponseOKWithID "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /user/register [post]
func (a *AuthController) Register(c *gin.Context) {
	var payload model.RegisterRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, api.GeneralParsingError())
		return
	}

	id, serviceErr := a.service.Register(&payload, c)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr))
		return
	}

	c.JSON(http.StatusOK, base.ResponseOKWithID{
		Status: http.StatusText(http.StatusOK),
		ID:     *id,
	})
}

// Login User authorisation user-api
// @Summary      User authorisation
// @Description  User authorisation
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param payload body model.LoginRequest true "User request"
// @Success      200  {object}  model.LoginResponse "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /user/login [post]
func (a *AuthController) Login(c *gin.Context) {
	var payload model.LoginRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, api.GeneralParsingError())
		return
	}

	deviceID := helpers.GetDeviceID(c)

	token, refresh, serviceErr := a.service.Login(&payload, deviceID, c)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr))
		return
	}

	c.JSON(http.StatusOK, model.LoginResponse{
		ResponseOK: base.ResponseOK{
			Status: http.StatusText(http.StatusOK),
		},
		JWT:          *token,
		RefreshToken: *refresh,
	})
}

// Logout Unauthorized users user-api
// @Summary      Unauthorized users
// @Description  Unauthorized users
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param payload body model.RecreateJWTRequest true "User request"
// @Success      200  {object}  base.ResponseOK "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /user/logout [post]
func (a *AuthController) Logout(c *gin.Context) {

	var payload model.RecreateJWTRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, api.GeneralParsingError())
		return
	}

	if serviceErr := a.service.Logout(payload.RefreshToken, c); serviceErr != nil {
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr))
		return
	}

	c.JSON(http.StatusOK, base.ResponseOK{
		Status: http.StatusText(http.StatusOK),
	})
}

// RecreateJWT Re-create refresh token user-api
// @Summary      Re-create refresh token
// @Description  Re-create refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param payload body model.RecreateJWTRequest true "User request"
// @Success      200  {object}  model.LoginResponse "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /user/refresh [post]
func (a *AuthController) RecreateJWT(c *gin.Context) {
	var payload model.RecreateJWTRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, api.GeneralParsingError())
		return
	}

	token, newRefresh, _, serviceErr := a.service.RefreshJWT(payload.RefreshToken, c)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr))
		return
	}

	c.JSON(http.StatusOK, model.LoginResponse{
		ResponseOK: base.ResponseOK{
			Status: http.StatusText(http.StatusOK),
		},
		JWT:          *token,
		RefreshToken: *newRefresh,
	})
}
