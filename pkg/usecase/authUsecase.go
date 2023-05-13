package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/otp"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/token"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type authUseCase struct {
	authRepo interfaces.AuthRepository

	userRepo  interfaces.UserRepository
	adminRepo interfaces.AdminRepository
	tokenAuth token.TokenAuth
	otpVerify otp.OtpVerification
}

func NewAuthUseCase(authRepo interfaces.AuthRepository, tokenAuth token.TokenAuth,
	userRepo interfaces.UserRepository, adminRepo interfaces.AdminRepository,
	otpVeriy otp.OtpVerification) service.AuthUseCase {

	return &authUseCase{
		userRepo:  userRepo,
		adminRepo: adminRepo,
		tokenAuth: tokenAuth,
		authRepo:  authRepo,
		otpVerify: otpVeriy,
	}
}

func (c *authUseCase) UserLogin(ctx context.Context, loginDetails req.Login) (userID uint, err error) {

	var user domain.User

	if loginDetails.Email != "" {
		user, err = c.userRepo.FindUserByEmail(ctx, loginDetails.Email)
	} else if loginDetails.UserName != "" {
		user, err = c.userRepo.FindUserByUserName(ctx, loginDetails.UserName)
	} else if loginDetails.Phone != "" {
		user, err = c.userRepo.FindUserByPhoneNumber(ctx, loginDetails.Phone)
	} else {
		return userID, fmt.Errorf("all user login unique fields are empty")
	}

	if err != nil {
		return userID, fmt.Errorf("an error found when find user \nerror: %v", err.Error())
	}

	if user.ID == 0 {
		return userID, fmt.Errorf("user not exist with given lgoin details")
	}

	if user.BlockStatus {
		return userID, fmt.Errorf("the user blocked by admin")
	}

	err = utils.ComparePasswordWithHashedPassword(loginDetails.Password, user.Password)
	if err != nil {
		return userID, fmt.Errorf("given password is wrong")
	}

	return user.ID, err
}

func (c *authUseCase) UserLoginOtpSend(ctx context.Context, loginDetails req.OTPLogin) (otpRes res.OTPResponse, err error) {

	var user domain.User
	if loginDetails.Email != "" {
		user, err = c.userRepo.FindUserByEmail(ctx, loginDetails.Email)
	} else if loginDetails.UserName != "" {
		user, err = c.userRepo.FindUserByUserName(ctx, loginDetails.UserName)
	} else if loginDetails.Phone != "" {
		user, err = c.userRepo.FindUserByPhoneNumber(ctx, loginDetails.Phone)
	} else {
		return otpRes, fmt.Errorf("all user login unique fields are empty")
	}

	if err != nil {
		return otpRes, fmt.Errorf("can't find the user \nerror:%v", err.Error())
	} else if user.ID == 0 {
		return otpRes, errors.New("user not exist with this details")
	}

	// check user block_status user is blocked or not
	if user.BlockStatus {
		return otpRes, errors.New("user blocked by admin")
	}

	_, err = c.otpVerify.SentOtp("+91" + user.Phone)

	if err != nil {
		return otpRes, fmt.Errorf("faild to send otp \nerrors:%v", err.Error())
	}

	otpRes.OTPID, err = uuid.NewRandom()
	if err != nil {
		return otpRes, fmt.Errorf("faild to create otp_id")
	}

	otpSession := domain.OtpSession{
		OTPID:    otpRes.OTPID,
		UserID:   user.ID,
		Phone:    user.Phone,
		ExpireAt: time.Now().Add(time.Minute * 2),
	}
	err = c.authRepo.SaveOtpSession(ctx, otpSession)

	if err != nil {
		return otpRes, fmt.Errorf("faild to save otp session \nerror:%v", err.Error())
	}
	return otpRes, nil
}

func (c *authUseCase) LoginOtpVerify(ctx context.Context, otpVeirifyDetails req.OTPVerify) (userID uint, err error) {

	otpSession, err := c.authRepo.FindOtpSession(ctx, otpVeirifyDetails.OTPID)
	if err != nil {
		return userID, fmt.Errorf("faild to get otp session \nerror:%v", err.Error())
	}

	if otpSession.ID == 0 {
		return userID, fmt.Errorf("invlaid otp sid")
	}

	if time.Since(otpSession.ExpireAt) > 0 {
		return userID, fmt.Errorf("opt expired")
	}

	err = c.otpVerify.VerifyOtp("+91"+otpSession.Phone, otpVeirifyDetails.OTP)
	if err != nil {
		return userID, fmt.Errorf("faild to verify otp \nerror:%v", err.Error())
	}

	return otpSession.UserID, err
}

