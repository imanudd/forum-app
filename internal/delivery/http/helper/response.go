package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSONResponse struct {
	Code       int         `json:"code,omitempty"`
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, code int, data ...interface{}) {
	hte := JSONResponse{
		StatusCode: code,
		Message:    http.StatusText(code),
	}

	if len(data) > 0 {
		hte.Data = data[0]
	}

	c.JSON(code, hte)
}

func Error(c *gin.Context, code int, message string) {
	hte := JSONResponse{
		Code:       code,
		StatusCode: code,
		Message:    message,
	}

	if code == http.StatusUnauthorized {
		defer c.Abort()
	}

	c.JSON(code, hte)
}

func InternalError(c *gin.Context, err error) {
	response := JSONResponse{
		Code:       http.StatusInternalServerError,
		StatusCode: http.StatusInternalServerError,
		Message:    err.Error(),
	}

	c.JSON(response.Code, response)
}
