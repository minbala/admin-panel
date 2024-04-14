package service

import (
	"admin-panel/internal/user/domain"
	commonImport "admin-panel/pkg/common"
	"admin-panel/pkg/config"
	"admin-panel/pkg/resources"
	"context"
	"strings"
)

type Authenticator struct {
	employeeManager UserManager
	jwtManager      commonImport.JwtManager
	config          *config.Config
	userSession     domain.SessionRepositoryInterface
}

type UserManager interface {
	GetUserData(ctx context.Context, userID int64) (commonImport.UserData, error)
}

func isStringValid(testString string) bool {
	if m := strings.TrimSpace(testString); len(m) == 0 || len(testString) == 0 {
		return false
	}
	return true
}

func parseToken(token string) (string, bool) {
	if !isStringValid(token) {
		return "", false
	}
	validationTokens := strings.Split(token, " ")
	if len(validationTokens) != 2 {
		return "", false
	}
	return validationTokens[1], true
}

func (a Authenticator) validAndGetTokenData(ctx context.Context, token string) (data commonImport.AuthJwtTokenInput, err error) {
	if !isStringValid(token) {
		return commonImport.AuthJwtTokenInput{}, resources.ErrUnAuthorized
	}
	token, ok := parseToken(token)
	if !ok {
		return commonImport.AuthJwtTokenInput{}, resources.ErrUnAuthorized
	}
	authData, err := a.jwtManager.ParseAuthJwtToken(token, a.config.App.JWTSECRET)
	if err != nil {
		return commonImport.AuthJwtTokenInput{}, resources.ErrUnAuthorized
	}
	err = a.validateToken(ctx, authData.UserID, token)
	if err != nil {
		return commonImport.AuthJwtTokenInput{}, resources.ErrUnAuthorized
	}

	return authData.AuthJwtTokenInput, nil
}

func (a Authenticator) validateToken(ctx context.Context, userID int64, token string) error {
	_, err := a.userSession.Get(ctx, userID, token)
	if err != nil {
		return resources.ErrUnAuthorized
	}
	return nil
}

func (a Authenticator) Authenticate(ctx context.Context, token string) (commonImport.UserData, error) {
	tokenData, err := a.validAndGetTokenData(ctx, token)
	if err != nil {
		return commonImport.UserData{}, err
	}
	return a.employeeManager.GetUserData(ctx, tokenData.UserID)
}

func (a Authenticator) AuthenticateAdmin(ctx context.Context, token string) (commonImport.UserData, error) {
	userData, err := a.Authenticate(ctx, token)
	if err != nil {
		return commonImport.UserData{}, err
	}
	if resources.UserRoleType(userData.UserRole) != resources.UserRoleTypeAdmin {
		return commonImport.UserData{}, resources.ErrUnAuthorized
	}

	return userData, nil
}

func ProvideAuthenticator(employeeManager UserManager,
	jwtManager commonImport.JwtManager, config *config.Config, sessionRepo domain.SessionRepositoryInterface) commonImport.Authenticator {
	return &Authenticator{
		jwtManager:      jwtManager,
		config:          config,
		employeeManager: employeeManager,
		userSession:     sessionRepo}
}
