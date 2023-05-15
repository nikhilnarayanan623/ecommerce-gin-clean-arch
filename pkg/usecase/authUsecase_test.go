package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/mock/mockRepository"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/mock/mockTokenAuth"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/token"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
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

func createRandomLoginDetail(loginKey LoginKey) req.Login {

	if loginKey == Email {
		return req.Login{
			Email:    utils.GenerateRandomString(7),
			Password: utils.GenerateRandomString(12),
		}
	} else if loginKey == PhoneNumber {
		return req.Login{
			Phone:    utils.GenerateRandomString(10),
			Password: utils.GenerateRandomString(12),
		}
	}
	return req.Login{
		UserName: utils.GenerateRandomString(12),
		Password: utils.GenerateRandomString(12),
	}
}

func TestUserLogin(t *testing.T) {

	tests := []struct {
		testName       string
		input          req.Login
		expectedOutput uint
		buildStub      func(mockRepo *mockRepository.MockUserRepository, loginDetails req.Login)
		expectedError  error
	}{
		{
			testName: "EmptyLoginStructReturnErrorOfNoUniqueFields",
			input:    req.Login{},
			buildStub: func(mockRepo *mockRepository.MockUserRepository, loginDetails req.Login) {
				//not expecting any call to mockRepo
			},
			expectedOutput: 0,
			expectedError:  fmt.Errorf("all user login unique fields are empty"),
		},

		{
			testName:       "EmailExistShouldCallFindUserByEmailWithGivenEmail",
			input:          req.Login{Email: "emailExist@gmail.com", Password: "password"},
			expectedOutput: 1,
			buildStub: func(mockRepo *mockRepository.MockUserRepository, loginDetails req.Login) {
				hashedPasword, err := utils.GetHashedPassword(loginDetails.Password)
				assert.NoError(t, err)
				outputUser := domain.User{ID: 1, Email: loginDetails.Email, Password: hashedPasword}
				mockRepo.EXPECT().FindUserByEmail(gomock.Any(), loginDetails.Email).
					Times(1).Return(outputUser, nil)
			},
			expectedError: nil,
		},
		{
			testName:       "FindUserErrorShouldReturnError",
			input:          createRandomLoginDetail(UserName),
			expectedOutput: 0,
			buildStub: func(mockRepo *mockRepository.MockUserRepository, loginDetails req.Login) {
				dbError := fmt.Errorf("error from find user on database")
				mockRepo.EXPECT().
					FindUserByUserName(gomock.Any(), loginDetails.UserName).
					Times(1).Return(domain.User{}, dbError)
			},
			expectedError: fmt.Errorf("an error found when find user \nerror: %v", "error from find user on database"),
		},
		{
			testName:       "NotExistingEmailShouldReturnErrorOfUserNotExist",
			input:          req.Login{Email: "nonExistingEmail@gmail.com"},
			expectedOutput: 0,
			buildStub: func(mockRepo *mockRepository.MockUserRepository, loginDetails req.Login) {
				outputUser := domain.User{}
				mockRepo.EXPECT().FindUserByEmail(gomock.Any(), loginDetails.Email).
					Times(1).Return(outputUser, nil)
			},
			expectedError: fmt.Errorf("user not exist with given login details"),
		},
		{
			testName:       "UserBlockedByAdminShouldReturnError",
			input:          createRandomLoginDetail(UserName),
			expectedOutput: 0,
			buildStub: func(mockRepo *mockRepository.MockUserRepository, loginDetails req.Login) {
				outputUser := domain.User{ID: 1, BlockStatus: true}
				mockRepo.EXPECT().FindUserByUserName(gomock.Any(), loginDetails.UserName).
					Times(1).Return(outputUser, nil)
			},
			expectedError: fmt.Errorf("the user blocked by admin"),
		},
		{
			testName:       "UserExistPasswordNotMatchWithHashedPasswordShouldReturnError",
			input:          createRandomLoginDetail(Email),
			expectedOutput: 0,
			buildStub: func(mockRepo *mockRepository.MockUserRepository, loginDetails req.Login) {
				_, err := utils.GetHashedPassword(loginDetails.Password)
				assert.NoError(t, err)
				outputUser := domain.User{ID: 1, Password: "hashedPassword"}
				mockRepo.EXPECT().FindUserByEmail(gomock.Any(), loginDetails.Email).
					Times(1).Return(outputUser, nil)
			},
			expectedError: fmt.Errorf("given password is wrong"),
		},
		{
			testName:       "ValidLoginDetailsShouldReturnUserIDWithNorError",
			input:          createRandomLoginDetail(PhoneNumber),
			expectedOutput: 1,
			buildStub: func(mockRepo *mockRepository.MockUserRepository, loginDetails req.Login) {
				hashedPassword, err := utils.GetHashedPassword(loginDetails.Password)
				assert.NoError(t, err)
				outputUser := domain.User{ID: 1, Password: hashedPassword}
				mockRepo.EXPECT().FindUserByPhoneNumber(gomock.Any(), loginDetails.Phone).
					Times(1).Return(outputUser, nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			t.Parallel()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			userMockRepo := mockRepository.NewMockUserRepository(ctl)
			test.buildStub(userMockRepo, test.input)

			authUseCase := NewAuthUseCase(nil, nil, userMockRepo, nil, nil)
			actualOutput, actualError := authUseCase.UserLogin(context.Background(), test.input)

			if test.expectedError != nil {
				assert.Equal(t, test.expectedError, actualError)
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
		testName           string
		inputField         service.GenerateTokenParams
		buildStubAuthRepo  func(t *testing.T, mockRepo *mockRepository.MockAuthRepository)
		buildStubTokenAuth func(t *testing.T, tokenAuth *mockTokenAuth.MockTokenAuth)
		checkOutput        func(t *testing.T, tokenString string, err error)
	}{
		{

			testName: "FaildToCreateTokenShouldReturnError",
			inputField: service.GenerateTokenParams{
				UserID:     2,
				UserType:   token.TokenForUser,
				ExpireDate: time.Now().Add(time.Minute * 20),
			},
			buildStubTokenAuth: func(t *testing.T, tokenAuth *mockTokenAuth.MockTokenAuth) {
				tokenAuth.EXPECT().CreateToken(gomock.Any(), gomock.Any()).
					Times(1).Return("", errors.New("faild to create token"))
			},
			buildStubAuthRepo: func(t *testing.T, mockRepo *mockRepository.MockAuthRepository) {
				// not expecting calls to any functions
			},
			checkOutput: func(t *testing.T, tokenString string, actualError error) {
				assert.Empty(t, tokenString)
				assert.Equal(t, errors.New("faild to create token"), actualError)
			},
		},
		{
			testName: "FaildToSaveRefreshTokenOnSession",
			inputField: service.GenerateTokenParams{UserID: 1,
				UserType:   token.TokenForUser,
				ExpireDate: time.Now().Add(time.Minute * 15),
			},
			buildStubTokenAuth: func(t *testing.T, tokenAuth *mockTokenAuth.MockTokenAuth) {
				tokenAuth.EXPECT().CreateToken(gomock.Any(), gomock.Any()).
					Times(1).Return("generatedAccessToken", nil)
			},
			buildStubAuthRepo: func(t *testing.T, mockRepo *mockRepository.MockAuthRepository) {
				mockRepo.EXPECT().SaveRefreshSession(gomock.Any(), gomock.Any()).
					Times(1).Return(errors.New("faild to save refresh_session on database"))
			},
			checkOutput: func(t *testing.T, tokenString string, err error) {
				assert.EqualError(t, err, "faild to save refresh_session on database")
				assert.Empty(t, tokenString)
			},
		},
		{
			testName: "SuccesfullRefreshTokenCreationAndSaveReturnToken",
			inputField: service.GenerateTokenParams{UserID: 1,
				UserType:   token.TokenForUser,
				ExpireDate: time.Now().Add(time.Minute * 15),
			},
			buildStubTokenAuth: func(t *testing.T, tokenAuth *mockTokenAuth.MockTokenAuth) {
				tokenAuth.EXPECT().CreateToken(gomock.Any(), gomock.Any()).
					Times(1).Return("generatedAccessToken", nil)
			},
			buildStubAuthRepo: func(t *testing.T, mockRepo *mockRepository.MockAuthRepository) {
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
			mockAuthRepo := mockRepository.NewMockAuthRepository(ctl)
			mockTokenAuth := mockTokenAuth.NewMockTokenAuth(ctl)
			test.buildStubAuthRepo(t, mockAuthRepo)
			test.buildStubTokenAuth(t, mockTokenAuth)

			authUseCase := NewAuthUseCase(mockAuthRepo, mockTokenAuth, nil, nil, nil)
			tokenString, err := authUseCase.GenerateRefreshToken(context.Background(), test.inputField)

			test.checkOutput(t, tokenString, err)
		})
	}
}

func TestVerifyAndGetRefreshTokenSession(t *testing.T) {

	tokenUsedFor := token.TokenForUser
	randomUUID := uuid.New()
	tests := []struct {
		testName       string
		refreshToken   string
		buildStub      func(authMockRepo *mockRepository.MockAuthRepository, tokenMockAuth *mockTokenAuth.MockTokenAuth)
		expectedOutput domain.RefreshSession
		expectedError  error
	}{
		{
			testName:     "InvalidRefreshTokenShouldReturnError",
			refreshToken: "invalidRefreshToken",
			buildStub: func(authMockRepo *mockRepository.MockAuthRepository, tokenMockAuth *mockTokenAuth.MockTokenAuth) {
				tokenMockAuth.EXPECT().VerifyToken("invalidRefreshToken", tokenUsedFor).
					Times(1).Return(token.Payload{}, errors.New("invalid refresh token"))
			},
			expectedOutput: domain.RefreshSession{},
			expectedError:  errors.New("invalid refresh token"),
		},
		{
			testName:     "FaildToFindRefreshSessionFromDatabase",
			refreshToken: "validRefreshToken",
			buildStub: func(authMockRepo *mockRepository.MockAuthRepository, tokenMockAuth *mockTokenAuth.MockTokenAuth) {
				tokenMockAuth.EXPECT().VerifyToken("validRefreshToken", tokenUsedFor).
					Times(1).Return(token.Payload{TokenID: randomUUID, UserID: 12}, nil)
				authMockRepo.EXPECT().FindRefreshSessionByTokenID(gomock.Any(), randomUUID).
					Times(1).Return(domain.RefreshSession{}, errors.New("error when finding refresh token"))
			},
			expectedOutput: domain.RefreshSession{},
			expectedError:  errors.New("error when finding refresh token"),
		},
		{
			testName:     "RefreshTokenNotExistInRefreshSession",
			refreshToken: "NonExistingRefreshToken",
			buildStub: func(authMockRepo *mockRepository.MockAuthRepository, tokenMockAuth *mockTokenAuth.MockTokenAuth) {
				tokenMockAuth.EXPECT().VerifyToken("NonExistingRefreshToken", tokenUsedFor).
					Times(1).Return(token.Payload{TokenID: randomUUID, UserID: 12}, nil)
				authMockRepo.EXPECT().FindRefreshSessionByTokenID(gomock.Any(), randomUUID).
					Times(1).Return(domain.RefreshSession{}, nil)
			},
			expectedOutput: domain.RefreshSession{},
			expectedError:  errors.New("there is no refresh token session for this token"),
		},
		{
			testName:     "BlockedRefeshTokenShouldReturnError",
			refreshToken: "validRefreshToken",
			buildStub: func(authMockRepo *mockRepository.MockAuthRepository, tokenMockAuth *mockTokenAuth.MockTokenAuth) {
				tokenMockAuth.EXPECT().VerifyToken("validRefreshToken", tokenUsedFor).
					Times(1).Return(token.Payload{TokenID: randomUUID, UserID: 12}, nil)
				authMockRepo.EXPECT().FindRefreshSessionByTokenID(gomock.Any(), randomUUID).
					Times(1).Return(domain.RefreshSession{TokenID: randomUUID, IsBlocked: true,
					ExpireAt: time.Now().Add(time.Hour * 2)}, nil)
			},
			expectedOutput: domain.RefreshSession{},
			expectedError:  errors.New("given refresh token is blocked"),
		},
		{
			testName:     "RefreshTokenSessionExpiredShouldReturnError",
			refreshToken: "validExistingRefresh",
			buildStub: func(authMockRepo *mockRepository.MockAuthRepository, tokenMockAuth *mockTokenAuth.MockTokenAuth) {
				tokenMockAuth.EXPECT().VerifyToken("validExistingRefresh", tokenUsedFor).
					Times(1).Return(token.Payload{TokenID: randomUUID, UserID: 12}, nil)
				expiredTokenSession := domain.RefreshSession{TokenID: randomUUID,
					ExpireAt: time.Date(2000, 12, 12, 12, 12, 12, 12, time.UTC)}
				authMockRepo.EXPECT().FindRefreshSessionByTokenID(gomock.Any(), randomUUID).
					Times(1).Return(expiredTokenSession, nil)
			},
			expectedOutput: domain.RefreshSession{},
			expectedError:  errors.New("given refresh token's session expired"),
		},
		{
			testName:     "ValidExistingTokenIdShouldReturnRefreshSession",
			refreshToken: "validExistingRefresh",
			buildStub: func(authMockRepo *mockRepository.MockAuthRepository, tokenMockAuth *mockTokenAuth.MockTokenAuth) {
				tokenMockAuth.EXPECT().VerifyToken("validExistingRefresh", tokenUsedFor).
					Times(1).Return(token.Payload{TokenID: randomUUID, UserID: 12}, nil)
				refreshSession := domain.RefreshSession{TokenID: randomUUID, IsBlocked: false,
					ExpireAt: time.Date(3000, 12, 12, 12, 12, 12, 12, time.UTC)}
				authMockRepo.EXPECT().FindRefreshSessionByTokenID(gomock.Any(), randomUUID).
					Times(1).Return(refreshSession, nil)
			},
			expectedOutput: domain.RefreshSession{TokenID: randomUUID, IsBlocked: false,
				ExpireAt: time.Date(3000, 12, 12, 12, 12, 12, 12, time.UTC)},
			expectedError: nil,
		},
	}

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			ctl := gomock.NewController(t)
			authMockRepo := mockRepository.NewMockAuthRepository(ctl)
			tokenMockAuth := mockTokenAuth.NewMockTokenAuth(ctl)

			authUseCase := NewAuthUseCase(authMockRepo, tokenMockAuth, nil, nil, nil)

			test.buildStub(authMockRepo, tokenMockAuth)

			actualOutput, actualError := authUseCase.
				VerifyAndGetRefreshTokenSession(context.Background(), test.refreshToken, tokenUsedFor)

			assert.Equal(t, test.expectedError, actualError)
			assert.Equal(t, test.expectedOutput, actualOutput)
		})
	}
}
