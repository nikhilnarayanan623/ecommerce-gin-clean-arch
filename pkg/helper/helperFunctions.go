package helper

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// there is 3 options username email number // count ho many are empty if its three its not valid
var count int

func Reset() { // reset it after validation for next login
	count = 0
	fmt.Println("reset count")
}

func CustomLoginValidator(fl validator.FieldLevel) bool {

	value := fl.Field().String()
	//fmt.Println(fl.GetTag(), fl.FieldName())

	if value == "" && count < 2 { //this for avoid two empty fied from user input
		count++
		return true
	} else if value != "" && len(value) >= 3 {
		return true
	}

	return false

}

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
