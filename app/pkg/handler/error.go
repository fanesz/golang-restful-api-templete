package handler

import (
	"sync"

	"github.com/gin-gonic/gin"
)

type ApiError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponseAPI struct {
	Status     bool       `json:"status"`
	StatusCode int        `json:"status_code"`
	Message    string     `json:"message"`
	Errors     []ApiError `json:"errors"`
}

var errorResponsePool = sync.Pool{
	New: func() interface{} {
		return &ErrorResponseAPI{}
	},
}

func Error(c *gin.Context, status int, message string, errors ...ApiError) {
	res := errorResponsePool.Get().(*ErrorResponseAPI)
	res.Status = false
	res.StatusCode = status
	res.Message = message
	res.Errors = errors
	c.AbortWithStatusJSON(status, res)
	errorResponsePool.Put(res)
}
