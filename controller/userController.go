package controller

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/middleware"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/model"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

type UserController struct {
	logger      *zap.Logger
	userService *service.UserService
}

func NewUserController(
	logger *zap.Logger,
	userService *service.UserService) *UserController {
	return &UserController{
		logger:      logger,
		userService: userService,
	}
}

// Get user-api
// @Summary      Get all users
// @Description  Get all users
// @Tags         User
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200  {object}  model.GetUsersResponse "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /user/get [get]
func (a *UserController) Get(c *gin.Context) {
	users, total, serviceErr := a.userService.Get(dataProcessing.GetOptions(c), c)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr))
		return
	}

	c.JSON(http.StatusOK, model.GetUsersResponse{
		ResponseOK: base.ResponseOK{
			Status: http.StatusText(http.StatusOK),
		},
		Total: total,
		Users: users,
	})
}

// DeleteUser user-api
// @Summary      Delete User
// @Description  Delete User
// @Tags         User
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        user-id path string true "User ID"
// @Success      200  {object}  base.ResponseOK "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /user/{user-id}/delete [delete]
func (a *UserController) DeleteUser(c *gin.Context) {
	adminID, _ := c.Get(middleware.AdminsID)
	userID, err := uuid.Parse(c.Params.ByName("user-id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, api.GeneralParsingError())
		return
	}

	if adminID.(uuid.UUID) == userID {
		c.JSON(http.StatusOK, base.ResponseFailure{
			Status:  http.StatusText(http.StatusForbidden),
			Blame:   base.BlameUser,
			Message: "no access",
		})
		return
	}

	if serviceErr := a.userService.DeleteUser(userID, c); serviceErr != nil {
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr))
		return
	}

	c.JSON(http.StatusOK, base.ResponseOK{
		Status: http.StatusText(http.StatusOK),
	})
}

// Update
// @Summary      Update User All Fields
// @Description  Update User All Fields
// @Tags         User
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        payload body   model.UpdateUserAllFieldRequest true "User data"
// @Success      200  {object}  base.ResponseOK "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /user/field/update [post]
func (a *UserController) Update(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)

	var payload model.UpdateUserAllFieldRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, base.ResponseFailure{
			Status:  http.StatusText(http.StatusBadRequest),
			Blame:   base.BlameUser,
			Message: "failed to parse json",
		})
		return
	}

	if serviceErr := a.userService.Update(userID.(uuid.UUID), payload, c); serviceErr != nil {
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr))
		return
	}

	c.JSON(http.StatusOK, base.ResponseOK{
		Status: http.StatusText(http.StatusOK),
	})
}

// RetrieveUser user-api
// @Summary     Retrieve data of an authorised user
// @Description  Retrieve data of an authorised user
// @Tags         User
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200  {object}  model.GetUserResponse "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /user/retrieve [get]
func (a *UserController) RetrieveUser(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)
	user, serviceErr := a.userService.RetrieveUser(userID.(uuid.UUID), c)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr))
		return
	}

	c.JSON(http.StatusOK, model.GetUserResponse{
		ResponseOK: base.ResponseOK{
			Status: http.StatusText(http.StatusOK),
		},
		User: *user,
	})
}

// GetUserByIdList private-user-api
// @Summary      Retrieve user information by id list
// @Description  Retrieve user information by id list
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        payload body   model.UsersByIdListRequest true "User data"
// @Success      200  {object}  model.GetUsersResponse "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /usersByIdList [post]
func (a *UserController) GetUserByIdList(c *gin.Context) {
	var payload model.UsersByIdListRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, base.ResponseFailure{
			Status:  http.StatusText(http.StatusBadRequest),
			Blame:   base.BlameUser,
			Message: "failed to parse json",
		})
		return
	}
	if len(payload.Ids) == 0 {
		c.JSON(http.StatusOK, model.GetUsersResponse{
			ResponseOK: base.ResponseOK{
				Status: http.StatusText(http.StatusOK),
			},
			Total: int64(len(payload.Ids)),
			Users: []model.UserObject{},
		})
	}

	res, serviceErr := a.userService.GetUsersById(payload.Ids, c)
	if serviceErr != nil {
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr))
		return
	}
	c.JSON(http.StatusOK, model.GetUsersResponse{
		ResponseOK: base.ResponseOK{
			Status: http.StatusText(http.StatusOK),
		},
		Total: int64(len(payload.Ids)),
		Users: res,
	})
}

func (a *UserController) GetUserById(c *gin.Context) {
	userId, err := uuid.Parse(c.Params.ByName("user-id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, api.GeneralParsingError())
		return
	}
	user, serviceError := a.userService.RetrieveUser(userId, c)
	if serviceError != nil {
		if serviceError.Code != 404 {
			c.JSON(serviceError.Code, serviceError)
			return
		}
		c.JSON(serviceError.Code, base.ResponseOK{
			Status: http.StatusText(http.StatusOK),
		})
		return
	}
	c.JSON(http.StatusOK, model.GetUserResponse{
		ResponseOK: base.ResponseOK{
			Status: http.StatusText(http.StatusOK),
		},
		User: *user,
	})
}
