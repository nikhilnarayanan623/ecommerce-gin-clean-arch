package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/mock/mockRepository"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/token"
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
			expectedError: fmt.Errorf("user not exist with given lgoin details"),
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
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			userMockRepo := mockRepository.NewMockUserRepository(ctl)
			test.buildStub(userMockRepo, test.input)

			authUseCase := NewAuthUseCase(userMockRepo, nil, nil, config.Config{})
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
	type inputFields struct {
		userID       uint
		userType     token.UserType
		timeDuration time.Duration
	}
	tests := []struct {
		testName    string
		inputField  inputFields
		buildStub   func(t *testing.T, mockRepo *mockRepository.MockAuthRepository)
		checkOutput func(t *testing.T, tokenString string, err error)
	}{
		{
			testName:   "SaveRefreshTokenSessionError",
			inputField: inputFields{userID: 1, userType: "admin"},
			buildStub: func(t *testing.T, mockRepo *mockRepository.MockAuthRepository) {
				mockRepo.EXPECT().SaveRefreshSession(gomock.Any(), gomock.Any()).
					Times(1).Return(errors.New("faild to save refresh_session on database"))
			},
			checkOutput: func(t *testing.T, tokenString string, err error) {
				assert.EqualError(t, err, "faild to save refresh_session on database")
				assert.Empty(t, tokenString)
			},
		},
		{
			testName:   "SuccesfullRefreshTokenCreationAndSaveReturnToken",
			inputField: inputFields{userID: 1, userType: "admin"},
			buildStub: func(t *testing.T, mockRepo *mockRepository.MockAuthRepository) {
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
			mockRepo := mockRepository.NewMockAuthRepository(ctl)

			test.buildStub(t, mockRepo)

			authUseCase := NewAuthUseCase(nil, nil, mockRepo, config.Config{JWT: "secretKey"})

			tokenString, err := authUseCase.
				GenerateRefreshToken(context.Background(), test.inputField.userID, test.inputField.userType, test.inputField.timeDuration)

			test.checkOutput(t, tokenString, err)
		})
	}
}
