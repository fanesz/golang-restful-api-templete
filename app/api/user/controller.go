package user

import (
	"backend/app/common/consts"
	"backend/app/common/utils"
	"backend/app/pkg/handler"
	"backend/app/pkg/middleware"
	"backend/app/pkg/validator"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Controller(router *gin.Engine, db *gorm.DB) {
	v1 := router.Group("/v1")

	v1.GET("/users", middleware.ValidateHeader(db, true), func(c *gin.Context) {
		var filters UserFilters
		if err := validator.BindParams(c, &filters); err {
			return
		}

		GetUser(db, c, &filters)
	})

	v1.GET("/users/:uuid", middleware.ValidateHeader(db, true), func(c *gin.Context) {
		var userUUID UserGetByUUID
		if err := validator.BindUri(c, &userUUID); err {
			return
		}

		if !utils.IsValidUUID(userUUID.UUID) {
			handler.Error(c, http.StatusBadRequest, "Invalid UUID")
			return
		}

		GetUserByUUID(db, c, userUUID.UUID)
	})

	v1.POST("/signup", func(c *gin.Context) {
		var reqBody UserCreate
		if err := validator.BindBody(c, &reqBody); err {
			return
		}

		CreateUser(db, c, reqBody)
	})

	v1.PUT("/users/:uuid", middleware.ValidateHeader(db, true), func(c *gin.Context) {
		var reqBody UserUpdate
		if err := validator.BindBody(c, &reqBody); err {
			return
		}
		if err := validator.BindUri(c, &reqBody); err {
			return
		}
		userUUID := c.GetString("uuid")
		reqBody.Role = consts.Role(strings.ToUpper(string(reqBody.Role)))

		if reqBody.UUID != userUUID || reqBody.Role != consts.USER {
			err, res := middleware.IsAdmin(db, c, &userUUID)
			if err {
				return
			}
			if !res {
				handler.Error(c, http.StatusUnauthorized, "Required admin role")
				return
			}
		}

		UpdateUser(db, c, reqBody.UUID, reqBody)
	})

	v1.POST("/signin", func(c *gin.Context) {
		var reqBody UserLogin
		if err := validator.BindBody(c, &reqBody); err {
			return
		}

		if reqBody.Username == "" && reqBody.Email == "" {
			handler.Error(c, http.StatusBadRequest, "Username or Email is required")
			return
		}

		SignIn(db, c, reqBody)
	})

	v1.POST("/signout", middleware.ValidateHeader(db, false), func(c *gin.Context) {
		SignOut(db, c)
	})

	v1.POST("/check-email", func(c *gin.Context) {
		var reqBody CheckEmail
		if err := validator.BindBody(c, &reqBody); err {
			return
		}

		err, exist := IsEmailExists(db, c, reqBody.Email, false)
		if err {
			return
		}

		response := gin.H{"isExist": exist}
		handler.Success(c, http.StatusOK, "Success checking email", &response)
	})

	v1.POST("/check-username", func(c *gin.Context) {
		var reqBody CheckUsername
		if err := validator.BindBody(c, &reqBody); err {
			return
		}

		err, exist := IsUsernameExists(db, c, reqBody.Username, false)
		if err {
			return
		}

		response := gin.H{"isExist": exist}
		handler.Success(c, http.StatusOK, "Success checking username", &response)
	})

	v1.POST("/reset-password", func(c *gin.Context) {
		var reqBody PasswordResetRequest
		if err := validator.BindBody(c, &reqBody); err {
			return
		}

		RequestResetPassword(db, c, reqBody)
	})

	v1.GET("/check-reset-token/:token", middleware.ValidateResetPassword(), func(c *gin.Context) {
		handler.Success(c, http.StatusOK, "Reset token is valid", nil)
	})

	v1.PUT("/reset-password", middleware.ValidateResetPassword(), func(c *gin.Context) {
		var reqBody PasswordReset
		if err := validator.BindBody(c, &reqBody); err {
			return
		}

		ResetPassword(db, c, reqBody)
	})

}

func responseFormatter(user *User) *UserResponse {
	return &UserResponse{
		UUID:     user.UUID.String(),
		Username: user.Username,
		Email:    user.Email,
		Role:     string(user.Role.Name),
	}
}
