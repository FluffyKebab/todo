package testauth

import (
	"github.com/FluffyKebab/todo/app/auth"
)

type Authenticator struct {
	HasAccessReturn bool
}

var _ auth.Authenticator = Authenticator{}

func (a Authenticator) CreateToken(userID string) (string, error) {
	return userID, nil
}

func (a Authenticator) GetUserID(token string) (string, error) {
	return token, nil
}

func (a Authenticator) HasAccess(token string, userID string) (bool, error) {
	return a.HasAccessReturn, nil
}
