package utils

import "golang.org/x/crypto/bcrypt"

func GenerateHashFromPassword(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hash), err
}

func VerifyHashAndPassword(hashedPassword, password string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}
