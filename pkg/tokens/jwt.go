package tokens

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
)

func JWTClaims(id int, login string, exp int64) jwt.MapClaims {
	return jwt.MapClaims{
		"id":    id,
		"login": login,
		"exp":   exp,
	}
}

func CreateJWT(claims jwt.MapClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("token.SignedString: %w", err)
	}
	return tokenString, nil
}

func IsAuthorized(token string, secret string) (bool, error) {
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractID(requestToken string, secret string) (string, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	id := claims["id"].(float64)
	return strconv.FormatFloat(id, 'f', 0, 64), nil
}
