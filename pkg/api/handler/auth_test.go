package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/mock/mockusecase"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/service/token"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase"
	"github.com/stretchr/testify/assert"
)

func TestUserLogin(t *testing.T) {

	tests := map[string]struct {
		loginDetails  request.Login
		buildStub     func(useCaseMock *mockusecase.MockAuthUseCase, loginDetails request.Login)
		checkResponse func(t *testing.T, responseRecorder *httptest.ResponseRecorder)
	}{
		"EmptyLoginDetailsShouldReturnBadRequestError": {
			loginDetails: request.Login{},
			buildStub: func(useCaseMock *mockusecase.MockAuthUseCase, loginDetails request.Login) {
				// not expecting any call to useCase
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
			},
		},
		"InvalidEmailShouldReturnErrorOnBinding": {
			loginDetails: request.Login{Email: "invalidEmail666", Password: "password"},
			buildStub: func(useCaseMock *mockusecase.MockAuthUseCase, loginDetails request.Login) {
				//not expecting any calls
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
			},
		},
		"ValidLoginDetailsAndUserNotExist": {
			loginDetails: request.Login{Email: "validNonExistEmail@gmail.com", Password: "validPassword"},

			buildStub: func(useCaseMock *mockusecase.MockAuthUseCase, loginDetails request.Login) {
				useCaseMock.EXPECT().UserLogin(gomock.Any(), loginDetails).
					Times(1).Return(uint(0), usecase.ErrUserNotExist)
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
			},
		},

		"FailedToGenerateAccessTokenShouldReturnInternalServerError": {
			loginDetails: request.Login{UserName: "userName", Password: "password"},
			buildStub: func(useCaseMock *mockusecase.MockAuthUseCase, loginDetails request.Login) {
				useCaseMock.EXPECT().UserLogin(gomock.Any(), loginDetails).
					Times(1).Return(uint(1), nil)
				useCaseMock.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).
					Times(1).Return("", errors.New("failed to generate access token"))
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
			},
		},
		"FailedToGenerateRefreshTokenShouldReturnInternalServerError": {
			loginDetails: request.Login{Phone: "9078659867", Password: "password"},
			buildStub: func(useCaseMock *mockusecase.MockAuthUseCase, loginDetails request.Login) {

				useCaseMock.EXPECT().UserLogin(gomock.Any(), loginDetails).
					Times(1).Return(uint(1), nil)
				useCaseMock.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).
					Times(1).Return("accessToken", nil)
				useCaseMock.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any()).
					Times(1).Return("", errors.New("faild to generate access_token"))
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
			},
		},
		"SuccessfulLoginShouldSetTokenOnHeaderAndResponse": {
			loginDetails: request.Login{UserName: "userName", Password: "password"},
			buildStub: func(useCaseMock *mockusecase.MockAuthUseCase, loginDetails request.Login) {
				useCaseMock.EXPECT().UserLogin(gomock.Any(), loginDetails).
					Times(1).Return(uint(1), nil)
				useCaseMock.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).
					Times(1).Return("accessTokenFromGenerateAccessToken", nil)
				useCaseMock.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any()).
					Times(1).Return("refreshTokenFromGenerateRefreshToken", nil)
			},

			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, responseRecorder.Code)
				//assert.NotEmpty(t, responseRecorder.Header().Get(authorizationType))
				expectedOutput := response.TokenResponse{
					AccessToken:  "accessTokenFromGenerateAccessToken",
					RefreshToken: "refreshTokenFromGenerateRefreshToken",
				}

				responseStruct, err := getResponseStructFromResponseBody(responseRecorder.Body)
				assert.NoError(t, err)
				dataFields := responseStruct.Data.([]interface{})
				tokenData := dataFields[0].(map[string]interface{})

				assert.Equal(t, expectedOutput.AccessToken, tokenData["access_token"])
				assert.Equal(t, expectedOutput.RefreshToken, tokenData["refresh_token"])
			},
		},
	}

	for testName, test := range tests {
		test := test
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			ctl := gomock.NewController(t)
			mockUseCase := mockusecase.NewMockAuthUseCase(ctl)
			test.buildStub(mockUseCase, test.loginDetails)

			authHandler := NewAuthHandler(mockUseCase, config.Config{})
			server := gin.New()
			url := "/login"
			server.POST(url, authHandler.UserLogin)

			jsonData, err := json.Marshal(test.loginDetails)
			assert.NoError(t, err)
			body := bytes.NewBuffer(jsonData)

			mockRequest, err := http.NewRequest(http.MethodPost, url, body)
			assert.NoError(t, err)
			responseRecorder := httptest.NewRecorder()

			server.ServeHTTP(responseRecorder, mockRequest)

			test.checkResponse(t, responseRecorder)
		})
	}
}

