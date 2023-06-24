package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/token"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

type Middleware interface {
	GetUserAuthMiddleware() gin.HandlerFunc
	GetAdminAuthMiddleware() gin.HandlerFunc
}

type middleware struct {
	tokenService token.TokenService
}

func NewMiddleware(tokenService token.TokenService) Middleware {
	return &middleware{
		tokenService: tokenService,
	}
}

const (
	authorizationHeaderKey string = "authorization"
	authorizationType      string = "bearer"
)

// Get User Auth middleware
func (c *middleware) GetUserAuthMiddleware() gin.HandlerFunc {
	// return c.authorize(token.User)
	return c.middlewareUsingCookie(token.User)
}

// Get Admin Auth middleware
func (c *middleware) GetAdminAuthMiddleware() gin.HandlerFunc {
	// return c.authorize(token.Admin)
	return c.middlewareUsingCookie(token.Admin)
}

// authorize request on request header using user type
// func (c *middleware) authorize(tokenUser token.UserType) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {

// 		authorizationValues := ctx.GetHeader(authorizationHeaderKey)

// 		authFields := strings.Fields(authorizationValues)
// 		if len(authFields) < 2 {

// 			err := errors.New("authorization token not provided")

// 			response.ErrorResponse(ctx, http.StatusUnauthorized, "Failed to authorize request", err, nil)
// 			ctx.Abort()
// 			return
// 		}

// 		authType := authFields[0]
// 		accessToken := authFields[1]

// 		if !strings.EqualFold(authType, authorizationType) {
// 			err := errors.New("invalid authorization type")
// 			response.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized user", err, nil)
// 			ctx.Abort()
// 			return
// 		}

// 		tokenVerifyReq := token.VerifyTokenRequest{
// 			TokenString: accessToken,
// 			UsedFor:     tokenUser,
// 		}

// 		verifyRes, err := c.tokenService.VerifyToken(tokenVerifyReq)

// 		if err != nil {
// 			response.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized user", err, nil)
// 			ctx.Abort()
// 			return
// 		}

// 		ctx.Set("userId", verifyRes.UserID)
// 	}
// }

func (c *middleware) middlewareUsingCookie(tokenUser token.UserType) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		cookieName := "auth-" + string(tokenUser)
		accessToken, err := ctx.Cookie(cookieName)

		if err != nil || accessToken == "" {
			if accessToken == "" {
				err = errors.Join(err, errors.New("access token not found"))
			}
			response.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized user", err, nil)
			ctx.Abort()
			return
		}
		verifyRes, err := c.tokenService.VerifyToken(token.VerifyTokenRequest{
			TokenString: accessToken,
			UsedFor:     tokenUser,
		})

		if err != nil {
			response.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized user", err, nil)
			ctx.Abort()
			return
		}

		ctx.Set("userId", verifyRes.UserID)
	}

}
