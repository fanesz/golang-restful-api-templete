package middleware

import (
	"backend/app/common/utils"
	"backend/app/pkg/handler"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PayloadToken struct {
	UUID       uuid.UUID `json:"uuid"`
	Email      string    `json:"email"`
	ResetToken string    `json:"reset_token"`
}

type TokenExpired = int64

var (
	ACCESS_TOKEN_EXPIRED  TokenExpired = int64(time.Minute * 1)
	REFRESH_TOKEN_EXPIRED TokenExpired = int64(time.Hour * 24)
	RESETPW_TOKEN_EXPIRED TokenExpired = int64(time.Minute * 20)
)

func ValidateHeader(db *gorm.DB, isAutoRefresh bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		headerUUID := c.GetHeader("uuid")
		userUUID := utils.ParseUUID(&headerUUID)
		userIP := c.ClientIP()

		if utils.IsUUIDNil(userUUID) {
			handler.Error(c, http.StatusUnauthorized, "UUID is required")
			return
		}

		accessToken := c.GetHeader("access_token")
		if accessToken == "" || accessToken == "null" {
			handler.Error(c, http.StatusUnauthorized, "Access Token is required")
			return
		}

		accessTokenPayload, err := validateToken(&accessToken)
		if err != nil && !isTokenExpired(err) {
			handler.Error(c, http.StatusUnauthorized, "Invalid Token")
			return
		}
		if err != nil && isTokenExpired(err) {
			refreshToken := c.GetHeader("refresh_token")
			if refreshToken == "" || refreshToken == "null" {
				handler.Error(c, http.StatusUnauthorized, "Refresh Token is required")
				return
			}

			refreshTokenPayload, err := validateToken(&refreshToken)
			if err != nil && !isTokenExpired(err) {
				handler.Error(c, http.StatusUnauthorized, "Invalid Token")
				return
			}
			if err != nil && isTokenExpired(err) {
				res, err := isUserLogin(db, c, refreshToken, userIP)
				if err {
					return
				}
				if !res {
					handler.Error(c, http.StatusUnauthorized, "Unauthorized")
					return
				}

				if isAutoRefresh {
					newRefreshToken, errs := GenerateToken(c, userUUID, REFRESH_TOKEN_EXPIRED)
					if errs {
						return
					}
					if updateLoginToken(db, c, &refreshToken, &newRefreshToken, &userIP) {
						return
					}

					c.Header("refresh_token", newRefreshToken)
				} else {
					handler.Error(c, http.StatusUnauthorized, "Token is expired")
					return
				}
			}
			if err == nil && refreshTokenPayload.UUID != userUUID {
				handler.Error(c, http.StatusUnauthorized, "Invalid Token")
				return
			}
			if err == nil && refreshTokenPayload.UUID == userUUID {
				res, err := isUserLogin(db, c, refreshToken, userIP)
				if err {
					return
				}
				if !res {
					handler.Error(c, http.StatusUnauthorized, "Unauthorized")
					return
				}
				c.Set("uuid", refreshTokenPayload.UUID.String())
			}

			newAccessToken, errs := GenerateToken(c, userUUID, ACCESS_TOKEN_EXPIRED)
			if errs {
				return
			}
			if isAutoRefresh {
				c.Header("access_token", newAccessToken)
			}
		}
		if err == nil && accessTokenPayload.UUID != userUUID {
			handler.Error(c, http.StatusUnauthorized, "Invalid Token")
			return
		}
		if err == nil && accessTokenPayload.UUID == userUUID {
			c.Set("uuid", accessTokenPayload.UUID.String())
		}

		c.Next()
	}
}

func ValidateResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("resetpw_token")
		if token == "" || token == "null" {
			handler.Error(c, http.StatusUnauthorized, "Reset Password is required")
			return
		}

		payload, err := validateToken(&token)
		if err != nil && !isTokenExpired(err) {
			handler.Error(c, http.StatusUnauthorized, "Invalid Token")
			return
		}
		if err != nil && isTokenExpired(err) {
			handler.Error(c, http.StatusUnauthorized, "Token is expired")
			return
		}

		c.Set("email", payload.Email)
		c.Set("reset_token", payload.ResetToken)
		c.Next()
	}
}