func TestUserRenewRefreshToken(t *testing.T) {
	tokenUsedFor := token.User
	tests := []struct {
		testName      string
		refreshToken  string
		buildStub     func(mockAuthUseCase *mockusecase.MockAuthUseCase)
		checkResponse func(t *testing.T, responseRecorder *httptest.ResponseRecorder)
	}{
		{
			testName:     "NoRefreshTokenShouldReturnBadRequest",
			refreshToken: "",
			buildStub: func(mockAuthUseCase *mockusecase.MockAuthUseCase) {
				// no call expecting to usecase
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
			},
		},
		{
			testName:     "InvalidTokenShouldReturnBadRequest",
			refreshToken: "invalidRefreshToken",
			buildStub: func(mockAuthUseCase *mockusecase.MockAuthUseCase) {
				mockAuthUseCase.EXPECT().VerifyAndGetRefreshTokenSession(gomock.Any(), "invalidRefreshToken", tokenUsedFor).
					Times(1).Return(domain.RefreshSession{}, usecase.ErrInvalidRefreshToken)
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
			},
		},
		{
			testName:     "FailedToGenerateAccessTokenShouldReturnInternalServerError",
			refreshToken: "validRefreshToken",
			buildStub: func(mockAuthUseCase *mockusecase.MockAuthUseCase) {
				mockAuthUseCase.EXPECT().VerifyAndGetRefreshTokenSession(gomock.Any(), "validRefreshToken", tokenUsedFor).
					Times(1).Return(domain.RefreshSession{}, nil)

				mockAuthUseCase.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).
					Times(1).Return("", errors.New("failed to generate access token for refresh token"))
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
			},
		},
		{
			testName:     "SuccessFullRenewAccessTokenShouldReturnAccessTokenWithStatusOk",
			refreshToken: "validRefreshToken",
			buildStub: func(mockAuthUseCase *mockusecase.MockAuthUseCase) {
				mockAuthUseCase.EXPECT().VerifyAndGetRefreshTokenSession(gomock.Any(), "validRefreshToken", tokenUsedFor).
					Times(1).Return(domain.RefreshSession{}, nil)
				mockAuthUseCase.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).
					Times(1).Return("generated_access_token", nil)
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, responseRecorder.Code)
				responseStruct, err := getResponseStructFromResponseBody(responseRecorder.Body)
				assert.NoError(t, err)

				dataFields := responseStruct.Data.([]interface{})
				tokenField := dataFields[0].(map[string]interface{})

				jsonData, err := json.Marshal(tokenField)
				assert.NoError(t, err)

				var tokenResponse response.TokenResponse
				err = json.Unmarshal(jsonData, &tokenResponse)
				assert.NoError(t, err)

				assert.Equal(t, "generated_access_token", tokenResponse.AccessToken)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			ctl := gomock.NewController(t)
			mockAuthUseCase := mockusecase.NewMockAuthUseCase(ctl)
			test.buildStub(mockAuthUseCase)

			authHandler := NewAuthHandler(mockAuthUseCase, config.Config{})

			engine := gin.New()
			url := "/renew-access-token"
			engine.POST(url, authHandler.UserRenewAccessToken())

			responseRecorder := httptest.NewRecorder()

			data := request.RefreshToken{
				RefreshToken: test.refreshToken,
			}
			byteData, err := json.Marshal(data)
			assert.NoError(t, err)
			requestBody := bytes.NewBuffer(byteData)

			mockRequest, err := http.NewRequest(http.MethodPost, url, requestBody)
			assert.NoError(t, err)

			engine.ServeHTTP(responseRecorder, mockRequest)
			test.checkResponse(t, responseRecorder)
		})
	}
}

func getResponseStructFromResponseBody(responseBody *bytes.Buffer) (responseStruct response.Response, err error) {

	data, err := io.ReadAll(responseBody)
	if err := json.Unmarshal(data, &responseStruct); err != nil {
		return response.Response{}, nil
	}

	return
}
