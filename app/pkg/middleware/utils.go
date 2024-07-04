package middleware

import (
	"backend/app/common/utils"
	"backend/app/pkg/handler"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func GenerateToken(c *gin.Context, UUID uuid.UUID, expired TokenExpired) (string, bool) {
	claims := jwt.MapClaims{
		"uuid": UUID,
		"exp":  time.Now().Add(time.Duration(expired)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(utils.GetEnv("JWT_SECRET_KEY")))
	if err != nil {
		handler.Error(c, http.StatusInternalServerError, err.Error())
		return "", true
	}
	return tokenString, false
}

func GenerateResetPWToken(c *gin.Context, resetToken string, email string) (string, bool) {
	claims := jwt.MapClaims{
		"reset_token": resetToken,
		"email":       email,
		"exp":         time.Now().Add(time.Duration(RESETPW_TOKEN_EXPIRED)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(utils.GetEnv("JWT_SECRET_KEY")))
	if err != nil {
		handler.Error(c, http.StatusInternalServerError, err.Error())
		return "", true
	}
	return tokenString, false
}

func validateToken(tokenInput *string) (*PayloadToken, error) {
	res, err := jwt.Parse(*tokenInput, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(utils.GetEnv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := res.Claims.(jwt.MapClaims)
	if !ok || !res.Valid {
		return nil, errors.New("Unauthorized")
	}

	var payloadToken = PayloadToken{}
	payloadByte, _ := json.Marshal(claims)
	err = json.Unmarshal(payloadByte, &payloadToken)
	if err != nil {
		return nil, err
	}

	return &payloadToken, nil
}

func isTokenExpired(err error) bool {
	if err != nil {
		if err.Error() == "Token is expired" {
			return true
		}
	}
	return false
}
