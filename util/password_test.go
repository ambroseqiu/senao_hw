package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	randomPassword := RandomPassword(8)
	hashedPassword, err := HashedPassword(randomPassword)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = CheckPassword(randomPassword, hashedPassword)
	require.NoError(t, err)

	wrongPassword := RandomPassword(10)
	err = CheckPassword(wrongPassword, hashedPassword)
	require.EqualError(t, err, ErrMismatchedPassword.Error())
}

func TestPasswordDifferentHashedPassword(t *testing.T) {
	password := RandomPassword(8)
	hashedPassword1, err := HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)
	hashedPassword2, err := HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
