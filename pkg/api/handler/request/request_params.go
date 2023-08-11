package request

import (
	"fmt"
	"mime/multipart"
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
	uintVal, err := strconv.ParseUint(value, 10, 32)

	if err != nil || uintVal == 0 {
		return 0, fmt.Errorf("failed to get %s from request body as int", name)
	}

	return uint(uintVal), nil
}

// Get query values as uint from request
func GetQueryValueAsUint(ctx *gin.Context, key string) (uint, error) {

	value := ctx.Query(key)
	uintVal, err := strconv.ParseUint(value, 10, 32)
	if err != nil || uintVal == 0 {
		return 0, fmt.Errorf("failed to get %s from query as int", key)
	}

	return uint(uintVal), nil
}

// Get query params as uint from request url
func GetParamAsUint(ctx *gin.Context, key string) (uint, error) {

	param := ctx.Param(key)
	value, err := strconv.ParseUint(param, 10, 32)

	if err != nil || value == 0 {
		return 0, fmt.Errorf("failed to get %s from param as int", key)
	}

	return uint(value), nil
}

// Get multiple form files from request body
func GetArrayOfFromFiles(ctx *gin.Context, name string) ([]*multipart.FileHeader, error) {

	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil {

		return nil, fmt.Errorf("failed to parse form data err: %w", err)
	}

	files, ok := ctx.Request.MultipartForm.File[name]

	if !ok || len(files) == 0 {
		return nil, fmt.Errorf("failed to get %s files from request form data", name)
	}

	return files, nil
}

// Get value from request form as array slice
func GetArrayFormValueAsUint(ctx *gin.Context, name string) ([]uint, error) {

	values, ok := ctx.GetPostFormArray(name)
	if !ok || len(values) == 0 {
		return nil, fmt.Errorf("failed to get %s of array from request form data", name)
	}

	uintValues := make([]uint, len(values))

	for i := range values {
		fmt.Println("value: ", i)

		num, err := strconv.ParseUint(values[i], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("request value is not and integer for %s values", name)
		}
		uintValues[i] = uint(num)
	}

	return uintValues, nil
}
