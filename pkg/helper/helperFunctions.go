package helper

import (
	"fmt"

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

	if value == "" || count < 2 {
		count++
		return true
	}

	return false
}
