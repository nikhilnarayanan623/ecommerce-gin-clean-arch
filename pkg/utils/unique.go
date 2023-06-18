package utils

import "github.com/google/uuid"

func GenerateUniqueString() string {

	return uuid.NewString()
}
