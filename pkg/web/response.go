package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Response(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

func Success(c *gin.Context, status int, data interface{}) {
	Response(c, status, SuccessResponse{Data: data})
}

// NewErrorf creates a new error with the given status code and the message
// formatted according to args and format.
func Error(c *gin.Context, status int, format string, args ...interface{}) {
	err := ErrorResponse{
		Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(status)), " ", "_"),
		Message: fmt.Sprintf(format, args...),
		Status:  status,
	}

	Response(c, status, err)
}

func RespondWithError(c *gin.Context, status int, format string, args ...interface{}) {
	err := ErrorResponse{
		Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(status)), " ", "_"),
		Message: fmt.Sprintf(format, args...),
		Status:  status,
	}
	c.AbortWithStatusJSON(status, err)
}
