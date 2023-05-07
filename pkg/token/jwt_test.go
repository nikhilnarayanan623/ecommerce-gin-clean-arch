package token

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {

	tests := []struct {
		testName               string
		inputPayload           *Payload
		inputUserType          UserType
		IsExpectingTokenString bool
		expectedError          error
	}{
		{
			testName:               "NilPayloadInputShouldReturnError",
			inputPayload:           nil,
			inputUserType:          TokenForAdmin,
			IsExpectingTokenString: false,
			expectedError:          errors.New("payload should not be nil"),
		},
		{
			testName:               "InvalidUserTypeShouldReturnError",
			inputPayload:           &Payload{TokenID: uuid.New(), UserID: 2},
			inputUserType:          "invalidUseType",
			IsExpectingTokenString: false,
			expectedError:          errors.New("invalid user_type"),
		},
		{
			testName:               "ValidPayloadShouldReturnTokenString",
			inputUserType:          TokenForUser,
			inputPayload:           &Payload{TokenID: uuid.New(), UserID: 1},
			IsExpectingTokenString: true,
			expectedError:          nil,
		},
	}

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {

			cfg := config.Config{JWTAdmin: "adminSecret", JWTUser: "userSecret"}
			tokenAuth := NewJWTAuth(cfg)
			tokenString, actualError := tokenAuth.CreateToken(test.inputPayload, test.inputUserType)

			assert.Equal(t, test.expectedError, actualError)
			if test.IsExpectingTokenString {
				assert.NotEmpty(t, tokenString)
			} else {
				assert.Empty(t, tokenString)
			}
		})
	}

}

func TestVerifyToken(t *testing.T) {

	tests := []struct {
		testName       string
		userType       UserType
		expectedOutput Payload
		buildStub      func(t *testing.T, tokenAuth TokenAuth) (tokenString string)
		expectedError  error
	}{
		{
			testName:       "EmptyTokenStringShouldRetunError",
			expectedOutput: Payload{},
			userType:       TokenForUser,
			buildStub:      func(t *testing.T, tokenAuth TokenAuth) (tokenString string) { return },
			expectedError:  errInvalidToken,
		},
		{
			testName:       "ExpiredTokenShouldReturnExpiredError",
			userType:       TokenForUser,
			expectedOutput: Payload{},
			buildStub: func(t *testing.T, tokenAuth TokenAuth) (tokenString string) {
				tokenString, err := tokenAuth.CreateToken(&Payload{
					TokenID:  uuid.New(),
					UserID:   12,
					ExpireAt: time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC),
				}, TokenForUser)
				assert.NoError(t, err)
				return
			},
			expectedError: errExpiredToken,
		},
		{
			testName:       "ChangedSigninMethodShouldReturnInvalidTokenError",
			expectedOutput: Payload{},
			userType:       TokenForAdmin,
			buildStub: func(t *testing.T, tokenAuth TokenAuth) (tokenString string) {
				token := jwt.NewWithClaims(jwt.SigningMethodNone, &Payload{})

				tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
				assert.NoError(t, err)

				return
			},
			expectedError: errInvalidToken,
		},
		{
			testName:       "InvalidUserTypeShouldReturnEror",
			userType:       "inalidUserType",
			expectedOutput: Payload{},
			buildStub: func(t *testing.T, tokenAuth TokenAuth) (tokenString string) {
				return
			},
			expectedError: errors.New("invalid user_type"),
		},
		{
			testName:       "ValidTokenShouldReturnPalyload",
			userType:       TokenForUser,
			expectedOutput: Payload{UserID: 1, ExpireAt: time.Now().Add(time.Hour * 1)},
			buildStub: func(t *testing.T, tokenAuth TokenAuth) (tokenString string) {
				payload := &Payload{
					UserID:   1,
					ExpireAt: time.Now().Add(time.Hour * 1),
				}
				tokenString, err := tokenAuth.CreateToken(payload, TokenForUser)
				assert.NoError(t, err)
				return
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			cfg := config.Config{JWTAdmin: "adminSecret", JWTUser: "userSecret"}
			tokenAuth := NewJWTAuth(cfg)

			tokenString := test.buildStub(t, tokenAuth)

			payload, actualError := tokenAuth.VerifyToken(tokenString, test.userType)

			assert.Equal(t, test.expectedError, actualError)

			assert.Equal(t, test.expectedOutput.UserID, payload.UserID)
			assert.Equal(t, test.expectedOutput.ExpireAt.Day(), payload.ExpireAt.Day())

		})
	}
}
