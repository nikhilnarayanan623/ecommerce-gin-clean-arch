package middleware

// import (
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang/mock/gomock"
// 	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/token"
// 	"github.com/stretchr/testify/assert"
// )

// func TestGetUserMiddleware(t *testing.T) {

// 	tests := []struct {
// 		testName         string
// 		buildStubRequest func(t *testing.T, mockRequest *http.Request, tokenAuth *mockTokenAuth.MockTokenAuth)
// 		checkResponse    func(t *testing.T, responseRecorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			testName: "EmptyCookieWithCookieNameShoulReturnUnAuthorized",
// 			buildStubRequest: func(t *testing.T, mockRequest *http.Request, tokenAuth *mockTokenAuth.MockTokenAuth) {
// 				cokkieName := "auth-" + string(token.TokenForUser)
// 				cookie := &http.Cookie{Name: cokkieName, Value: ""}
// 				mockRequest.AddCookie(cookie)
// 			},
// 			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
// 				assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
// 			},
// 		},
// 		{
// 			testName: "InvalidTokenShouldReturnUnAuthorized",
// 			buildStubRequest: func(t *testing.T, mockRequest *http.Request, tokenAuth *mockTokenAuth.MockTokenAuth) {
// 				cokkieName := "auth-" + string(token.TokenForUser)
// 				accessToken := "invalidAccessToken"
// 				cookie := &http.Cookie{Name: cokkieName, Value: accessToken}
// 				mockRequest.AddCookie(cookie)

// 				tokenAuth.EXPECT().VerifyToken(accessToken, token.TokenForUser).
// 					Times(1).Return(token.Payload{}, errors.New("invalid Token"))
// 			},
// 			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
// 				assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
// 			},
// 		},
// 		{
// 			testName: "ValidTokenShould",
// 			buildStubRequest: func(t *testing.T, mockRequest *http.Request, tokenAuth *mockTokenAuth.MockTokenAuth) {
// 				cokkieName := "auth-" + string(token.TokenForUser)
// 				accessToken := "validAccessToken"
// 				cookie := &http.Cookie{Name: cokkieName, Value: accessToken}
// 				mockRequest.AddCookie(cookie)

// 				tokenAuth.EXPECT().VerifyToken(accessToken, token.TokenForUser).
// 					Times(1).Return(token.Payload{UserID: 2}, nil)
// 			},
// 			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
// 				assert.Equal(t, http.StatusOK, responseRecorder.Code)
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.testName, func(t *testing.T) {

// 			ctl := gomock.NewController(t)
// 			tokenAuth := mockTokenAuth.NewMockTokenAuth(ctl)
// 			mockMiddleware := NewMiddleware(tokenAuth)
// 			responseRecorder := httptest.NewRecorder()

// 			url := "/for/auth"
// 			mockRequest, err := http.NewRequest(http.MethodGet, url, nil)
// 			assert.NoError(t, err)
// 			test.buildStubRequest(t, mockRequest, tokenAuth)

// 			server := gin.New()
// 			server.Any(url, mockMiddleware.GetUserMiddleware())

// 			server.ServeHTTP(responseRecorder, mockRequest)

// 			test.checkResponse(t, responseRecorder)
// 		})
// 	}
// }
