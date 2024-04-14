//go:build wireinject
// +build wireinject

package dependency_manager

import (
	userPostgresAdapter "admin-panel/internal/user/adapter/postgresAdapter"
	userService "admin-panel/internal/user/service"
	commonImport "admin-panel/pkg/common"
	"admin-panel/pkg/config"
	"admin-panel/pkg/jwt"
	passwordPackage "admin-panel/pkg/password"
	"admin-panel/pkg/postgres"
	"github.com/google/wire"
)

var SuperSet = wire.NewSet(commonImport.ProvideLogger, commonImport.ProvideHTTPOperation, commonImport.ProvideContainer,
	postgres.ProvideSqlx, config.ProvideConfig, postgres.ProvideClient, jwt.ProvideTokenParser, userPostgresAdapter.ProvideUserRepositoryInterface,
	userPostgresAdapter.ProvideUserLogRepositoryInterface, userPostgresAdapter.ProvideUserSessionRepository, userService.ProvideUserManager,
	userService.ProvideUserService, userService.ProvideUserLogRecorder, userService.ProvideAuthenticator, passwordPackage.ProvidePasswordManager)

func InitializeContainer() *commonImport.Container {
	wire.Build(SuperSet)
	return &commonImport.Container{}
}
