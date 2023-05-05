package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/mock/mockUseCase"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/token"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
	"github.com/stretchr/testify/assert"
)

func TestUserLogin(t *testing.T) {

	tests := []struct {
		testName     string
		loginDetails req.Login
		buildStub    func(useCaseMock *mockUseCase.MockAuthUseCase, loginDetails req.Login)
		checkReponse func(t *testing.T, responseRecorder *httptest.ResponseRecorder)
	}{
		{
			testName:     "EmptyLoginDetailsShouldReturnBadRequestError",
			loginDetails: req.Login{},
			buildStub: func(useCaseMock *mockUseCase.MockAuthUseCase, loginDetails req.Login) {
				// not expecting any call to useCase
			},
			checkReponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
			},
		},
		{
			testName:     "InvalidEmailShouldReturnErrorOnBinding",
			loginDetails: req.Login{Email: "invalidEmail666", Password: "password"},
			buildStub: func(useCaseMock *mockUseCase.MockAuthUseCase, loginDetails req.Login) {
				//not expecting any calls
			},
			checkReponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
			},
		},
		{
			testName:     "ValidLoginDetailsAndUserNotExist",
			loginDetails: req.Login{Email: "validNonExistEmail@gmail.com", Password: "validPassword"},
			buildStub: func(useCaseMock *mockUseCase.MockAuthUseCase, loginDetails req.Login) {
				useCaseMock.EXPECT().UserLogin(gomock.Any(), loginDetails).
					Times(1).Return(uint(0), errors.New("user with this details not exist"))
			},
			checkReponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

			},
		},

		{
			testName:     "FaildToGenerateAccessTokenShouldReturnInternalServerError",
			loginDetails: req.Login{UserName: "userName", Password: "password"},
			buildStub: func(useCaseMock *mockUseCase.MockAuthUseCase, loginDetails req.Login) {
				useCaseMock.EXPECT().UserLogin(gomock.Any(), loginDetails).
					Times(1).Return(uint(1), nil)

				useCaseMock.EXPECT().GenerateAccessToken(gomock.Any(), uint(1), token.TokenForUser, time.Minute*15).
					Times(1).Return("", errors.New("faild to generate access token"))
			},
			checkReponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
			},
		},
		{
			testName:     "FaildToGenerateRefreshTokenShouldReturnInternalServerError",
			loginDetails: req.Login{Phone: "9078659867", Password: "password"},
			buildStub: func(useCaseMock *mockUseCase.MockAuthUseCase, loginDetails req.Login) {
				useCaseMock.EXPECT().UserLogin(gomock.Any(), loginDetails).
					Times(1).Return(uint(1), nil)
				useCaseMock.EXPECT().GenerateAccessToken(gomock.Any(), uint(1), token.TokenForUser, time.Minute*15).
					Times(1).Return("accessToken", nil)
				useCaseMock.EXPECT().GenerateRefreshToken(gomock.Any(), uint(1), token.TokenForUser, time.Hour*24*7).
					Times(1).Return("", errors.New("faild to generate access_token"))
			},
			checkReponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
			},
		},
		{
			testName:     "SuccessfullLoginShouldSetTokenOnHeaderAndResponse",
			loginDetails: req.Login{UserName: "userName", Password: "password"},
			buildStub: func(useCaseMock *mockUseCase.MockAuthUseCase, loginDetails req.Login) {
				useCaseMock.EXPECT().UserLogin(gomock.Any(), loginDetails).
					Times(1).Return(uint(1), nil)
				useCaseMock.EXPECT().GenerateAccessToken(gomock.Any(), uint(1), token.TokenForUser, time.Minute*15).
					Times(1).Return("accessTokenFromGenerateAccessToken", nil)
				useCaseMock.EXPECT().GenerateRefreshToken(gomock.Any(), uint(1), token.TokenForUser, time.Hour*24*7).
					Times(1).Return("refreshTokenFromGenerateRefreshToken", nil)
			},

			checkReponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, responseRecorder.Code)
				//assert.NotEmpty(t, responseRecorder.Header().Get(authorizationType))
				expectedOutput := res.TokenResponse{
					AccessToken:  "accessTokenFromGenerateAccessToken",
					RefreshToken: "refreshTokenFromGenerateRefreshToken",
				}

				var responseStruct, err = getResponseStructFromResponseBody(responseRecorder.Body)
				assert.NoError(t, err)
				dataFields := responseStruct.Data.([]interface{})
				tokenData := dataFields[0].(map[string]interface{})

				assert.Equal(t, expectedOutput.AccessToken, tokenData["access_token"])
				assert.Equal(t, expectedOutput.RefreshToken, tokenData["refresh_token"])
			},
		},
	}

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {

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

func getResponseStructFromResponseBody(responseBody *bytes.Buffer) (responseStruct res.Response, err error) {
	data, err := io.ReadAll(responseBody)
	json.Unmarshal(data, &responseStruct)
	return
}
