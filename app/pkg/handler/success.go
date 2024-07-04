package handler

import (
	"github.com/gin-gonic/gin"
)

type SuccessResponseAPI struct {
	Status     bool        `json:"status"`
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Payload    interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Page       int `json:"page" form:"page"`
	Limit      int `json:"limit" form:"limit"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

func Success(c *gin.Context, status int, message string, data interface{}, pagination ...Pagination) {
	response := SuccessResponseAPI{
		Status:     true,
		StatusCode: status,
		Message:    message,
		Payload:    data,
	}
	if len(pagination) > 0 && pagination[0].Limit != 0 && pagination[0].Page != 0 {
		response.Pagination = &pagination[0]
	}

	c.JSON(status, response)
}
