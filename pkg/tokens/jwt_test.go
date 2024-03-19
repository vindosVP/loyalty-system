package tokens

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestExtractID(t *testing.T) {
	userID := 1
	userLogin := "someLogin"
	jwtSecret := "superSecret"

	token, err := CreateJWT(JWTClaims(userID, userLogin, time.Now().Add(time.Hour*72).Unix()), jwtSecret)
	assert.NoError(t, err)

	strID, err := ExtractID(token, jwtSecret)
	assert.NoError(t, err)

	id, err := strconv.Atoi(strID)
	assert.NoError(t, err)
	assert.Equal(t, userID, id)
}
