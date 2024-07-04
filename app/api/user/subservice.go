package user

import (
	"backend/app/common/consts"
	"backend/app/pkg/handler"
	"backend/app/pkg/validator"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func IsUsernameExists(Db *gorm.DB, c *gin.Context, username string, isHandleReturn bool) (bool, bool) {
	queryUsernameCheck := Db.Model(User{}).Where("username = ?", username)
	val, count := validator.Query(queryUsernameCheck, c, true)
	if !val {
		return true, false
	}
	if count > 0 {
		if isHandleReturn {
			handler.Error(c, http.StatusBadRequest, "Username already exists")
		}
		return false, true
	}

	return false, false
}

func IsEmailExists(Db *gorm.DB, c *gin.Context, email string, isHandleReturn bool) (bool, bool) {
	queryEmailCheck := Db.Model(User{}).Where("email = ?", email)
	val, count := validator.Query(queryEmailCheck, c, true)
	if !val {
		return true, false
	}
	if count > 0 {
		if isHandleReturn {
			handler.Error(c, http.StatusBadRequest, "Email already exists")
		}
		return false, true
	}

	return false, false
}

func GetRoleID(db *gorm.DB, c *gin.Context, roleName consts.Role) (bool, string) {
	query := db.Model(Role{}).Where("name = ?", roleName)
	val, count := validator.Query(query, nil, true)
	if !val {
		return true, ""
	}
	if count == 0 {
		handler.Error(c, http.StatusBadRequest, "Role not initialized")
		return true, ""
	}
	var role Role
	query.First(&role)

	return false, role.ID
}
