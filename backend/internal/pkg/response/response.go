package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
}

type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Body{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

func SuccessPage(c *gin.Context, data interface{}, pagination Pagination) {
	c.JSON(http.StatusOK, Body{
		Code:       200,
		Message:    "success",
		Data:       data,
		Pagination: pagination,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Body{
		Code:    201,
		Message: "created",
		Data:    data,
	})
}

func Error(c *gin.Context, httpStatus, code int, message string) {
	c.JSON(httpStatus, Body{
		Code:    code,
		Message: message,
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, 400, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, 401, message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, 403, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, 404, message)
}

func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, 500, message)
}
