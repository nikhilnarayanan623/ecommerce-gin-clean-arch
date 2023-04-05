package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
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

func CompareUsers(user, checkUser domain.User) (err error) {
	if checkUser.Email == user.Email {
		err = errors.Join(err, errors.New("user already exist with this email"))
	}
	if checkUser.UserName == user.UserName {
		err = errors.Join(err, errors.New("user already exist with this user name"))
	}
	if checkUser.Phone == user.Phone {
		err = errors.Join(err, errors.New("user already exist with this phone"))
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

func VeifyRazorPaySignature(orderId, paymentId, signature string) error {

	razorPaySecret := config.GetCofig().RazorPaySecret
	data := orderId + "|" + paymentId

	h := hmac.New(sha256.New, []byte(razorPaySecret))
	_, err := h.Write([]byte(data))
	if err != nil {
		return errors.New("faild to veify signature")
	}

	sha := hex.EncodeToString(h.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(sha), []byte(signature)) == 1 {
		return nil
	}
	return errors.New("razorpay signature not match")
}
