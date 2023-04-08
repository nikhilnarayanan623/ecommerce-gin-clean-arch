package res

import "strings"

type Response struct {
	StatusCode int         `json:"status"`
	Message    string      `json:"message"`
	Errors     interface{} `json:"errors,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func SuccessResponse(statusCode int, message string, data ...interface{}) Response {

	return Response{
		StatusCode: statusCode,
		Message:    message,
		Errors:     nil,
		Data:       data,
	}
}

func ErrorResponse(statusCode int, message string, err string, data interface{}) Response {

	spiltedError := strings.Split(err, "\n")
	return Response{
		StatusCode: statusCode,
		Message:    message,
		Errors:     spiltedError,
		Data:       data,
	}
}
