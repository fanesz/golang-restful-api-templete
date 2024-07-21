package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func corsConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []string{"Content-Type", "refresh_token", "access_token", "uuid"},
			ExposeHeaders:    []string{"Content-Length", "refresh_token", "access_token", "uuid"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		})
	}
}
