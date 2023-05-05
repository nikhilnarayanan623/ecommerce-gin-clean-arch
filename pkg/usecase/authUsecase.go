package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/token"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
)

type authUseCase struct {
	authRepo  interfaces.AuthRepository
	userRepo  interfaces.UserRepository
	adminRepo interfaces.AdminRepository
	tokenAuth token.TokenAuth
}

func NewAuthUseCase(userRepo interfaces.UserRepository, adminRepo interfaces.AdminRepository, authRepo interfaces.AuthRepository, cfg config.Config) service.AuthUseCase {
	tokenAuth := token.NewJWTAuth(cfg.JWTAdmin, cfg.JWTUser)
	return &authUseCase{
		userRepo:  userRepo,
		adminRepo: adminRepo,
		tokenAuth: tokenAuth,
		authRepo:  authRepo,
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

func (c *authUseCase) GenerateAccessToken(ctx context.Context, userID uint, userType token.UserType, expireTimeDuration time.Duration) (tokenString string, err error) {

	uniqueID, err := uuid.NewRandom()
	if err != nil {
		return tokenString, nil
	}
	payload := &token.Payload{
		TokenID:  uniqueID,
		UserID:   userID,
		ExpireAt: time.Now().Add(expireTimeDuration),
	}
	tokenString, err = c.tokenAuth.CreateToken(payload, userType)

	return
}
func (c *authUseCase) GenerateRefreshToken(ctx context.Context, userID uint, userType token.UserType, expireTimeDuration time.Duration) (tokenString string, err error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return tokenString, nil
	}
	refresTokenExpire := time.Now().Add(expireTimeDuration)

	payload := &token.Payload{
		TokenID:  tokenID,
		UserID:   userID,
		ExpireAt: refresTokenExpire,
	}
	tokenString, err = c.tokenAuth.CreateToken(payload, userType)
	if err != nil {
		return tokenString, err
	}

	err = c.authRepo.SaveRefreshSession(ctx, domain.RefreshSession{
		TokenID:      tokenID,
		RefreshToken: tokenString,
		ExpireAt:     refresTokenExpire,
	})
	if err != nil {
		return "", err
	}
	log.Printf("successfully refresh token created and refresh session stored in database")
	return
}
