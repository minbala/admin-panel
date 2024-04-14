package commonimport

import (
	"context"
)

type Authenticator interface {
	Authenticate(ctx context.Context, token string) (UserData, error)
	AuthenticateAdmin(ctx context.Context, token string) (UserData, error)
}

type JwtManager interface {
	ParseAuthJwtToken(token string, secretKey string) (*AuthJwtClaims, error)
	CreateAuthToken(data AuthJwtTokenInput) (string, error)
}

type UserLogRecorder interface {
	Record(ctx context.Context, log UserLog) error
}

type PasswordManager interface {
	CheckPassword(password string, hashPassword string) error
	Hash(plain string) (string, error)
}
