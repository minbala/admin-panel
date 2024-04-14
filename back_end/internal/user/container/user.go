package container

import (
	"admin-panel/internal/user/service"
	commonImport "admin-panel/pkg/common"
	"sync"
)

type UserContainer struct {
	Operation   *commonImport.HTTPOperation
	UserService *service.UserService
}

func ProvideAuthAuthorizeContainer(userService *service.UserService,
	operation *commonImport.HTTPOperation) *UserContainer {
	return &UserContainer{
		UserService: userService,
		Operation:   operation,
	}
}

var userContainer *UserContainer
var setUpAuthAuthorizeContainerOnce sync.Once

func GetContainer() *UserContainer {
	setUpAuthAuthorizeContainerOnce.Do(func() {
		userContainer = InitializeContainer()
	})
	return userContainer
}
