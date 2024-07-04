package middleware

import (
	"backend/app/common/consts"
	"backend/app/pkg/handler"
	"backend/app/pkg/mailer"
	"backend/app/pkg/mailer/templete"
	"backend/app/pkg/validator"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type loggedUser struct {
	Email     string `gorm:"column:email"`
	IpAddress string `gorm:"column:ip_address"`
}

func isUserLogin(db *gorm.DB, c *gin.Context, token string, ipAddress string) (bool, bool) {
	var loggedUser loggedUser
	query := db.Table("users").Select("email, ip_address").Where("login_token = ?", token).Scan(&loggedUser)
	val, count := validator.Query(query, c, true)
	if !val {
		return false, true
	}
	if count == 0 {
		return false, false
	}
	if loggedUser.IpAddress != ipAddress {
		mailer.SendMail(mailer.MailInfo{
			EmailTarget: []string{loggedUser.Email},
			Subject:     "Security Alert!",
			Body: fmt.Sprintf(
				templete.TokenStolen,
				ipAddress,
			),
		})
		return false, false
	}

	return true, false
}

func IsAdmin(db *gorm.DB, c *gin.Context, userUUID *string) (bool, bool) {
	if *userUUID == "" {
		handler.Error(c, http.StatusBadRequest, "UUID is required")
		return true, false
	}
	var roleID string
	query := db.Table("users").Select("role_id").Where("uuid = ?", userUUID).Scan(&roleID)
	val, _ := validator.Query(query, c, false)
	if !val {
		return true, false
	}

	var role consts.Role
	queryRoleName := db.Table("roles").Select("name").Where("id = ?", roleID).Scan(&role)
	val, _ = validator.Query(queryRoleName, c, false)
	if !val {
		return true, false
	}

	if role == consts.ADMIN {
		return false, true
	}

	return false, false
}

func updateLoginToken(db *gorm.DB, c *gin.Context, oldToken *string, newToken *string, ipAddress *string) bool {
	query := db.Table("users").Where("login_token = ? AND ip_address = ?", oldToken, ipAddress).Update("login_token", newToken)
	val, _ := validator.Query(query, c, false)
	return !val
}
