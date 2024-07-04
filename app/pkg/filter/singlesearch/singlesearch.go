package singlesearch

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Build(c *gin.Context, query *gorm.DB, searchField string, searchValue *string) {
	if searchValue == nil || *searchValue == "" {
		return
	}

	query.Where("LOWER(username) LIKE LOWER(?)", "%"+*searchValue+"%")
}
