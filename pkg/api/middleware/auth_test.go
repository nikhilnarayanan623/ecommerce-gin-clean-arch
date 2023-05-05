package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMiddlware(t *testing.T) {

	tests := []struct {
		testName      string
		buildStub     func(mockRequest *http.Request)
		checkResponse func(t *testing.T, responseRecorder *httptest.ResponseRecorder)
	}{
		// {
		// 	testName: "NoAuthorizationTypeOnHeaderShouldReturnUnauthorized",
		// 	buildStub: func(mockRequest *http.Request) {
		// 		mo
		// 	},
		// },
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			responseRecorder := httptest.NewRecorder()

			url := "/for/auth"
			mockRequest, err := http.NewRequest(http.MethodConnect, url, nil)
			assert.NoError(t, err)
			test.buildStub(mockRequest)

			server := gin.New()
			server.Any(url)

			server.ServeHTTP(responseRecorder, mockRequest)

			test.checkResponse(t, responseRecorder)
		})
	}
}
