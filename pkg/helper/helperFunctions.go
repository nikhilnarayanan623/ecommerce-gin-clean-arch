package helper

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
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

func CompareUsers(user, checkUser domain.User) error {
	var err error
	if checkUser.Email == user.Email {
		err = errors.Join(err, errors.New("user already exist with this email"))
	}
	if checkUser.UserName == user.UserName {
		errors.Join(err, errors.New("user already exist with this user name"))
	}
	if checkUser.Phone == user.Phone {
		errors.Join(err, errors.New("user already exist with this phone"))
	}

	return err
}
