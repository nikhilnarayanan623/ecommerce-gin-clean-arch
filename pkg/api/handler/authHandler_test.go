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
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/mock/mockUseCase"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/token"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
	"github.com/stretchr/testify/assert"
)

func TestUserLogin(t *testing.T) {

	tests := map[string]struct {
		loginDetails req.Login
		buildStub    func(useCaseMock *mockUseCase.MockAuthUseCase, loginDetails req.Login)
		checkReponse func(t *testing.T, responseRecorder *httptest.ResponseRecorder)
	}{
		"EmptyLoginDetailsShouldReturnBadRequestError": {
			loginDetails: req.Login{},
			buildStub: func(useCaseMock *mockUseCase.MockAuthUseCase, loginDetails req.Login) {
				// not expecting any call to useCase
			},
			checkReponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
			},
		},
		"InvalidEmailShouldReturnErrorOnBinding": {
			loginDetails: req.Login{Email: "invalidEmail666", Password: "password"},
			buildStub: func(useCaseMock *mockUseCase.MockAuthUseCase, loginDetails req.Login) {
				//not expecting any calls
			},
			checkReponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
			},
		},
		"ValidLoginDetailsAndUserNotExist": {
			loginDetails: req.Login{Email: "validNonExistEmail@gmail.com", Password: "validPassword"},
			buildStub: func(useCaseMock *mockUseCase.MockAuthUseCase, loginDetails req.Login) {
				useCaseMock.EXPECT().UserLogin(gomock.Any(), loginDetails).
					Times(1).Return(uint(0), errors.New("user with this details not exist"))
			},
			checkReponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

			},
		},

		"FaildToGenerateAccessTokenShouldReturnInternalServerError": {
			loginDetails: req.Login{UserName: "userName", Password: "password"},
			buildStub: func(useCaseMock *mockUseCase.MockAuthUseCase, loginDetails req.Login) {
				useCaseMock.EXPECT().UserLogin(gomock.Any(), loginDetails).
					Times(1).Return(uint(1), nil)
				useCaseMock.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).
					Times(1).Return("", errors.New("faild to generate access token"))
			},
			checkReponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
			},
		},
		"FaildToGenerateRefreshTokenShouldReturnInternalServerError": {
			loginDetails: req.Login{Phone: "9078659867", Password: "password"},
			buildStub: func(useCaseMock *mockUseCase.MockAuthUseCase, loginDetails req.Login) {

				useCaseMock.EXPECT().UserLogin(gomock.Any(), loginDetails).
					Times(1).Return(uint(1), nil)
				useCaseMock.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).
					Times(1).Return("accessToken", nil)
				useCaseMock.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any()).
					Times(1).Return("", errors.New("faild to generate access_token"))
			},
			checkReponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
			},
		},
		"SuccessfullLoginShouldSetTokenOnHeaderAndResponse": {
			loginDetails: req.Login{UserName: "userName", Password: "password"},
			buildStub: func(useCaseMock *mockUseCase.MockAuthUseCase, loginDetails req.Login) {
				useCaseMock.EXPECT().UserLogin(gomock.Any(), loginDetails).
					Times(1).Return(uint(1), nil)
				useCaseMock.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).
					Times(1).Return("accessTokenFromGenerateAccessToken", nil)
				useCaseMock.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any()).
					Times(1).Return("refreshTokenFromGenerateRefreshToken", nil)
			},

			checkReponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, responseRecorder.Code)
				//assert.NotEmpty(t, responseRecorder.Header().Get(authorizationType))
				expectedOutput := res.TokenResponse{
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
		t.Run(testName, func(t *testing.T) {
			test := test
			t.Parallel()
			ctl := gomock.NewController(t)
			mockUseCase := mockUseCase.NewMockAuthUseCase(ctl)
			test.buildStub(mockUseCase, test.loginDetails)

			authHandler := NewAuthHandler(mockUseCase)
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

			test.checkReponse(t, responseRecorder)
		})
	}
}

