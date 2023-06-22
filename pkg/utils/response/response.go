package response

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  bool        `json:"success"`
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(ctx *gin.Context, statusCode int, message string, data ...interface{}) {

	log.Printf("\033[0;32m%s\033[0m\n", message)

	response := Response{
		Status:  true,
		Message: message,
		Error:   nil,
		Data:    data,
	}
	ctx.JSON(statusCode, response)
}

func ErrorResponse(ctx *gin.Context, statusCode int, message string, err error, data interface{}) {

	log.Printf("\033[0;31m%s\033[0m\n", err.Error())

	errFields := strings.Split(err.Error(), "\n")
	response := Response{
		Status:  false,
		Message: message,
		Error:   errFields,
		Data:    data,
	}

	ctx.JSON(statusCode, response)
}
