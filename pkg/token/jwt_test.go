package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {

	tests := []struct {
		name            string
		input           GenerateTokenRequest
		expectingOutput bool
		expectedError   error
	}{

		{
			name:            "InvalidUserTypeShouldReturnError",
			input:           GenerateTokenRequest{UserID: 1, UsedFor: "Invalid user type"},
			expectingOutput: false,
			expectedError:   ErrInvalidUserType,
		},
		{
			name:            "ValidPayloadShouldReturnTokenString",
			input:           GenerateTokenRequest{UserID: 12, UsedFor: User},
			expectingOutput: true,
			expectedError:   nil,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			cfg := config.Config{AdminAuthKey: "adminSecret", UserAuthKey: "userSecret"}
			tokenAuth := NewTokenService(cfg)
			actualOutput, actualError := tokenAuth.GenerateToken(test.input)

			assert.Equal(t, test.expectedError, actualError)
			if test.expectingOutput {
				assert.NotEmpty(t, actualOutput)
			} else {
				assert.Empty(t, actualOutput)
			}
		})
	}

}

func TestVerifyToken(t *testing.T) {

	tests := []struct {
		name           string
		tokenUser      UserType
		expectedOutput VerifyTokenResponse
		buildStub      func(t *testing.T, tokenAuth TokenService) (tokenString string)
		expectedError  error
	}{
		{
			name:           "EmptyTokenStringShouldReturnError",
			tokenUser:      Admin,
			expectedOutput: VerifyTokenResponse{},
			buildStub:      func(t *testing.T, tokenAuth TokenService) (tokenString string) { return },
			expectedError:  ErrInvalidToken,
		},
		{
			name:           "ExpiredTokenShouldReturnExpiredError",
			tokenUser:      User,
			expectedOutput: VerifyTokenResponse{},
			buildStub: func(t *testing.T, tokenAuth TokenService) string {
				response, err := tokenAuth.GenerateToken(GenerateTokenRequest{
					UserID:   12,
					ExpireAt: time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC),
					UsedFor:  User,
				})
				assert.NoError(t, err)
				return response.TokenString
			},
			expectedError: ErrExpiredToken,
		},
		{
			name:           "ChangedSigningMethodShouldReturnInvalidTokenError",
			expectedOutput: VerifyTokenResponse{},
			tokenUser:      User,
			buildStub: func(t *testing.T, tokenAuth TokenService) (tokenString string) {
				// create a token with unsafe signature
				token := jwt.NewWithClaims(jwt.SigningMethodNone, &jwtClaims{})
				tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
				assert.NoError(t, err)

				return
			},
			expectedError: ErrInvalidToken,
		},
		{
			name:           "InvalidUserTypeShouldReturnError",
			tokenUser:      "Invalid User Type",
			expectedOutput: VerifyTokenResponse{},
			buildStub: func(t *testing.T, tokenAuth TokenService) (tokenString string) {
				return
			},
			expectedError: ErrInvalidUserType,
		},
		{
			name:           "ValidTokenShouldReturnResponse",
			tokenUser:      Admin,
			expectedOutput: VerifyTokenResponse{UserID: 12, TokenID: "token_id"},
			buildStub: func(t *testing.T, tokenAuth TokenService) string {
				request := GenerateTokenRequest{
					UserID:   12,
					UsedFor:  Admin,
					ExpireAt: time.Now().Add(time.Hour * 1),
				}
				response, err := tokenAuth.GenerateToken(request)
				assert.NoError(t, err)
				return response.TokenString
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			cfg := config.Config{AdminAuthKey: "adminSecret", UserAuthKey: "userSecret"}
			tokenAuth := NewTokenService(cfg)

			tokenString := test.buildStub(t, tokenAuth)

			verifyRequest := VerifyTokenRequest{
				TokenString: tokenString,
				UsedFor:     test.tokenUser,
			}

			actualOutput, actualError := tokenAuth.VerifyToken(verifyRequest)

			assert.Equal(t, test.expectedError, actualError)

			assert.Equal(t, test.expectedOutput.UserID, actualOutput.UserID)

		})
	}
}
