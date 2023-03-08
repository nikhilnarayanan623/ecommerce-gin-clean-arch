package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
)

// type jwtClaims struct {
// 	ID uint
// 	jwt.StandardClaims
// }

func GenerateJWT(id uint) (map[string]string, error) {
	expireTime := time.Now().Add(10 * time.Minute).Unix()

	// claims := &jwtClaims{
	// 	ID:             id,
	// 	StandardClaims: jwt.StandardClaims{

	// 	},
	// }
	// create token with expire time and claims id as user id
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: expireTime,
		Id:        fmt.Sprint(id),
	})

	// conver the token into signed string
	tokenString, err := token.SignedString([]byte(config.GetJWTCofig()))

	if err != nil {
		return nil, err
	}
	// refresh token add next time
	return map[string]string{"accessToken": tokenString}, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {

	return jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(config.GetJWTCofig()), nil
		},
	)
	// if err != nil || !token.Valid {
	// 	return nil, errors.New("not valid token")
	// }
}