func TestUserRenewRefreshToken(t *testing.T) {
	tokenUsedFor := token.TokenForUser
	tests := []struct {
		testName     string
		refreshToken string
		buildStub    func(mockAuthUseCase *mockUseCase.MockAuthUseCase)
		checkReponse func(t *testing.T, responseRecorder *httptest.ResponseRecorder)
	}{
		{
			testName:     "NoRefreshTokenShouldReturnBadRequest",
			refreshToken: "",
			buildStub: func(mockAuthUseCase *mockUseCase.MockAuthUseCase) {
				// no call expecting to usecase
			},
			checkReponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
			},
		},
		{
			testName:     "InvalidTokenShouldReturnBadRequest",
			refreshToken: "invalidRefreshToken",
			buildStub: func(mockAuthUseCase *mockUseCase.MockAuthUseCase) {
				mockAuthUseCase.EXPECT().VerifyAndGetRefreshTokenSession(gomock.Any(), "invalidRefreshToken", tokenUsedFor).
					Times(1).Return(domain.RefreshSession{}, errors.New("invalid refresh token"))
			},
			checkReponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
			},
		},
		{
			testName:     "FaildToGenerateAccessTokenShouldReturnInternalServerError",
			refreshToken: "validRefreshToken",
			buildStub: func(mockAuthUseCase *mockUseCase.MockAuthUseCase) {
				mockAuthUseCase.EXPECT().VerifyAndGetRefreshTokenSession(gomock.Any(), "validRefreshToken", tokenUsedFor).
					Times(1).Return(domain.RefreshSession{}, nil)

				mockAuthUseCase.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).
					Times(1).Return("", errors.New("faild to generate access token for refresh token"))
			},
			checkReponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
			},
		},
		{
			testName:     "SuccessFullRenewAccessTokenShouldReturnAccessTokenWithStatusOk",
			refreshToken: "validRefreshToken",
			buildStub: func(mockAuthUseCase *mockUseCase.MockAuthUseCase) {
				mockAuthUseCase.EXPECT().VerifyAndGetRefreshTokenSession(gomock.Any(), "validRefreshToken", tokenUsedFor).
					Times(1).Return(domain.RefreshSession{}, nil)
				mockAuthUseCase.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).
					Times(1).Return("newAccessToken", nil)
			},
			checkReponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, responseRecorder.Code)
				responseStruct, err := getResponseStructFromResponseBody(responseRecorder.Body)
				assert.NoError(t, err)

				dataFields := responseStruct.Data.([]interface{})
				tokenField := dataFields[0].(map[string]interface{})

				jsonData, err := json.Marshal(tokenField)
				assert.NoError(t, err)

				var tokenReponse res.TokenResponse
				json.Unmarshal(jsonData, &tokenReponse)

				assert.Equal(t, "newAccessToken", tokenReponse.AccessToken)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			ctl := gomock.NewController(t)
			mockAuthUseCase := mockUseCase.NewMockAuthUseCase(ctl)
			test.buildStub(mockAuthUseCase)

			authHandler := NewAuthHandler(mockAuthUseCase)

			engine := gin.New()
			url := "/renew-access-token"
			engine.POST(url, authHandler.UserRenewAccessToken())

			responseRecorder := httptest.NewRecorder()

			data := req.RefreshToken{
				RefreshToken: test.refreshToken,
			}
			byteData, err := json.Marshal(data)
			assert.NoError(t, err)
			requestBody := bytes.NewBuffer(byteData)

			mockRequest, err := http.NewRequest(http.MethodPost, url, requestBody)
			assert.NoError(t, err)

			engine.ServeHTTP(responseRecorder, mockRequest)
			test.checkReponse(t, responseRecorder)
		})
	}
}

func getResponseStructFromResponseBody(responseBody *bytes.Buffer) (responseStruct res.Response, err error) {
	data, err := io.ReadAll(responseBody)
	json.Unmarshal(data, &responseStruct)
	return
}
