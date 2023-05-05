package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/token"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

var tokenAuth token.TokenAuth

func SetupMiddleware(cfg config.Config) {
	tokenAuth = token.NewJWTAuth(cfg.JWTAdmin, cfg.JWTUser)
}

// const (
// 	authorizationHeaderKey string = "authorization"
// 	authorizationType      string = "bearer"
// )

func GetUserMiddleware() gin.HandlerFunc {
	//return middleware(token.TokenForUser)
	return middlewareUsingCookie(token.TokenForUser)
}

func GetAdminMiddleware() gin.HandlerFunc {
	//return middleware(token.TokenForAdmin)
	return middlewareUsingCookie(token.TokenForAdmin)
}

// func middleware(tokenUser token.UserType) gin.HandlerFunc {
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

// 		payload, err := tokenAuth.VerifyToken(headerAccessToken, tokenUser)

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

func middlewareUsingCookie(tokenUser token.UserType) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		accessToken, err := ctx.Cookie("user-auth")
		if err != nil || accessToken == "" {
			if accessToken == "" {
				err = errors.Join(err, errors.New("there is no access token"))
			}
			response := res.ErrorResponse(401, "faild to authenticate", err.Error(), nil)
			ctx.JSON(http.StatusUnauthorized, response)
		}
		payload, err := tokenAuth.VerifyToken(accessToken, tokenUser)

		if err != nil {
			ctx.Abort()
			response := res.ErrorResponse(401, "faild to authenticate", err.Error(), nil)
			ctx.JSON(http.StatusUnauthorized, response)
			return
		}
		fmt.Println(payload)
		ctx.Set("userId", payload.UserID)
	}

}
