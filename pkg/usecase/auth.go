package usecase

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/service/otp"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/service/token"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
)

const (
	countryCode       = "+91"
	otpExpireDuration = time.Minute * 2
)

type authUseCase struct {
	authRepo interfaces.AuthRepository

	userRepo     interfaces.UserRepository
	adminRepo    interfaces.AdminRepository
	tokenService token.TokenService
	optAuth      otp.OtpAuth
}

func NewAuthUseCase(authRepo interfaces.AuthRepository, tokenService token.TokenService,
	userRepo interfaces.UserRepository, adminRepo interfaces.AdminRepository,
	optAuth otp.OtpAuth) service.AuthUseCase {

	return &authUseCase{
		userRepo:     userRepo,
		adminRepo:    adminRepo,
		tokenService: tokenService,
		authRepo:     authRepo,
		optAuth:      optAuth,
	}
}

const (
	AccessTokenDuration  = time.Minute * 20
	RefreshTokenDuration = time.Hour * 24 * 7
)

func (c *authUseCase) UserLogin(ctx context.Context, loginDetails request.Login) (uint, error) {

	var (
		user domain.User
		err  error
	)
	switch {
	case loginDetails.Email != "":
		user, err = c.userRepo.FindUserByEmail(ctx, loginDetails.Email)
	case loginDetails.UserName != "":
		user, err = c.userRepo.FindUserByUserName(ctx, loginDetails.UserName)
	case loginDetails.Phone != "":
		user, err = c.userRepo.FindUserByPhoneNumber(ctx, loginDetails.Phone)
	default:
		return 0, ErrEmptyLoginCredentials
	}

	if err != nil {
		return 0, utils.PrependMessageToError(err, "failed to find user from database")
	}

	if user.ID == 0 {
		return 0, ErrUserNotExist
	}

	if !user.Verified {
		return 0, ErrUserNotVerified
	}

	if user.BlockStatus {
		return 0, ErrUserBlocked
	}

	err = utils.ComparePasswordWithHashedPassword(loginDetails.Password, user.Password)
	if err != nil {
		return 0, ErrWrongPassword
	}

	return user.ID, nil
}

func (c *authUseCase) UserLoginOtpSend(ctx context.Context, loginDetails request.OTPLogin) (string, error) {

	var (
		user domain.User
		err  error
	)

	switch {

	case loginDetails.Email != "":
		user, err = c.userRepo.FindUserByEmail(ctx, loginDetails.Email)
	case loginDetails.UserName != "":
		user, err = c.userRepo.FindUserByUserName(ctx, loginDetails.UserName)
	case loginDetails.Phone != "":
		user, err = c.userRepo.FindUserByPhoneNumber(ctx, loginDetails.Phone)
	default:
		return "", ErrEmptyLoginCredentials
	}

	if err != nil {
		return "", fmt.Errorf("can't find the user \nerror:%v", err.Error())
	}

	if user.ID == 0 {
		return "", ErrUserNotExist
	}

	if user.BlockStatus {
		return "", ErrUserBlocked
	}

	errChan := make(chan error, 2)
	wait := sync.WaitGroup{}
	wait.Add(2)

	go func() {
		defer wait.Done()
		_, err := c.optAuth.SentOtp(countryCode + user.Phone)
		if err != nil {
			errChan <- fmt.Errorf("failed to send otp \nerrors:%v", err.Error())
		}
	}()
	otpID := uuid.NewString()

	go func() {
		defer wait.Done()
		otpSession := domain.OtpSession{
			OtpID:    otpID,
			UserID:   user.ID,
			Phone:    user.Phone,
			ExpireAt: time.Now().Add(otpExpireDuration), // 2 minutes expire for otp
		}
		err := c.authRepo.SaveOtpSession(ctx, otpSession)
		if err != nil {
			errChan <- fmt.Errorf("failed to save otp session \nerror:%v", err.Error())
		}
	}()

	wait.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return "", err
		}
	}

	return otpID, nil
}

