package helper

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

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

// chack probability for coupons (chance 0.0 to 1.0 : low to hign)
func CheckProbability(channce float64) bool {

	rand.Seed(time.Now().UnixMilli())

	// rangom.Float64() gives value of 0 to 1
	return channce > rand.Float64()
}

// random coupons
func CreateRandomCouponCode(couponCodeLenth int) string {
	// letter for coupns
	letters := `ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890`
	rand.Seed(time.Now().UnixMilli())

	// creat a byte array of couponCodeLength
	couponCode := make([]byte, couponCodeLenth)

	// loop through the array and randomly pic letter and add to array
	for i := range couponCode {
		couponCode[i] = letters[rand.Intn(len(letters))]
	}
	// convert into string and return the random letter array
	return string(couponCode)
}

// select a rangom number from start to end
func SelectRandomNumber(min, max int) int {

	rand.Seed(time.Now().UnixMilli())

	return rand.Intn(max-min) + min
}
