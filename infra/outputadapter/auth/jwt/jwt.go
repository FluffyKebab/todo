package jwt

import (
	"github.com/FluffyKebab/todo/app/auth"

	"github.com/golang-jwt/jwt/v5"
)

type Authenticator struct {
	secretKey string
}

var _ auth.Authenticator = Authenticator{}

func NewAuthenticator(secretKey string) Authenticator {
	return Authenticator{
		secretKey: secretKey,
	}
}

func (a Authenticator) CreateToken(userID string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userID": userID,
		},
	)

	return t.SignedString([]byte(a.secretKey))
}

func (a Authenticator) GetUserID(token string) (string, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.NewParser().ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.secretKey), nil
	})
	if err != nil {
		return "", err
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return "", auth.ErrMissingUserIDInToken
	}

	return userID, nil
}

func (a Authenticator) HasAccess(token string, userID string) (bool, error) {
	tokenID, err := a.GetUserID(token)
	if err != nil {
		return false, err
	}

	return tokenID == userID, nil
}
