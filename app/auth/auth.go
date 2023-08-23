package auth

import "errors"

var ErrMissingUserIDInToken = errors.New("missing userID in claims of token")

type Authenticator interface {
	CreateToken(userID string) (string, error)
	GetUserID(token string) (string, error)
	HasAccess(token string, userID string) (bool, error)
}
