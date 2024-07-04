package user

import (
	"backend/app/common/consts"
	"backend/app/common/utils"
	"backend/app/pkg/handler"
	"backend/app/pkg/mailer"
	"backend/app/pkg/mailer/templete"
	"backend/app/pkg/middleware"
	"backend/app/pkg/validator"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUser(db *gorm.DB, c *gin.Context, filters *UserFilters) {
	var users []User
	query := db.Model(users).Preload("Role")
	if val, _ := validator.Query(query, c, false); !val {
		return
	}

	filterService(c, query, filters)
	query.Find(&users)
	result := make([]*UserResponse, len(users))
	for i, user := range users {
		result[i] = responseFormatter(&user)
	}

	var response interface{} = result
	handler.Success(c, http.StatusOK, "Success getting users", &response, filters.Pagination)
}

func GetUserByUUID(db *gorm.DB, c *gin.Context, userUUID string) {
	var user User
	query := db.Model(&user).Where("uuid = ?", userUUID).Preload("Role")
	val, count := validator.Query(query, c, true)
	if !val {
		return
	}
	if count == 0 {
		handler.Error(c, http.StatusBadRequest, "User not found")
		return
	}

	query.First(&user)
	handler.Success(c, http.StatusOK, "Success getting user", responseFormatter(&user))
}

func CreateUser(db *gorm.DB, c *gin.Context, reqBody UserCreate) {
	if err, exist := IsEmailExists(db, c, reqBody.Email, true); err || exist {
		return
	}
	if err, exist := IsUsernameExists(db, c, reqBody.Username, true); err || exist {
		return
	}

	user := User{
		UUID:     utils.GenerateUUID(),
		Username: reqBody.Username,
		Email:    reqBody.Email,
		Password: reqBody.Password,
	}
	utils.Encrypt(&user.Password)

	if err, userRoleID := GetRoleID(db, c, consts.USER); err {
		return
	} else {
		user.RoleID = userRoleID
	}

	queryCreate := db.Create(&user)
	if val, _ := validator.Query(queryCreate, c, false); !val {
		return
	}

	queryCreatedUser := db.Preload("Role").First(&user, "uuid = ?", user.UUID)
	if val, _ := validator.Query(queryCreatedUser, c, false); !val {
		return
	}

	handler.Success(c, http.StatusCreated, "Success creating user", responseFormatter(&user))
}

func UpdateUser(db *gorm.DB, c *gin.Context, uuid string, reqBody UserUpdate) {
	var user User
	query := db.Model(&user).Where("uuid = ?", uuid)
	val, count := validator.Query(query, c, true)
	if !val {
		return
	}
	if count == 0 {
		handler.Error(c, http.StatusBadRequest, "User not found")
		return
	}

	query.First(&user)
	if user.Email != reqBody.Email {
		if err, exists := IsEmailExists(db, c, reqBody.Email, true); err || exists {
			return
		}
	}
	if user.Username != reqBody.Username {
		if err, exist := IsUsernameExists(db, c, reqBody.Username, true); err || exist {
			return
		}
	}

	user.Email = reqBody.Email
	user.Username = reqBody.Username

	if reqBody.Password != "" && reqBody.OldPassword != "" {
		if !utils.CompareEncrypted(&user.Password, &reqBody.OldPassword) {
			handler.Error(c, http.StatusBadRequest, "Invalid old password")
			return
		}
		if !utils.CompareEncrypted(&user.Password, &reqBody.Password) {
			utils.Encrypt(&reqBody.Password)
			user.Password = reqBody.Password
			mailer.SendMail(mailer.MailInfo{
				EmailTarget: []string{user.Email},
				Subject:     "Password Changed",
				Body:        templete.PasswordChanged,
			})
		}
	}

	if err, userRoleID := GetRoleID(db, c, reqBody.Role); err {
		return
	} else {
		user.RoleID = userRoleID
	}

	queryUpdate := db.Save(&user)
	if val, _ := validator.Query(queryUpdate, c, false); !val {
		return
	}

	queryUpdatedUser := db.Preload("Role").First(&user, "uuid = ?", user.UUID)
	if val, _ := validator.Query(queryUpdatedUser, c, false); !val {
		return
	}

	handler.Success(c, http.StatusOK, "Success updating user", responseFormatter(&user))
}

func SignIn(db *gorm.DB, c *gin.Context, reqBody UserLogin) {
	var user User
	query := db.Model(&user).Where("username = ? OR email = ?", reqBody.Username, reqBody.Email).Find(&user)
	val, count := validator.Query(query, c, true)
	if !val {
		return
	}
	if count == 0 {
		handler.Error(c, http.StatusBadRequest, "User not found")
		return
	}

	if !utils.CompareEncrypted(&user.Password, &reqBody.Password) {
		handler.Error(c, http.StatusBadRequest, "Invalid password")
		return
	}

	refreshToken, err := middleware.GenerateToken(c, user.UUID, middleware.REFRESH_TOKEN_EXPIRED)
	if err {
		return
	}
	accessToken, err := middleware.GenerateToken(c, user.UUID, middleware.ACCESS_TOKEN_EXPIRED)
	if err {
		return
	}

	user.LoginToken = refreshToken
	user.IPAddress = c.ClientIP()
	queryUpdate := db.Model(&user).Where("uuid = ?", user.UUID).Updates(&user)
	if val, _ := validator.Query(queryUpdate, c, false); !val {
		return
	}

	c.Header("uuid", user.UUID.String())
	c.Header("refresh_token", refreshToken)
	c.Header("access_token", accessToken)

	handler.Success(c, http.StatusOK, "Success signing in", nil)
}

func SignOut(db *gorm.DB, c *gin.Context) {
	var user User
	query := db.Model(&user).Where("uuid = ?", c.GetString("uuid"))
	val, count := validator.Query(query, c, true)
	if !val {
		return
	}
	if count == 0 {
		handler.Error(c, http.StatusBadRequest, "User not found or not logged in")
		return
	}

	query.First(&user)
	queryUpdate := db.Model(&user).Update("login_token", "")
	if val, _ := validator.Query(queryUpdate, c, false); !val {
		return
	}

	handler.Success(c, http.StatusOK, "Success signing out", nil)
}

func RequestResetPassword(db *gorm.DB, c *gin.Context, reqBody PasswordResetRequest) {
	var user User
	query := db.Model(&user).Where("email = ?", reqBody.Email)
	val, count := validator.Query(query, c, true)
	if !val {
		return
	}
	if count == 0 {
		handler.Error(c, http.StatusBadRequest, "User not found")
		return
	}

	query.First(&user)
	user.ResetPWToken = utils.GenerateRandomStringWithNumber(64)
	resetPWToken, err := middleware.GenerateResetPWToken(c, user.ResetPWToken, user.Email)
	if err {
		return
	}

	resetPasswordLink := utils.GetEnv("RESETPW_FE_ENDPOINT") + "/" + resetPWToken
	mailer.SendMail(mailer.MailInfo{
		EmailTarget: []string{user.Email},
		Subject:     "Password Reset Request",
		Body: fmt.Sprintf(
			templete.PasswordResetRequest,
			resetPasswordLink,
		),
	})

	go func() {
		db.Model(&user).Update("reset_pw_token", user.ResetPWToken)
	}()

	handler.Success(c, http.StatusOK, "Success sending reset password to email", nil)
}

func ResetPassword(db *gorm.DB, c *gin.Context, reqBody PasswordReset) {
	userEmailReq := c.GetString("email")
	resetTokenReq := c.GetString("reset_token")
	if userEmailReq == "" || resetTokenReq == "" {
		handler.Error(c, http.StatusBadRequest, "Email and reset token are required")
		return
	}

	var user User
	query := db.Model(&user).Where("email = ?", userEmailReq)
	val, count := validator.Query(query, c, true)
	if !val {
		return
	}
	if count == 0 {
		handler.Error(c, http.StatusBadRequest, "User not found")
		return
	}

	query.First(&user)
	if user.ResetPWToken != resetTokenReq {
		handler.Error(c, http.StatusBadRequest, "Invalid reset token")
		return
	}

	user.Password = reqBody.NewPassword
	utils.Encrypt(&user.Password)

	queryUpdate := db.Debug().Model(&user).Updates(map[string]interface{}{
		"Password":     user.Password,
		"ResetPWToken": "",
		"LoginToken":   "",
		"IPAddress":    "",
	})
	if val, _ := validator.Query(queryUpdate, c, false); !val {
		return
	}

	mailer.SendMail(mailer.MailInfo{
		EmailTarget: []string{user.Email},
		Subject:     "Password Changed",
		Body:        templete.PasswordResetted,
	})

	handler.Success(c, http.StatusOK, "Success updating user password", nil)
}
