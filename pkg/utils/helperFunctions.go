package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/razorpay/razorpay-go"
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

// // chack probability for coupons (chance 0.0 to 1.0 : low to hign)
// func CheckProbability(channce float64) bool {

// 	rand.Seed(time.Now().UnixMilli())

// 	// rangom.Float64() gives value of 0 to 1
// 	return channce > rand.Float64()
// }

// random coupons
func GenerateCouponCode(couponCodeLenth int) string {
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

// // select a rangom number from start to end
// func SelectRandomNumber(min, max int) int {

// 	rand.Seed(time.Now().UnixMilli())

// 	return rand.Intn(max-min) + min
// }

// function for generate razorpay order
func GenerateRazorpayOrder(razorPayAmount uint, recieptIdOptional string) (razorpayOrderID interface{}, err error) {
	// get razor pay key and secret
	razorpayKey := config.GetCofig().RazorPayKey
	razorpaySecret := config.GetCofig().RazorPaySecret

	//create a razorpay client
	client := razorpay.NewClient(razorpayKey, razorpaySecret)

	data := map[string]interface{}{
		"amount":   razorPayAmount,
		"currency": "INR",
		"receipt":  recieptIdOptional,
	}
	// create an order on razor pay
	razorpayRes, err := client.Order.Create(data, nil)
	if err != nil {
		return razorpayOrderID, fmt.Errorf("fadil to create razorpay order for amount %v", razorPayAmount)
	}

	razorpayOrderID = razorpayRes["id"]

	return razorpayOrderID, nil
}

// //verify the razorpay signature
// err = utils.VeifyRazorPaySignature(razorpayOrderID, razorpayPaymentID, razorpaySignature)
// if err != nil {
// 	respones := res.ErrorResponse(400, "faild to veify razorpay payment", err.Error(), nil)
// 	ctx.JSON(http.StatusBadRequest, respones)
// 	return
// }

// // get razorpay key and secret
// razorpayKey, razorPaySecret := config.GetCofig().RazorPayKey, config.GetCofig().RazorPaySecret
// // create a new razorpay client
// razorpayClient := razorpay.NewClient(razorpayKey, razorPaySecret)
// payment, err := razorpayClient.Payment.Fetch(razorpayPaymentID, nil, nil)
// if err != nil {
// 	response := res.ErrorResponse(400, "faild to get payment details", err.Error(), nil)
// 	ctx.JSON(400, response)
// 	return
// }

//	if payment["status"] != "captured" {
//		response := res.ErrorResponse(400, "payment faild", "payment not got on razorpay", payment)
//		ctx.JSON(400, response)
//		return
//	}
func VeifyRazorpayPayment(razorpayOrderID, razorpayPaymentID, razorpaySignatur string) error {

	razorpayKey := config.GetCofig().RazorPayKey
	razorPaySecret := config.GetCofig().RazorPaySecret

	//varify signature
	data := razorpayOrderID + "|" + razorpayPaymentID
	h := hmac.New(sha256.New, []byte(razorPaySecret))
	_, err := h.Write([]byte(data))
	if err != nil {
		return errors.New("faild to veify signature")
	}

	sha := hex.EncodeToString(h.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(sha), []byte(razorpaySignatur)) != 1 {
		return errors.New("razorpay signature not match")
	}

	// then vefiy payment
	client := razorpay.NewClient(razorpayKey, razorPaySecret)

	// fetch payment and vefify
	payment, err := client.Payment.Fetch(razorpayPaymentID, nil, nil)

	if err != nil {
		return err
	}

	// check payment status
	if payment["status"] != "captured" {
		return fmt.Errorf("faild to varify payment \nrazorpay payment with payment_id %v", razorpayPaymentID)
	}

	return nil
}
