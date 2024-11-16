package pagination

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

const (
	OptionsContextKey = "pagination_options"
	pageKey           = "page"
	limitKey          = "limit"
)

// ParsePaginationArgument
//
// This is a middleware function that handles pagination arguments for a Gin web server.
// It takes a logger and a default limit as input parameters and returns a Gin handler function.
//
// The middleware function extracts the limit and page query parameters from the URL query string.
// If these parameters are not provided or their values are equal to their default values, the middleware sets the IsToApply field of the Options struct to false and returns.
//
// If the limit and page parameters are provided and have valid values, they are parsed to integers.
// The middleware constructs an Options struct with the IsToApply field set to true, and the limit and offset fields set according to the parsed values.
//
// Finally, the middleware sets the Options struct in the Gin context so that it can be accessed by subsequent handlers.
// If there are errors in parsing the parameters, the middleware logs a warning and returns a 400 HTTP response with a relevant error message.
func ParsePaginationArgument(logger zap.Logger, defaultLimit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		limitStr := c.DefaultQuery(limitKey, strconv.Itoa(defaultLimit))
		pageStr := c.DefaultQuery(pageKey, "1")

		if !strings.Contains(c.Request.RequestURI, "page") || !strings.Contains(c.Request.RequestURI, "limit") {
			c.Set(OptionsContextKey, Options{
				IsToApply: false,
			})
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, api.GeneralPaginationError())
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, api.GeneralPaginationError())
			return
		}

		options := Options{
			IsToApply: true,
			Limit:     limit,
			Offset:    limit * (page - 1),
		}

		c.Set(OptionsContextKey, options)
	}
}
