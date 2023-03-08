package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/auth"
)

func AuthenticateUser(ctx *gin.Context) {

	tokenString, _ := ctx.Cookie("user-auth")

	token, err := auth.ValidateToken(tokenString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"StatusCode": 401,
			"msg":        "not an autherized user",
		})
		return
	}

	// if its a valid token then claim it
	claims, ok := token.Claims.(*jwt.StandardClaims)

	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"StatusCode": 401,
			"msg":        "can't claim the token",
		})
		return
	}

	// claim the userId and set it on context
	ctx.Set("userId", claims.Id)
}

func AuthenticateAdmin(ctx *gin.Context) {
	tokenString, _ := ctx.Cookie("admin-auth")

	_, err := auth.ValidateToken(tokenString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"StatusCode": 401,
			"msg":        "not an autherized user",
		})
		return
	}

	// // if its a valid token then claim it
	// claims, ok := token.Claims.(*jwt.StandardClaims)

	// if !ok {
	// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
	// 		"StatusCode": 401,
	// 		"msg":        "can't claim the token",
	// 	})
	// 	return
	// }
}
