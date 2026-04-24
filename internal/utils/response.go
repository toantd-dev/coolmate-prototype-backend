package utils

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type PaginationMeta struct {
	Total       int64 `json:"total"`
	Page        int   `json:"page"`
	PageSize    int   `json:"pageSize"`
	TotalPages  int64 `json:"totalPages"`
}

type PaginatedResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    interface{}     `json:"data"`
	Meta    PaginationMeta  `json:"meta"`
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessPaginatedResponse(c *gin.Context, statusCode int, message string, data interface{}, meta PaginationMeta) {
	c.JSON(statusCode, PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string, errors interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

func BadRequest(c *gin.Context, message string) {
	ErrorResponse(c, 400, message, nil)
}

func Unauthorized(c *gin.Context, message string) {
	ErrorResponse(c, 401, message, nil)
}

func Forbidden(c *gin.Context, message string) {
	ErrorResponse(c, 403, message, nil)
}

func NotFound(c *gin.Context, message string) {
	ErrorResponse(c, 404, message, nil)
}

func InternalServerError(c *gin.Context, message string) {
	ErrorResponse(c, 500, message, nil)
}
