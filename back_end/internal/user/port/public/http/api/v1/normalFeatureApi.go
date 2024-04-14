package v1

import (
	"admin-panel/internal/user/container"
	"admin-panel/internal/user/port/public/http/schema"
	"admin-panel/internal/user/service"
	commonImport "admin-panel/pkg/common"
	"admin-panel/pkg/resources"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type NormalFeatureAPI struct {
}

func NewNormalFeatureAPI() *NormalFeatureAPI {
	return &NormalFeatureAPI{}
}

// UserLoginV1   godoc
//
//	@Summary		User Login
//	@Description	User Login Version 1
//	@Tags			NormalFeatureAPI
//	@Accept			application/json
//	@Produce		application/json
//	@Param			requestData	body		schema.LoginRequest	true	"Request body in JSON format"
//	@Success		201			{object}	schema.LoginResponse
//	@Failure		400			{object}	app.ResponseMessage
//	@Failure		500			{object}	app.ResponseMessage
//	@Router			/v1/login [post]
func (NormalFeatureAPI) UserLoginV1(c *gin.Context) {
	app := container.GetContainer()
	logContext := app.Operation.GetLogDataContext(c)

	var requestData schema.LoginRequest
	err := app.Operation.BindAndValid(c, &requestData)
	userLog := commonImport.GetUserLog(c, 0, requestData)
	if err != nil {
		app.Operation.Error(c, userLog, err)
		return
	}

	userLog.Data = requestData
	logContext = app.Operation.AddPayLoadToLog(logContext, requestData)

	var accessToken string
	accessToken, userLog.UserID, err = app.UserService.LoginUser(logContext, service.LoginCmd{
		Email:    requestData.Email,
		Password: requestData.Password,
	})
	if err != nil {
		app.Operation.Error(c, userLog, err)
		return
	}
	app.Operation.Success(c, userLog, http.StatusCreated, schema.LoginResponse{
		AccessToken: accessToken,
	})
}

// UserLogOutV1 go doc
//
//	@Summary		User logout
//	@Description	User logout Version 1
//	@Tags			NormalFeatureAPI
//	@Accept			application/json
//	@Produce		application/json
//	@Success		204	{object}	nil
//	@Failure		400	{object}	app.ResponseMessage
//	@Failure		500	{object}	app.ResponseMessage
//	@Router			/v1/logout [delete]
//
//	@Security		Bearer
func (NormalFeatureAPI) UserLogOutV1(c *gin.Context) {
	app := container.GetContainer()
	userData := app.Operation.GetUserData(c)
	logContext := app.Operation.GetLogDataContext(c)
	userLog := commonImport.GetUserLog(c, userData.UserID, nil)
	logContext = app.Operation.AddPayLoadToLog(logContext, c.GetHeader("Authorization"))
	buffer := strings.Split(c.GetHeader("Authorization"), " ")
	if len(buffer) < 2 {
		app.Operation.Error(c, userLog, resources.ErrClient)
		return
	}
	err := app.UserService.LogoutUser(app.Operation.GetLogDataContext(c), buffer[1], userData.UserID)
	if err != nil {
		app.Operation.Error(c, userLog, err)
		return
	}
	app.Operation.Success(c, userLog, http.StatusNoContent, nil)

}
