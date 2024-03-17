package auth

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrNoAuthHeader      = fmt.Errorf("no auth header")
	ErrInvalidAuthFormat = fmt.Errorf("invalid auth format")
)

const bearerSchema = "Bearer "

func ParseBearerToken(r *http.Request) (string, error) {
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		return "", ErrNoAuthHeader
	}
	splitToken := strings.Split(reqToken, bearerSchema)
	if len(splitToken) != 2 {
		return "", ErrInvalidAuthFormat
	}
	return splitToken[1], nil
}
