package pagination

import (
	"backend/app/pkg/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	DEFAULT_LIMIT = 10
	DEFAULT_PAGE  = 1
)

func Build(c *gin.Context, query *gorm.DB, pagination *handler.Pagination) {
	if pagination.Limit == 0 {
		pagination.Limit = DEFAULT_LIMIT
	}
	if pagination.Page == 0 {
		pagination.Page = DEFAULT_PAGE
	}

	var count int64
	query.Count(&count)
	pagination.TotalItems = int(count)
	pagination.TotalPages = (pagination.TotalItems + pagination.Limit - 1) / pagination.Limit
	query.Limit(pagination.Limit).Offset((pagination.Page - 1) * pagination.Limit)
}
