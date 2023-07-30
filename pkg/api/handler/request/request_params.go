package request

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Get values from request form as string
func GetFormValuesAsString(ctx *gin.Context, name string) (value string, err error) {

	value = ctx.Request.PostFormValue(name)
	if value == "" {
		return "", fmt.Errorf("failed to get %s from request body", name)
	}

	return value, nil
}

// Get values from request form as string
func GetFormValuesAsUint(ctx *gin.Context, name string) (uint, error) {

	value := ctx.Request.PostFormValue(name)
	uintVal, err := strconv.ParseUint(value, 10, 64)

	if err != nil || uintVal == 0 {
		return 0, fmt.Errorf("failed to get %s from request body as int", name)
	}

	return uint(uintVal), nil
}

// Get query values as uint from request
func GetQueryValueAsUint(ctx *gin.Context, key string) (uint, error) {

	value := ctx.Query(key)
	uintVal, err := strconv.ParseUint(value, 10, 64)
	if err != nil || uintVal == 0 {
		return 0, fmt.Errorf("failed to get %s from query as int", key)
	}

	return uint(uintVal), nil
}

// Get query params as uint from request url
func GetParamAsUint(ctx *gin.Context, key string) (uint, error) {

	param := ctx.Param(key)
	value, err := strconv.ParseUint(param, 10, 64)

	if err != nil || value == 0 {
		return 0, fmt.Errorf("failed to get %s from param as int", key)
	}

	return uint(value), nil
}
