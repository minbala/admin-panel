//go:build wireinject
// +build wireinject

package container

import (
	"admin-panel/dependency_manager"
	"github.com/google/wire"
)

func InitializeContainer() *UserContainer {
	wire.Build(dependency_manager.SuperSet, ProvideAuthAuthorizeContainer)
	return &UserContainer{}
}
