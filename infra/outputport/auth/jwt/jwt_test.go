package jwt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJWT(t *testing.T) {
	t.Parallel()

	userID := "1234"
	a := NewAuthenticator("secret_key")
	token, err := a.CreateToken(userID)
	require.NoError(t, err)

	tokenUserID, err := a.GetUserID(token)
	require.NoError(t, err)
	require.Equal(t, userID, tokenUserID, "userID put in token not equal the userID in token claims")
}
