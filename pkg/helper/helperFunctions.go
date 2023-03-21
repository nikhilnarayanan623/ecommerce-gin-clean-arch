package helper

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// take userId from context
func GetUserIdFromContext(ctx *gin.Context) uint {
	userIdStr := ctx.GetString("userId")
	userIdInt, _ := strconv.Atoi(userIdStr)
	return uint(userIdInt)
}

func StringToUint(str string) (uint, error) {
	val, err := strconv.Atoi(str)
	return uint(val), err
}
