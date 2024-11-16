package dataProcessing

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing/filter"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing/pagination"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing/sort"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/common"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/enum"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"sync"
)

type DataProcessing struct {
	config *common.DataProcessingConfig
}

func NewDataProcessing(
	defaultSortField string,
	defaultSortOrder string,
	defaultLimit int,
) *DataProcessing {
	return &DataProcessing{
		config: common.NewDataProcessingConfig(defaultSortField, defaultSortOrder, defaultLimit),
	}
}

// ApplyMiddleware adds pagination, filter and sorting handlers to middleware chain.
func (p DataProcessing) ApplyMiddleware(
	logger zap.Logger,
	filterRules map[string]map[string]enum.ValidateType,
	sortRules map[string]enum.ValidateType,
) gin.HandlerFunc {
	handlers := gin.HandlersChain{
		pagination.ParsePaginationArgument(logger, p.config.DefaultLimit),
		sort.ParseSortingArgument(logger, p.config.DefaultSortField, p.config.DefaultSortOrder, sortRules),
	}

	if filterRules != nil {
		handlers = append(handlers, filter.ParseFilterArgument(logger, filterRules))
	}
	return func(c *gin.Context) {
		var wg sync.WaitGroup
		wg.Add(len(handlers))
		for _, h := range handlers {
			go func(handler gin.HandlerFunc) {
				handler(c)
				wg.Done()
			}(h)
		}
		wg.Wait()
		c.Next()
	}
}
