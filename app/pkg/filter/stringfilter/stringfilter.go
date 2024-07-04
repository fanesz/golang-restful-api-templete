package stringfilter

import (
	"backend/app/common/consts"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StringFilter struct {
	OrderBy consts.OrderBy `json:"order_by" form:"order_by"`
	Sort    string         `json:"sort" form:"sort"`
}

func Build(c *gin.Context, query *gorm.DB, fieldName string, filter *StringFilter) {
	if filter.OrderBy != consts.ASC && filter.OrderBy != consts.DESC {
		return
	}

	if filter.Sort != fieldName && filter.OrderBy != consts.ASC && filter.OrderBy != consts.DESC {
		return
	}

	query.Order(gorm.Expr("LOWER(?) ?", fieldName, filter.OrderBy))
}
