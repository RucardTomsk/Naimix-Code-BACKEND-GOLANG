package middleware

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/auth"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	//UserIDKey value type is uuid.UUID
	UserIDKey = "userID"
	IsAdmin   = "isAdmin"
	AdminsID  = "adminID"
)

// SetAuthorizationCheck adds authorization check to middleware chain.
func SetAuthorizationCheck(JWTManager *auth.JWTManager, logger zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		header := c.GetHeader(authorizationHeader)
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, base.ResponseFailure{
				Status:  http.StatusText(http.StatusUnauthorized),
				Blame:   base.BlameUser,
				Message: "unauthorized",
			})
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, base.ResponseFailure{
				Status:  http.StatusText(http.StatusUnauthorized),
				Blame:   base.BlameUser,
				Message: "unauthorized",
			})
			return
		}

		stringUserID, err := JWTManager.Parse(headerParts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, base.ResponseFailure{
				Status:  http.StatusText(http.StatusUnauthorized),
				Blame:   base.BlameUser,
				Message: "unauthorized",
			})
			return
		}

		userID, err := uuid.Parse(stringUserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, base.ResponseFailure{
				Status:  http.StatusText(http.StatusUnauthorized),
				Blame:   base.BlameUser,
				Message: "unauthorized",
			})
			return
		}

		c.Set(UserIDKey, userID)
		c.Next()
	}
}

func SetAuthorizationAdminCheck(JWTManager *auth.JWTManager, adminID uuid.UUID, logger zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		header := c.GetHeader(authorizationHeader)
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, base.ResponseFailure{
				Status:  http.StatusText(http.StatusUnauthorized),
				Blame:   base.BlameUser,
				Message: "unauthorized",
			})
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, base.ResponseFailure{
				Status:  http.StatusText(http.StatusUnauthorized),
				Blame:   base.BlameUser,
				Message: "unauthorized",
			})
			return
		}

		stringUserID, err := JWTManager.Parse(headerParts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, base.ResponseFailure{
				Status:  http.StatusText(http.StatusUnauthorized),
				Blame:   base.BlameUser,
				Message: "unauthorized",
			})
			return
		}

		userID, err := uuid.Parse(stringUserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, base.ResponseFailure{
				Status:  http.StatusText(http.StatusUnauthorized),
				Blame:   base.BlameUser,
				Message: "unauthorized",
			})
			return
		}

		if userID != adminID {
			c.AbortWithStatusJSON(http.StatusForbidden, base.ResponseFailure{
				Status:  http.StatusText(http.StatusForbidden),
				Blame:   base.BlameUser,
				Message: "no access",
			})
		}

		c.Set(IsAdmin, true)
		c.Set(AdminsID, userID)
		c.Next()
	}
}
