package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/mock/mockrepo"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/mock/mockservice"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/token"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/stretchr/testify/assert"
)

// func createRandonUser() domain.User {
// 	return domain.User{
// 		ID:       uint(utils.RandomInt(1, 10)),
// 		UserName: utils.GenerateRandomString(6),
// 	}
// }

type LoginKey string

const (
	Email       LoginKey = "Email"
	PhoneNumber LoginKey = "PhoneNumber"
	UserName    LoginKey = "UserName"
)

func createRandomLoginDetail(loginKey LoginKey) request.Login {

	if loginKey == Email {
		return request.Login{
			Email:    utils.GenerateRandomString(7),
			Password: utils.GenerateRandomString(12),
		}
	} else if loginKey == PhoneNumber {
		return request.Login{
			Phone:    utils.GenerateRandomString(10),
			Password: utils.GenerateRandomString(12),
		}
	}
	return request.Login{
		UserName: utils.GenerateRandomString(12),
		Password: utils.GenerateRandomString(12),
	}
}

func TestUserLogin(t *testing.T) {

	tests := []struct {
		testName       string
		input          request.Login
		expectedOutput uint
		buildStub      func(mockRepo *mockrepo.MockUserRepository, loginDetails request.Login)
		expectedError  error
	}{
		{
			testName: "EmptyLoginCredentialsShouldReturnError",
			input:    request.Login{},
			buildStub: func(mockRepo *mockrepo.MockUserRepository, loginDetails request.Login) {
				//not expecting any call to mockRepo
			},
			expectedOutput: 0,
			expectedError:  ErrEmptyLoginCredentials,
		},

		{
			testName:       "EmailExistShouldCallFindUserByEmailWithGivenEmail",
			input:          request.Login{Email: "emailExist@gmail.com", Password: "password"},
			expectedOutput: 1,
			buildStub: func(mockRepo *mockrepo.MockUserRepository, loginDetails request.Login) {
				hashedPassword, err := utils.GetHashedPassword(loginDetails.Password)
				assert.NoError(t, err)
				outputUser := domain.User{ID: 1, Email: loginDetails.Email, Password: hashedPassword}
				mockRepo.EXPECT().FindUserByEmail(gomock.Any(), loginDetails.Email).
					Times(1).Return(outputUser, nil)
			},
			expectedError: nil,
		},
		{
			testName:       "FindUserErrorShouldReturnError",
			input:          createRandomLoginDetail(UserName),
			expectedOutput: 0,
			buildStub: func(mockRepo *mockrepo.MockUserRepository, loginDetails request.Login) {
				dbError := fmt.Errorf("error from find user on database")
				mockRepo.EXPECT().
					FindUserByUserName(gomock.Any(), loginDetails.UserName).
					Times(1).Return(domain.User{}, dbError)
			},
			expectedError: fmt.Errorf("an error found when find user \nerror: %v", "error from find user on database"),
		},
		{
			testName:       "NonExistingEmailShouldReturnErrorOfUserNotExist",
			input:          request.Login{Email: "nonExistingEmail@gmail.com"},
			expectedOutput: 0,
			buildStub: func(mockRepo *mockrepo.MockUserRepository, loginDetails request.Login) {
				outputUser := domain.User{}
				mockRepo.EXPECT().FindUserByEmail(gomock.Any(), loginDetails.Email).
					Times(1).Return(outputUser, nil)
			},
			expectedError: ErrUserNotExist,
		},
		{
			testName:       "UserBlockedByAdminShouldReturnError",
			input:          createRandomLoginDetail(UserName),
			expectedOutput: 0,
			buildStub: func(mockRepo *mockrepo.MockUserRepository, loginDetails request.Login) {
				outputUser := domain.User{ID: 1, BlockStatus: true}
				mockRepo.EXPECT().FindUserByUserName(gomock.Any(), loginDetails.UserName).
					Times(1).Return(outputUser, nil)
			},
			expectedError: ErrUserBlocked,
		},
		{
			testName:       "UserExistPasswordNotMatchWithHashedPasswordShouldReturnError",
			input:          createRandomLoginDetail(Email),
			expectedOutput: 0,
			buildStub: func(mockRepo *mockrepo.MockUserRepository, loginDetails request.Login) {
				_, err := utils.GetHashedPassword(loginDetails.Password)
				assert.NoError(t, err)
				outputUser := domain.User{ID: 1, Password: "hashedPassword"}
				mockRepo.EXPECT().FindUserByEmail(gomock.Any(), loginDetails.Email).
					Times(1).Return(outputUser, nil)
			},
			expectedError: ErrWrongPassword,
		},
		{
			testName:       "ValidLoginDetailsShouldReturnUserIDWithNorError",
			input:          createRandomLoginDetail(PhoneNumber),
			expectedOutput: 1,
			buildStub: func(mockRepo *mockrepo.MockUserRepository, loginDetails request.Login) {
				hashedPassword, err := utils.GetHashedPassword(loginDetails.Password)
				assert.NoError(t, err)
				outputUser := domain.User{ID: 1, Password: hashedPassword}
				mockRepo.EXPECT().FindUserByPhoneNumber(gomock.Any(), loginDetails.Phone).
					Times(1).Return(outputUser, nil)
			},
		},
	}

	for _, test := range tests {
		var test = test
		t.Run(test.testName, func(t *testing.T) {
			t.Parallel()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			userMockRepo := mockrepo.NewMockUserRepository(ctl)
			test.buildStub(userMockRepo, test.input)

			authUseCase := NewAuthUseCase(nil, nil, userMockRepo, nil, nil)
			actualOutput, actualError := authUseCase.UserLogin(context.Background(), test.input)

			if test.expectedError != nil {
				assert.Error(t, test.expectedError, actualError)
			} else {
				assert.NoError(t, actualError)
			}
			assert.Equal(t, test.expectedOutput, actualOutput)
		})

	}

}

