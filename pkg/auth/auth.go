package auth

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	errNoAuthHeader      = fmt.Errorf("no auth header")
	errInvalidAuthFormat = fmt.Errorf("invalid auth format")
)

const bearerSchema = "Bearer "

func ParseBearerToken(r *http.Request) (string, error) {
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		return "", errNoAuthHeader
	}
	splitToken := strings.Split(reqToken, bearerSchema)
	if len(splitToken) != 2 {
		return "", errInvalidAuthFormat
	}
	return splitToken[1], nil
}
