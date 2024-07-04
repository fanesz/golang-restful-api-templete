package user

import (
	"backend/app/pkg/filter/datefilter"
	"backend/app/pkg/filter/singlesearch"
	"backend/app/pkg/filter/stringfilter"
	"backend/app/pkg/handler"
	"backend/app/pkg/pagination"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserFilters struct {
	Pagination     handler.Pagination
	DateFilter     datefilter.DateFilter
	StringFilter   stringfilter.StringFilter
	UsernameFilter searchByUsername
}

func filterService(c *gin.Context, query *gorm.DB, filters *UserFilters) {
	datefilter.Build(c, query, &filters.DateFilter)
	stringfilter.Build(c, query, "username", &filters.StringFilter)
	singlesearch.Build(c, query, "username", &filters.UsernameFilter.Username)

	pagination.Build(c, query, &filters.Pagination)
}

type searchByUsername struct {
	Username string `form:"username"`
}