// func TestGenerateAccessToken(t *testing.T) {
// 	type inputFields struct {
// 		userID   uint
// 		userType string
// 	}
// 	tests := []struct {
// 		testName    string
// 		inputFields inputFields
// 	}{}

// 	for _, test := range tests {
// 		t.Run(test.testName, func(t *testing.T) {

// 			authUseCase := NewAuthUseCase(nil, nil, nil, config.Config{JWT: "secretKey"})

// 			accessToken, err := authUseCase.GenerateAccessToken(context.Background(), test.inputFields.userID, test.inputFields.userType)

// 		})
// 	}
// }

func TestGenerateRefreshToken(t *testing.T) {

	tests := []struct {
		testName              string
		inputField            service.GenerateTokenParams
		buildStubAuthRepo     func(t *testing.T, mockRepo *mockrepo.MockAuthRepository)
		buildStubTokenService func(t *testing.T, tService *mockservice.MockTokenService)
		checkOutput           func(t *testing.T, tokenString string, err error)
	}{
		{

			testName: "FailedToCreateTokenShouldReturnError",
			inputField: service.GenerateTokenParams{
				UserID:   2,
				UserType: token.User,
			},
			buildStubTokenService: func(t *testing.T, tService *mockservice.MockTokenService) {
				tService.EXPECT().GenerateToken(gomock.Any()).
					Times(1).Return(token.GenerateTokenResponse{}, errors.New("failed to generate token"))
			},
			buildStubAuthRepo: func(t *testing.T, mockRepo *mockrepo.MockAuthRepository) {
				// not expecting calls to any functions
			},
			checkOutput: func(t *testing.T, tokenString string, actualError error) {
				assert.Empty(t, tokenString)
				assert.Equal(t, errors.New("failed to generate token"), actualError)
			},
		},
		{
			testName: "FailedToSaveRefreshTokenOnSession",
			inputField: service.GenerateTokenParams{UserID: 1,
				UserType: token.User,
			},
			buildStubTokenService: func(t *testing.T, tokenAuth *mockservice.MockTokenService) {
				tokenAuth.EXPECT().GenerateToken(gomock.Any()).
					Times(1).Return(token.GenerateTokenResponse{TokenString: "access_token"}, nil)
			},
			buildStubAuthRepo: func(t *testing.T, mockRepo *mockrepo.MockAuthRepository) {
				mockRepo.EXPECT().SaveRefreshSession(gomock.Any(), gomock.Any()).
					Times(1).Return(errors.New("failed to save refresh_session on database"))
			},
			checkOutput: func(t *testing.T, tokenString string, err error) {
				assert.EqualError(t, err, "failed to save refresh_session on database")
				assert.Empty(t, tokenString)
			},
		},
		{
			testName: "SuccessfulRefreshTokenCreationAndSaveReturnToken",
			inputField: service.GenerateTokenParams{UserID: 1,
				UserType: token.User,
			},
			buildStubTokenService: func(t *testing.T, tokenAuth *mockservice.MockTokenService) {
				tokenAuth.EXPECT().GenerateToken(gomock.Any()).
					Times(1).Return(token.GenerateTokenResponse{TokenID: "token_id", TokenString: "token"}, nil)
			},
			buildStubAuthRepo: func(t *testing.T, mockRepo *mockrepo.MockAuthRepository) {
				mockRepo.EXPECT().SaveRefreshSession(gomock.Any(), gomock.Any()).
					Times(1).Return(nil)
			},
			checkOutput: func(t *testing.T, tokenString string, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, tokenString)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			ctl := gomock.NewController(t)
			mockAuthRepo := mockrepo.NewMockAuthRepository(ctl)
			mockTokenAuth := mockservice.NewMockTokenService(ctl)
			test.buildStubAuthRepo(t, mockAuthRepo)
			test.buildStubTokenService(t, mockTokenAuth)

			authUseCase := NewAuthUseCase(mockAuthRepo, mockTokenAuth, nil, nil, nil)
			tokenString, err := authUseCase.GenerateRefreshToken(context.Background(), test.inputField)

			test.checkOutput(t, tokenString, err)
		})
	}
}