func (c *authUseCase) LoginOtpVerify(ctx context.Context, otpVerifyDetails request.OTPVerify) (uint, error) {

	otpSession, err := c.authRepo.FindOtpSession(ctx, otpVerifyDetails.OtpID)
	if err != nil {
		return 0, utils.PrependMessageToError(err, "failed to find otp session from database")
	}

	if time.Since(otpSession.ExpireAt) > 0 {
		return 0, ErrOtpExpired
	}

	valid, err := c.optAuth.VerifyOtp(countryCode+otpSession.Phone, otpVerifyDetails.Otp)
	if err != nil {
		return 0, utils.PrependMessageToError(err, "failed to verify otp")
	}
	if !valid {
		return 0, ErrInvalidOtp
	}

	return otpSession.UserID, nil
}

func (c *authUseCase) AdminLogin(ctx context.Context, loginDetails request.Login) (uint, error) {

	var (
		admin domain.Admin
		err   error
	)
	switch {
	case loginDetails.Email != "":
		admin, err = c.adminRepo.FindAdminByEmail(ctx, loginDetails.Email)
	case loginDetails.UserName != "":
		admin, err = c.adminRepo.FindAdminByUserName(ctx, loginDetails.UserName)
	default:
		return 0, ErrEmptyLoginCredentials
	}

	if err != nil {
		return 0, utils.PrependMessageToError(err, "failed to find admin")
	}

	if admin.ID == 0 {
		return 0, ErrUserNotExist
	}

	err = utils.ComparePasswordWithHashedPassword(loginDetails.Password, admin.Password)
	if err != nil {
		return 0, ErrWrongPassword
	}

	return admin.ID, nil
}

func (c *authUseCase) GenerateAccessToken(ctx context.Context, tokenParams service.GenerateTokenParams) (string, error) {

	tokenReq := token.GenerateTokenRequest{
		UserID:   tokenParams.UserID,
		UsedFor:  tokenParams.UserType,
		ExpireAt: time.Now().Add(AccessTokenDuration),
	}

	tokenRes, err := c.tokenService.GenerateToken(tokenReq)

	return tokenRes.TokenString, err
}
func (c *authUseCase) GenerateRefreshToken(ctx context.Context, tokenParams service.GenerateTokenParams) (string, error) {

	expireAt := time.Now().Add(RefreshTokenDuration)
	tokenReq := token.GenerateTokenRequest{
		UserID:   tokenParams.UserID,
		UsedFor:  tokenParams.UserType,
		ExpireAt: expireAt,
	}
	tokenRes, err := c.tokenService.GenerateToken(tokenReq)
	if err != nil {
		return "", err
	}

	err = c.authRepo.SaveRefreshSession(ctx, domain.RefreshSession{
		UserID:       tokenParams.UserID,
		TokenID:      tokenRes.TokenID,
		RefreshToken: tokenRes.TokenString,
		ExpireAt:     expireAt,
	})
	if err != nil {
		return "", err
	}
	log.Printf("successfully refresh token created and refresh session stored in database")
	return tokenRes.TokenString, nil
}

func (c *authUseCase) VerifyAndGetRefreshTokenSession(ctx context.Context, refreshToken string, usedFor token.UserType) (domain.RefreshSession, error) {

	verifyReq := token.VerifyTokenRequest{
		TokenString: refreshToken,
		UsedFor:     usedFor,
	}
	verifyRes, err := c.tokenService.VerifyToken(verifyReq)
	if err != nil {
		return domain.RefreshSession{}, utils.PrependMessageToError(ErrInvalidRefreshToken, err.Error())
	}

	refreshSession, err := c.authRepo.FindRefreshSessionByTokenID(ctx, verifyRes.TokenID)
	if err != nil {
		return refreshSession, err
	}

	if refreshSession.TokenID == "" {
		return refreshSession, ErrRefreshSessionNotExist
	}

	if time.Since(refreshSession.ExpireAt) > 0 {
		return domain.RefreshSession{}, ErrRefreshSessionExpired
	}

	if refreshSession.IsBlocked {
		return domain.RefreshSession{}, ErrRefreshSessionBlocked
	}

	return refreshSession, nil
}

