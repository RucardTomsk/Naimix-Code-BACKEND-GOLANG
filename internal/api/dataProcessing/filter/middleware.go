package filter

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/enum"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

const (
	OptionsContextKey = "filter_options"
)

/*
The addFilterArgs function takes a pointer to an Options struct, a fieldName string, and a gin.Context object as input parameters.
It parses the query string in the context for the specified fieldName, and extracts the filter parameters for that field.
It then adds each filter parameter to the Options struct using the AddField method.
The function handles various comparison operators such as =, !=, <, <=, >, and >=, by replacing them with the corresponding comparison operator in the AddField method.
*/
func addFilterArgs(options *Options, fieldName string, c *gin.Context) {
	filterParamArray := c.QueryArray(fieldName)

	if len(filterParamArray) != 0 {
		for _, filterParam := range filterParamArray {
			if strings.Contains(filterParam, OperatorEquality) {
				filterParam = strings.ReplaceAll(filterParam, OperatorEquality, "")
				options.AddField(fieldName, filterParam, "=")
			} else if strings.Contains(filterParam, OperatorNotEquality) {
				filterParam = strings.ReplaceAll(filterParam, OperatorNotEquality, "")
				options.AddField(fieldName, filterParam, "!=")
			} else if strings.Contains(filterParam, OperatorLowerThan) {
				filterParam = strings.ReplaceAll(filterParam, OperatorLowerThan, "")
				options.AddField(fieldName, filterParam, "<")
			} else if strings.Contains(filterParam, OperatorLowerThanEq) {
				filterParam = strings.ReplaceAll(filterParam, OperatorLowerThanEq, "")
				options.AddField(fieldName, filterParam, "<=")
			} else if strings.Contains(filterParam, OperatorGreaterThan) {
				filterParam = strings.ReplaceAll(filterParam, OperatorGreaterThan, "")
				options.AddField(fieldName, filterParam, ">")
			} else if strings.Contains(filterParam, OperatorGreaterThanEq) {
				filterParam = strings.ReplaceAll(filterParam, OperatorGreaterThanEq, "")
				options.AddField(fieldName, filterParam, ">=")
			} else if strings.Contains(filterParam, OperatorLike) {
				filterParam = strings.ReplaceAll(filterParam, OperatorLike, "")
				options.AddField(fieldName, filterParam, "ILIKE")
			} else if strings.Contains(filterParam, OperatorIN) {
				filterParam = strings.ReplaceAll(filterParam, OperatorIN, "")
				options.AddField(fieldName, filterParam, "IN")
			} else {
				options.AddField(fieldName, filterParam, "=")
			}
		}
	}
}

// ParseFilterArgument
//
// The Parse Filter Argument function is a middleware for a Gin web framework that handles filter arguments received in the URL query parameters.
// It takes in a logger and a map of field names and types to validate the input.
// The middleware parses the URL and extracts the filter arguments from it.
// Then, it adds the filter arguments to an options object and validates the fields against the given validation model.
// If there is an error with the validation, the middleware returns an HTTP 400 bad request error with a specific message.
// Otherwise, the middleware sets the options in the Gin context for later use by other handlers.
// The middleware is intended to be used as part of a chain of middleware functions in the Gin web framework.
func ParseFilterArgument(logger zap.Logger, filterRules map[string]map[string]enum.ValidateType) gin.HandlerFunc {
	return func(c *gin.Context) {
		options := Options{}
		options.IsToApply = true

		argsURL := strings.Split(c.Request.RequestURI, "?")
		if len(argsURL) == 1 {
			options.IsToApply = false
			c.Set(OptionsContextKey, options)
			return
		}

		argsMas := strings.Split(argsURL[1], "&")

		var fieldNameString strings.Builder
		for _, rules := range filterRules {
			for nameColumn := range rules {
				addFilterArgs(&options, nameColumn, c)
				fieldNameString.WriteString(nameColumn)
			}
		}

		fieldNameString.WriteString("sort")
		fieldNameString.WriteString("limit")
		fieldNameString.WriteString("page")

		for index := range argsMas {
			argsName := strings.Split(argsMas[index], "=")[0]
			if !strings.Contains(fieldNameString.String(), argsName) {
				c.AbortWithStatusJSON(http.StatusBadRequest, api.GeneralFilterError())
				return
			}
		}

		if len(options.Fields) == 0 {
			options.IsToApply = false
		} else {
			if err := options.ValidateField(filterRules); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, api.GeneralFilterError())
				return
			}
		}

		c.Set(OptionsContextKey, options)
	}
}
