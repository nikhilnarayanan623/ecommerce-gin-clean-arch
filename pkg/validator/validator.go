package validator

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

const (
	ValidationWhiteSpace = "whitespace"
)

// To register all custom validations
func RegisterAllCustomValidations() error {

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("failed to create validator")
	}

	v.RegisterValidation(ValidationWhiteSpace, ValidateWhitespace)

	return nil
}

// validate white space
func ValidateWhitespace(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	regex := regexp.MustCompile(`^\s*$`)
	return !regex.MatchString(str)
}