func TestVerifyAndGetRefreshTokenSession(t *testing.T) {

	tokenUser := token.User
	tests := []struct {
		testName       string
		refreshToken   string
		buildStub      func(authMockRepo *mockrepo.MockAuthRepository, tokenMockAuth *mockservice.MockTokenService)
		expectedOutput domain.RefreshSession
		expectedError  error
	}{
		{
			testName:     "InvalidRefreshTokenShouldReturnError",
			refreshToken: "invalidRefreshToken",
			buildStub: func(authMockRepo *mockrepo.MockAuthRepository, tokenMockAuth *mockservice.MockTokenService) {
				tokenMockAuth.EXPECT().VerifyToken(token.VerifyTokenRequest{TokenString: "invalidRefreshToken", UsedFor: tokenUser}).
					Times(1).Return(token.VerifyTokenResponse{}, errors.New("invalid refresh token"))
			},
			expectedOutput: domain.RefreshSession{},
			expectedError:  errors.New("invalid refresh token"),
		},
		{
			testName:     "FailedToFindRefreshSessionOnDBShouldReturnError",
			refreshToken: "refreshToken",
			buildStub: func(authMockRepo *mockrepo.MockAuthRepository, tokenMockAuth *mockservice.MockTokenService) {

				tokenMockAuth.EXPECT().VerifyToken(token.VerifyTokenRequest{TokenString: "refreshToken", UsedFor: tokenUser}).
					Times(1).Return(token.VerifyTokenResponse{TokenID: "token_id", UserID: 12}, nil)

				authMockRepo.EXPECT().FindRefreshSessionByTokenID(gomock.Any(), "token_id").
					Times(1).Return(domain.RefreshSession{}, errors.New("error when finding refresh token"))
			},
			expectedOutput: domain.RefreshSession{},
			expectedError:  errors.New("error when finding refresh token"),
		},
		{
			testName:     "RefreshTokenNotExistInRefreshSession",
			refreshToken: "NonExistingRefreshToken",
			buildStub: func(authMockRepo *mockrepo.MockAuthRepository, tokenMockAuth *mockservice.MockTokenService) {

				tokenMockAuth.EXPECT().VerifyToken(token.VerifyTokenRequest{TokenString: "NonExistingRefreshToken", UsedFor: token.User}).
					Times(1).Return(token.VerifyTokenResponse{TokenID: "no_existing_token_id", UserID: 12}, nil)
				authMockRepo.EXPECT().FindRefreshSessionByTokenID(gomock.Any(), "no_existing_token_id").
					Times(1).Return(domain.RefreshSession{}, nil)
			},
			expectedOutput: domain.RefreshSession{},
			expectedError:  ErrRefreshSessionNotExist,
		},
		{
			testName:     "BlockedRefreshTokenShouldReturnError",
			refreshToken: "validRefreshToken",
			buildStub: func(authMockRepo *mockrepo.MockAuthRepository, tokenMockAuth *mockservice.MockTokenService) {
				tokenMockAuth.EXPECT().VerifyToken(token.VerifyTokenRequest{TokenString: "validRefreshToken", UsedFor: tokenUser}).
					Times(1).Return(token.VerifyTokenResponse{TokenID: "token_id", UserID: 12}, nil)
				authMockRepo.EXPECT().FindRefreshSessionByTokenID(gomock.Any(), "token_id").
					Times(1).Return(domain.RefreshSession{TokenID: "token_id", IsBlocked: true,
					ExpireAt: time.Now().Add(time.Hour * 2)}, nil)
			},
			expectedOutput: domain.RefreshSession{},
			expectedError:  ErrRefreshSessionBlocked,
		},
		{
			testName:     "RefreshTokenSessionExpiredShouldReturnError",
			refreshToken: "validExistingRefresh",
			buildStub: func(authMockRepo *mockrepo.MockAuthRepository, tokenMockAuth *mockservice.MockTokenService) {

				tokenMockAuth.EXPECT().VerifyToken(token.VerifyTokenRequest{TokenString: "validExistingRefresh", UsedFor: tokenUser}).
					Times(1).Return(token.VerifyTokenResponse{TokenID: "token_id", UserID: 12}, nil)
				expiredTokenSession := domain.RefreshSession{TokenID: "token_id",
					ExpireAt: time.Date(2000, 12, 12, 12, 12, 12, 12, time.UTC)}

				authMockRepo.EXPECT().FindRefreshSessionByTokenID(gomock.Any(), "token_id").
					Times(1).Return(expiredTokenSession, nil)
			},
			expectedOutput: domain.RefreshSession{},
			expectedError:  ErrRefreshSessionExpired,
		},
		{
			testName:     "ValidExistingTokenIdShouldReturnRefreshSession",
			refreshToken: "validExistingRefresh",
			buildStub: func(authMockRepo *mockrepo.MockAuthRepository, tokenMockAuth *mockservice.MockTokenService) {

				tokenMockAuth.EXPECT().VerifyToken(token.VerifyTokenRequest{TokenString: "validExistingRefresh", UsedFor: tokenUser}).
					Times(1).Return(token.VerifyTokenResponse{TokenID: "token_id", UserID: 12}, nil)

				refreshSession := domain.RefreshSession{TokenID: "token_id", IsBlocked: false,
					ExpireAt: time.Date(3000, 12, 12, 12, 12, 12, 12, time.UTC)}

				authMockRepo.EXPECT().FindRefreshSessionByTokenID(gomock.Any(), "token_id").
					Times(1).Return(refreshSession, nil)
			},
			expectedOutput: domain.RefreshSession{TokenID: "token_id", IsBlocked: false,
				ExpireAt: time.Date(3000, 12, 12, 12, 12, 12, 12, time.UTC)},
			expectedError: nil,
		},
	}

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			ctl := gomock.NewController(t)
			authMockRepo := mockrepo.NewMockAuthRepository(ctl)
			tokenService := mockservice.NewMockTokenService(ctl)

			authUseCase := NewAuthUseCase(authMockRepo, tokenService, nil, nil, nil)

			test.buildStub(authMockRepo, tokenService)

			actualOutput, actualError := authUseCase.
				VerifyAndGetRefreshTokenSession(context.Background(), test.refreshToken, tokenUser)

			if test.expectedError == nil {
				assert.NoError(t, actualError)
			} else {
				assert.Error(t, test.expectedError, actualError)
			}
			assert.Equal(t, test.expectedOutput, actualOutput)
		})
	}
}
