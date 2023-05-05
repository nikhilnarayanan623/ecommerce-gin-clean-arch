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
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/razorpay/razorpay-go"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"golang.org/x/crypto/bcrypt"
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

// generate userName
func GenerateRandomUserName(FirstName string) string {

	suffix := make([]byte, 4)

	numbers := "1234567890"
	rand.Seed(time.Now().UnixMilli())

	for i := range suffix {
		suffix[i] = numbers[rand.Intn(10)]
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

func StringToTime(timeString string) (timeValue time.Time, err error) {

	// parse the string time to time
	timeValue, err = time.Parse(time.RFC3339Nano, timeString)

	if err != nil {
		return timeValue, fmt.Errorf("faild to parse given time %v to time variable \nivalid input", timeString)
	}
	return timeValue, err
}

func GenerateStipeClientSecret(amountToPay uint, recieptEmail string) (clientSecret string, err error) {
	// set up the stip secret key
	stripe.Key = config.GetCofig().StripSecretKey

	// create a payment param
	params := &stripe.PaymentIntentParams{

		Amount:       stripe.Int64(int64(amountToPay)),
		ReceiptEmail: stripe.String(recieptEmail),

		Currency: stripe.String(string(stripe.CurrencyINR)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	// creata new payment intent with this param
	paymentIntent, err := paymentintent.New(params)

	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("faild to create strip payment for amount %v", amountToPay)
	}

	clientSecret = paymentIntent.ClientSecret
	return clientSecret, nil
}

func VeifyStripePaymentIntentByID(paymentID string) error {

	stripe.Key = config.GetCofig().StripSecretKey

	// get payment by payment_id
	paymentIntent, err := paymentintent.Get(paymentID, nil)

	if err != nil {
		return fmt.Errorf("faild to get stripe paymentIntent of payment_id %v", paymentID)
	}

	// verify the payment intent
	if paymentIntent.Status != stripe.PaymentIntentStatusSucceeded && paymentIntent.Status != stripe.PaymentIntentStatusRequiresCapture {
		return fmt.Errorf("payment not not completed")
	}

	return nil
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
