package http

import (
	v1 "admin-panel/internal/user/port/public/http/api/v1"
	commonImport "admin-panel/pkg/common"
	"github.com/gin-gonic/gin"
)

func SetupAPI(router *gin.Engine, container *commonImport.Container) {
	normalFeatureAPI := v1.NewNormalFeatureAPI()
	userManagementAPI := v1.NewEmployeeManagementAPI()

	publicAPIV1 := router.Group("/v1")

	{
		publicAPIV1.POST("/login", normalFeatureAPI.UserLoginV1)

	}

	privateAPI := router.Group("/v1")
	privateAPI.Use(container.Operation.ValidateNormalUser())
	{

		privateAPI.DELETE("/logout", normalFeatureAPI.UserLogOutV1)
	}
	adminAPI := router.Group("/v1")
	adminAPI.Use(container.Operation.ValidateAdminUser())

	{

		adminAPI.PUT("/user", userManagementAPI.UpdateUser)
		adminAPI.POST("/user", userManagementAPI.CreateUserAccount)
		adminAPI.GET("/user", userManagementAPI.ListUsers)
		adminAPI.GET("/user-logs", userManagementAPI.ListUserLogs)
		adminAPI.GET("/user/:id", userManagementAPI.GetUserByID)
		adminAPI.DELETE("/user/:id", userManagementAPI.DeleteUserAccount)

	}
}
