package passwords

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCompare(t *testing.T) {
	pwd := "somePassword"

	encrypted, err := Encrypt(pwd)
	require.NoError(t, err)
	require.NotEmpty(t, encrypted)

	valid := Compare(pwd, encrypted)
	assert.True(t, valid)
}
