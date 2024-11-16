package sort

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/enum"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

const (
	ASC               = "ASC"
	DESC              = "DESC"
	OptionsContextKey = "sort_options"
	sortArgsKey       = "sort"
)

// ParseSortingArgument
//
// Function is a middleware for a Gin web framework in Go.
// It takes in a logger, default sorting field and order, and a map of validated fields as input and returns a Gin middleware function.
//
// When a client sends a request to the server, this middleware parses the sorting parameters in the request's query string.
// If the query string doesn't contain sorting parameters, the middleware sets the default sorting field and order.
// If the field isn't a valid sorting field, the middleware sends a "Bad Request" status response to the client.
//
// The middleware then creates an Options struct that contains the sorting field and order and sets it in the request context using the OptionsContextKey.
//
// This middleware ensures that the server handles the client's sorting parameters correctly and consistently.
func ParseSortingArgument(
	logger zap.Logger,
	defaultSortField string,
	defaultSortOrder string,
	validateMap map[string]enum.ValidateType) gin.HandlerFunc {
	return func(c *gin.Context) {
		sortArgs := c.Query(sortArgsKey)
		sortBy := defaultSortField
		sortOrder := defaultSortOrder

		if !checkThereField(validateMap, sortBy) {
			c.AbortWithStatusJSON(http.StatusBadRequest, api.GeneralSortError())
			return
		}

		if sortArgs != "" {
			sortArgsArray := strings.Split(sortArgs, ".")
			sortBy = sortArgsArray[0]
			if len(sortArgsArray) == 2 {
				sortOrder = strings.ToUpper(sortArgsArray[1])
				if strings.ToUpper(sortOrder) != ASC && strings.ToUpper(sortOrder) != DESC {
					c.AbortWithStatusJSON(http.StatusBadRequest, api.GeneralSortError())
					return
				}
			}
		}

		options := Options{
			Field: sortBy,
			Order: sortOrder,
		}

		c.Set(OptionsContextKey, options)
	}
}

/*
The checkThereField function takes a map of fields to validate as well as a field name as input and returns a boolean value.
It checks whether the given field name is present in the map of fields to validate, and in the EntityFiledString constant string.

The function builds a string that contains all the field names from the map of fields to validate, and the EntityFiledString string.
Then, it uses the strings.Contains function to check if the given field name is present in the built string.
If the field is present, the function returns true; otherwise, it returns false.
*/
func checkThereField(validateMap map[string]enum.ValidateType, fieldName string) bool {
	var validateString strings.Builder
	for nameField := range validateMap {
		validateString.WriteString(nameField)
	}
	validateString.WriteString("id created_at updated_at deleted_at archived_at")
	return strings.Contains(validateString.String(), fieldName)
}
