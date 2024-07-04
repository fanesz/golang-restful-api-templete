package validator

import (
	"backend/app/pkg/handler"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Query(query *gorm.DB, c *gin.Context, count bool) (bool, int64) {
	if query.Error != nil {
		handler.Error(c, http.StatusInternalServerError, query.Error.Error())
		return false, -1
	}

	if !count {
		return true, -1
	}

	var result int64
	if query.Count(&result); result == 0 {
		return true, 0
	}

	return true, result
}