func (c *authUseCase) UserSignUp(ctx context.Context, signUpDetails domain.User) (string, error) {

	existUser, err := c.userRepo.FindUserByUserNameEmailOrPhoneNotID(ctx, signUpDetails)
	if err != nil {
		return "", utils.PrependMessageToError(err, "failed to check user details already exist")
	}

	// if user credentials already exist and  verified then return it as errors
	if existUser.ID != 0 && existUser.Verified {
		err = utils.CompareUserExistingDetails(existUser, signUpDetails)
		err = utils.AppendMessageToError(ErrUserAlreadyExit, err.Error())
		return "", err
	}

	errChan := make(chan error, 2)
	wait := sync.WaitGroup{}
	wait.Add(2)

	go func() {
		defer wait.Done()
		_, err := c.optAuth.SentOtp(countryCode + signUpDetails.Phone)
		if err != nil {
			errChan <- fmt.Errorf("failed to send otp \nerrors:%v", err.Error())
		}
	}()

	userID := existUser.ID

	if userID == 0 { // if user not exist then save user on database
		hashPass, err := utils.GenerateHashFromPassword(signUpDetails.Password)
		if err != nil {
			return "", utils.PrependMessageToError(err, "failed to hash the password")
		}

		signUpDetails.Password = string(hashPass)
		userID, err = c.userRepo.SaveUser(ctx, signUpDetails)

		if err != nil {
			return "", utils.PrependMessageToError(err, "failed to save user details")
		}
	}

	otpID := uuid.NewString()

	go func() {
		defer wait.Done()
		otpSession := domain.OtpSession{
			OtpID:    otpID,
			UserID:   userID,
			Phone:    signUpDetails.Phone,
			ExpireAt: time.Now().Add(otpExpireDuration), // 2 minutes expire for otp
		}
		err := c.authRepo.SaveOtpSession(ctx, otpSession)
		if err != nil {
			errChan <- fmt.Errorf("failed to save otp session \nerror:%v", err.Error())
		}
	}()

	wait.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return "", err
		}
	}

	return otpID, nil
}

func (c *authUseCase) SingUpOtpVerify(ctx context.Context,
	otpVerifyDetails request.OTPVerify) (userID uint, err error) {

	otpSession, err := c.authRepo.FindOtpSession(ctx, otpVerifyDetails.OtpID)
	if err != nil {
		return 0, utils.PrependMessageToError(err, "failed to find otp session from database")
	}

	if time.Since(otpSession.ExpireAt) > 0 {
		return 0, ErrOtpExpired
	}

	valid, err := c.optAuth.VerifyOtp(countryCode+otpSession.Phone, otpVerifyDetails.Otp)
	if err != nil {
		return 0, utils.PrependMessageToError(err, "failed to verify otp")
	}
	if !valid {
		return 0, ErrInvalidOtp
	}

	err = c.userRepo.UpdateVerified(ctx, otpSession.ID)
	if err != nil {
		return 0, utils.PrependMessageToError(err, "failed to update user verified on database")
	}

	return otpSession.UserID, nil
}

// google login
func (c *authUseCase) GoogleLogin(ctx context.Context, user domain.User) (userID uint, err error) {

	existUser, err := c.userRepo.FindUserByEmail(ctx, user.Email)
	if err != nil {
		return userID, fmt.Errorf("failed to get user details with given email \nerror:%v", err.Error())
	}

	if existUser.ID != 0 {
		return existUser.ID, nil
	}

	// create a random user name for user based on user name
	user.UserName = utils.GenerateRandomUserName(user.FirstName)

	userID, err = c.userRepo.SaveUser(ctx, user)
	if err != nil {
		return userID, fmt.Errorf("failed to save user details \nerror:%v", err.Error())
	}

	return userID, nil
}