func (c *authUseCase) AdminLogin(ctx context.Context, loginDetails req.Login) (adminID uint, err error) {

	var admin domain.Admin

	if loginDetails.Email != "" {
		admin, err = c.adminRepo.FindAdminByEmail(ctx, loginDetails.Email)
	} else if loginDetails.UserName != "" {
		admin, err = c.adminRepo.FindAdminByUserName(ctx, loginDetails.UserName)
	} else {
		return adminID, fmt.Errorf("all admin login unique fields are empty")
	}

	if err != nil {
		return adminID, fmt.Errorf("an error found when find user \nerror: %v", err.Error())
	}

	if admin.ID == 0 {
		return adminID, fmt.Errorf("admin not exist with given lgoin details")
	}

	err = utils.ComparePasswordWithHashedPassword(loginDetails.Password, admin.Password)
	if err != nil {
		return adminID, fmt.Errorf("given password is wrong")
	}

	return admin.ID, err
}

func (c *authUseCase) GenerateAccessToken(ctx context.Context, tokenParams service.GenerateTokenParams) (tokenString string, err error) {

	uniqueID, err := uuid.NewRandom()
	if err != nil {
		return tokenString, nil
	}
	payload := &token.Payload{
		TokenID:  uniqueID,
		UserID:   tokenParams.UserID,
		ExpireAt: tokenParams.ExpireDate,
	}
	tokenString, err = c.tokenAuth.CreateToken(payload, tokenParams.UserType)

	return
}
func (c *authUseCase) GenerateRefreshToken(ctx context.Context, tokenParams service.GenerateTokenParams) (tokenString string, err error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	if time.Since(tokenParams.ExpireDate) > 0 {
		return
	}

	payload := &token.Payload{
		TokenID:  tokenID,
		UserID:   tokenParams.UserID,
		ExpireAt: tokenParams.ExpireDate,
	}
	tokenString, err = c.tokenAuth.CreateToken(payload, tokenParams.UserType)
	if err != nil {
		return "", err
	}

	err = c.authRepo.SaveRefreshSession(ctx, domain.RefreshSession{
		UserID:       payload.UserID,
		TokenID:      tokenID,
		RefreshToken: tokenString,
		ExpireAt:     tokenParams.ExpireDate,
	})
	if err != nil {
		return "", err
	}
	log.Printf("successfully refresh token created and refresh session stored in database")
	return tokenString, nil
}

func (c *authUseCase) VerifyAndGetRefreshTokenSession(ctx context.Context, refreshToken string, usedFor token.UserType) (domain.RefreshSession, error) {

	payload, err := c.tokenAuth.VerifyToken(refreshToken, usedFor)
	if err != nil {
		return domain.RefreshSession{}, err
	}

	refreshSession, err := c.authRepo.FindRefreshSessionByTokenID(ctx, payload.TokenID)
	if err != nil {
		return refreshSession, err
	}

	if refreshSession.TokenID == uuid.Nil {
		return refreshSession, errors.New("there is no refresh token session for this token")
	}

	if time.Since(refreshSession.ExpireAt) > 0 {
		return domain.RefreshSession{}, errors.New("given refresh token's session expired")
	}

	if refreshSession.IsBlocked {
		return domain.RefreshSession{}, errors.New("given refresh token is blocked")
	}

	return refreshSession, nil
}

// google login
func (c *authUseCase) GoogleLogin(ctx context.Context, user domain.User) (userID uint, err error) {

	existUser, err := c.userRepo.FindUserByEmail(ctx, user.Email)
	if err != nil {
		return userID, fmt.Errorf("faild to get user details with given email \nerror:%v", err.Error())
	} else if existUser.ID != 0 {
		return existUser.ID, nil
	}

	user.UserName = utils.GenerateRandomUserName(user.FirstName)
	userID, err = c.userRepo.SaveUser(ctx, user)
	if err != nil {
		return userID, fmt.Errorf("faild to save user details \nerror:%v", err.Error())
	}

	return userID, nil
}
