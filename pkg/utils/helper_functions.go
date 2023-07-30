package utils

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// take userId from context
func GetUserIdFromContext(ctx *gin.Context) uint {
	userID := ctx.GetUint("userId")
	return userID
}

func StringToUint(str string) (uint, error) {
	val, err := strconv.Atoi(str)
	return uint(val), err
}

// generate userName
func GenerateRandomUserName(FirstName string) string {

	suffix := make([]byte, 4)

	numbers := "1234567890"
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))

	for i := range suffix {
		suffix[i] = numbers[rng.Intn(10)]
	}

	userName := (FirstName + string(suffix))

	return strings.ToLower(userName)
}

// generate unique string for sku
func GenerateSKU() string {
	sku := make([]byte, 10)

	rand.Read(sku)

	return hex.EncodeToString(sku)
}

// random coupons
func GenerateCouponCode(couponCodeLength int) string {
	// letter for coupons
	letters := `ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890`
	rand.Seed(time.Now().UnixMilli())

	// create a byte array of couponCodeLength
	couponCode := make([]byte, couponCodeLength)

	// loop through the array and randomly pic letter and add to array
	for i := range couponCode {
		couponCode[i] = letters[rand.Intn(len(letters))]
	}
	// convert into string and return the random letter array
	return string(couponCode)
}

func StringToTime(timeString string) (timeValue time.Time, err error) {

	// parse the string time to time
	timeValue, err = time.Parse(time.RFC3339Nano, timeString)

	if err != nil {
		return timeValue, fmt.Errorf("faild to parse given time %v to time variable \nivalid input", timeString)
	}
	return timeValue, err
}

func GenerateRandomString(length int) string {
	sku := make([]byte, length)

	rand.Read(sku)

	return hex.EncodeToString(sku)
}

func RandomInt(min, max int) int {
	rand.Seed(time.Hour.Nanoseconds())

	return rand.Intn(max-min) + min
}

func GetHashedPassword(password string) (hashedPassword string, err error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return hashedPassword, err
	}
	hashedPassword = string(hash)
	return hashedPassword, nil
}

func ComparePasswordWithHashedPassword(actualpassword, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(actualpassword))
	return err
}
