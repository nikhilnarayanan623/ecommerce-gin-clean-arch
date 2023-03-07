package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
)

func Authentication(ctx *gin.Context) {

	token, _ := ctx.Cookie("jwt-auth")

	if err := validteToken(token); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"StatusCode": 401,
			"msg":        "Unauthorized user",
		})
	}
}

func validteToken(token string) error {
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.GetJWTCofig()), nil
	})

	return err
}
