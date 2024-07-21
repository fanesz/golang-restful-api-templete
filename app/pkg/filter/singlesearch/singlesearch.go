package singlesearch

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Build(c *gin.Context, query *gorm.DB, searchField string, searchValue *string) {
	if searchValue == nil || *searchValue == "" {
		return
	}

	field := "LOWER(" + searchField + ")"
	query.Where(field+" LIKE LOWER(?)", "%"+*searchValue+"%")
}
