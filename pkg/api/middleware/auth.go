package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/token"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type Middleware interface {
	GetUserMiddleware() gin.HandlerFunc
	GetAdminMiddleware() gin.HandlerFunc
}

type middleware struct {
	tokenAuth token.TokenAuth
}

func NewMiddleware(tokenAuth token.TokenAuth) Middleware {
	return &middleware{
		tokenAuth: tokenAuth,
	}
}

// const (
// 	authorizationHeaderKey string = "authorization"
// 	authorizationType      string = "bearer"
// )

func (c *middleware) GetUserMiddleware() gin.HandlerFunc {
	//return middleware(token.TokenForUser)
	return c.middlewareUsingCookie(token.TokenForUser)
}

func (c *middleware) GetAdminMiddleware() gin.HandlerFunc {
	//return middleware(token.TokenForAdmin)
	return c.middlewareUsingCookie(token.TokenForAdmin)
}

// func (c *middleware) middleware(tokenUser token.UserType) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {

// 		autherizationHeaderValue := ctx.GetHeader(authorizationHeaderKey)

// 		authHeaderFields := strings.Fields(autherizationHeaderValue)
// 		if len(authHeaderFields) < 2 {
// 			ctx.Abort()
// 			response := res.ErrorResponse(401, "faild to authenticate", "authentication token not provided", nil)
// 			ctx.JSON(http.StatusUnauthorized, response)
// 			return
// 		}
// 		fmt.Println(authHeaderFields)

// 		headerAuthType := authHeaderFields[0]
// 		headerAccessToken := authHeaderFields[1]

// 		if !strings.EqualFold(headerAuthType, authorizationType) {
// 			ctx.Abort()
// 			response := res.ErrorResponse(401, "faild to authenticate", "invalid authorization type", nil)
// 			ctx.JSON(http.StatusUnauthorized, response)
// 			return
// 		}

// 		payload, err := c.tokenAuth.VerifyToken(headerAccessToken, tokenUser)

// 		if err != nil {
// 			ctx.Abort()
// 			response := res.ErrorResponse(401, "faild to authenticate", err.Error(), nil)
// 			ctx.JSON(http.StatusUnauthorized, response)
// 			return
// 		}
// 		fmt.Println(payload)
// 		ctx.Set("userId", payload.UserID)
// 	}
// }

func (c *middleware) middlewareUsingCookie(tokenUser token.UserType) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		cookieName := "auth-" + string(tokenUser)
		accessToken, err := ctx.Cookie(cookieName)
		if err != nil || accessToken == "" {
			if accessToken == "" {
				err = errors.Join(err, errors.New("there is no access token"))
			}
			response := res.ErrorResponse(401, "faild to authenticate", err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		payload, err := c.tokenAuth.VerifyToken(accessToken, tokenUser)

		if err != nil {
			ctx.Abort()
			response := res.ErrorResponse(401, "faild to authenticate", err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		fmt.Println("payload", payload.UserID)
		ctx.Set("userId", payload.UserID)
	}

}
