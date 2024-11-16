package dataProcessing

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing/filter"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing/pagination"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/api/dataProcessing/sort"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Options struct {
	FilterOptions     *filter.Options
	SortOptions       *sort.Options
	PaginationOptions *pagination.Options
}

// GetOptions returns Options from context.
func GetOptions(c *gin.Context) *Options {
	filterOptionsCtx, ok := c.Get(filter.OptionsContextKey)
	if !ok {
		newFilterOptions := filter.Options{}
		newFilterOptions.IsToApply = false
		filterOptionsCtx = newFilterOptions
	}
	filterOptions := filterOptionsCtx.(filter.Options)

	sortOptionsCtx, _ := c.Get(sort.OptionsContextKey)
	sortOptions := sortOptionsCtx.(sort.Options)

	paginationOptionsCtx, _ := c.Get(pagination.OptionsContextKey)
	paginationOptions := paginationOptionsCtx.(pagination.Options)

	return &Options{
		FilterOptions:     &filterOptions,
		SortOptions:       &sortOptions,
		PaginationOptions: &paginationOptions,
	}
}

func GetDefaultOptions() *Options {
	return &Options{
		FilterOptions: &filter.Options{
			IsToApply: false,
		},
		SortOptions: &sort.Options{
			Field: "id",
			Order: "ASC",
		},
		PaginationOptions: &pagination.Options{
			IsToApply: false,
		},
	}
}

// UseProcessing use all Options function
func (o *Options) UseProcessing(tx *gorm.DB) (*gorm.DB, int64, error) {
	if o.FilterOptions.IsToApply {
		mapConditionsFilter, err := o.FilterOptions.CreateConditionsFilter()
		if err != nil {
			return nil, 0, err
		}
		for nameOperator, value := range mapConditionsFilter {
			tx.Where(nameOperator, value)
		}
	}

	var count int64
	if err := tx.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	flag := true
	if o.SortOptions.Order == "ASC" {
		flag = false
	}
	tx.Order(clause.OrderByColumn{Column: clause.Column{Name: o.SortOptions.Field}, Desc: flag})
	if o.PaginationOptions.IsToApply {

		tx.Limit(o.PaginationOptions.Limit).Offset(o.PaginationOptions.Offset)
	}

	return tx, count, nil
}
