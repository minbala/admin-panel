package jwt

import (
	commonImport "admin-panel/pkg/common"
	"admin-panel/pkg/config"
	"admin-panel/pkg/resources"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/golang-jwt/jwt/v5"
	"sync"
	"time"
)

type TokenParser struct {
	config *config.Config
}

var tokeParser commonImport.JwtManager
var setUpTokenParserOnce sync.Once

func ProvideTokenParser(config *config.Config) commonImport.JwtManager {
	setUpTokenParserOnce.Do(func() {
		tokeParser = &TokenParser{config}
	})
	return tokeParser
}

func GenerateJwtToken(claims jwt.Claims, secretKey string) (token string, err error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
}

// GenerateAuthJwtToken generate tokens
func (c TokenParser) generateAuthJwtToken(data commonImport.AuthJwtTokenInput, secret string, expiryTime time.Time) (string, error) {
	claims := commonImport.AuthJwtClaims{
		AuthJwtTokenInput: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
			Issuer:    "adminPanel",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := GenerateJwtToken(claims, secret)
	return token, err
}

// ParseAuthJwtToken parsing token
func (c TokenParser) ParseAuthJwtToken(token string, secretKey string) (*commonImport.AuthJwtClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &commonImport.AuthJwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, fmt.Errorf("%w. %w", jwt.ErrTokenExpired, resources.ErrClient)
		}
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*commonImport.AuthJwtClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func (c TokenParser) getToken(data commonImport.AuthJwtTokenInput, expiredDay int) (string, error) {
	return c.generateAuthJwtToken(data,
		config.AppSetting.JWTSECRET, time.Now().AddDate(0, 0, expiredDay))
}

func (c TokenParser) CreateAuthToken(data commonImport.AuthJwtTokenInput) (string, error) {
	return c.getToken(data, config.AppSetting.AccessTokenExpiredTime)
}
