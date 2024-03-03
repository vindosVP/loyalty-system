package tokens

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

func JWTClaims(id int, login string, exp int64) jwt.MapClaims {
	return jwt.MapClaims{
		"id":    id,
		"login": login,
		"exp":   exp,
	}
}

func CreateJwt(claims jwt.MapClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("token.SignedString: %w", err)
	}
	return tokenString, nil
}
