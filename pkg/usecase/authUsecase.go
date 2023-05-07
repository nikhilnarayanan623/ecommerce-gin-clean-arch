package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
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

func NewAuthUseCase(authRepo interfaces.AuthRepository, tokenAuth token.TokenAuth, userRepo interfaces.UserRepository, adminRepo interfaces.AdminRepository) service.AuthUseCase {

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
