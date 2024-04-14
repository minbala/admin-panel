package v1

import (
	"admin-panel/internal/user/container"
	"admin-panel/internal/user/domain"
	"admin-panel/internal/user/port/public/http/schema"
	"admin-panel/internal/user/service"
	commonImport "admin-panel/pkg/common"
	"admin-panel/pkg/resources"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserManagementAPI struct {
}

func NewEmployeeManagementAPI() *UserManagementAPI {
	return &UserManagementAPI{}
}

// CreateUserAccount go doc
//
//	@Summary		create employee account
//	@Description	create employee account, assign user to a role but u can't  assign to owner role
//	@Tags			UserManagementAPI
//	@Accept			application/json
//	@Produce		application/json
//	@Param			requestData	body	schema.UserCreateRequest	true	"Request body in JSON format"
//	@Security		Bearer
//	@Success		201	{object}	nil
//	@Failure		400	{object}	app.ResponseMessage
//	@Failure		500	{object}	app.ResponseMessage
//	@Router			/v1/user [post]
func (UserManagementAPI) CreateUserAccount(c *gin.Context) {
	app := container.GetContainer()
	userData := app.Operation.GetUserData(c)
	logContext := app.Operation.GetLogDataContext(c)

	var requestData schema.UserCreateRequest
	err := app.Operation.BindAndValid(c, &requestData)
	userLog := commonImport.GetUserLog(c, userData.UserID, requestData)
	if err != nil {
		app.Operation.Error(c, userLog, err)
		return
	}

	userLog.Data = requestData
	logContext = app.Operation.AddPayLoadToLog(logContext, requestData)

	_, err = app.UserService.CreateUser(logContext, service.UserCreateCmd{
		Email:    requestData.Email,
		Name:     requestData.Name,
		Password: requestData.Password,
		UserRole: resources.UserRoleType(requestData.UserRole),
	})
	if err != nil {
		app.Operation.Error(c, userLog, err)
		return
	}
	app.Operation.Success(c, userLog, http.StatusCreated, nil)
}

// ListUsers   godoc
//
//	@Summary		List User Information
//	@Description	List Users Information
//	@Tags			UserManagementAPI
//	@Accept			application/json
//	@Produce		application/json
//	@Param			name		query		string	false	"bala"
//	@Param			userRole	query		string	false	"user"
//	@Param			limit		query		int32	false	"20"
//	@Param			offset		query		int32	false	"0"
//
//	@Success		200			{object}	schema.ListUserResponse
//	@Failure		400			{object}	app.ResponseMessage
//	@Failure		500			{object}	app.ResponseMessage
//	@Router			/v1/user [get]
//
//	@Security		Bearer
func (UserManagementAPI) ListUsers(c *gin.Context) {
	app := container.GetContainer()
	userData := app.Operation.GetUserData(c)
	logContext := app.Operation.GetLogDataContext(c)

	var requestData schema.UsersListRequest
	err := app.Operation.BindAndValid(c, &requestData)
	userLog := commonImport.GetUserLog(c, userData.UserID, requestData)
	if err != nil {
		app.Operation.Error(c, userLog, err)
		return
	}
	requestData.SqlFilter.Setup()
	userLog.Data = requestData
	logContext = app.Operation.AddPayLoadToLog(logContext, requestData)

	users, count, err := app.UserService.ListUsers(logContext, domain.UserQuery{
		UserRole: resources.UserRoleType(requestData.UserRole),
		Name:     requestData.Name,
		Limit:    requestData.Limit,
		Offset:   requestData.Offset,
	})
	if err != nil {
		app.Operation.Error(c, userLog, err)
		return
	}
	app.Operation.Success(c, userLog, http.StatusOK, schema.ListUserResponse{
		Users: schema.MarshalUsersToSchema(users),
		Total: count,
	})
}

// GetUserByID go doc
//
//	@Summary		delete  user account
//	@Description	delete  user account
//	@Tags			UserManagementAPI
//	@Accept			application/json
//	@Produce		application/json
//	@Param			id	path	int64	true	"user id "
//	@Security		Bearer
//	@Success		200	{object}	schema.User
//	@Failure		400	{object}	app.ResponseMessage
//	@Failure		500	{object}	app.ResponseMessage
//	@Router			/v1/user/{id} [get]
func (UserManagementAPI) GetUserByID(c *gin.Context) {
	app := container.GetContainer()
	userData := app.Operation.GetUserData(c)
	logContext := app.Operation.GetLogDataContext(c)
	userLog := commonImport.GetUserLog(c, userData.UserID, struct {
		UserID int64
	}{UserID: 0})
	userID, err := app.Operation.GetID(c, "id")
	if err != nil {
		app.Operation.Error(c, userLog, err)
		return
	}
	logContext = app.Operation.AddPayLoadToLog(logContext, userID)
	userLog.Data = struct {
		UserID int64
	}{UserID: userID}

	user, err := app.UserService.GetUser(logContext, domain.UserQuery{ID: userID})
	if err != nil {
		app.Operation.Error(c, userLog, err)
		return
	}
	app.Operation.Success(c, userLog, http.StatusOK, schema.MarshalUserToSchema(user))
}

// DeleteUserAccount go doc
//
//	@Summary		delete  user account
//	@Description	delete  user account
//	@Tags			UserManagementAPI
//	@Accept			application/json
//	@Produce		application/json
//	@Param			id	path	int64	true	"user id"
//	@Security		Bearer
//	@Success		204
//	@Failure		400	{object}	app.ResponseMessage
//	@Failure		500	{object}	app.ResponseMessage
//	@Router			/v1/user/{id} [delete]
func (UserManagementAPI) DeleteUserAccount(c *gin.Context) {
	app := container.GetContainer()
	userData := app.Operation.GetUserData(c)
	logContext := app.Operation.GetLogDataContext(c)
	userLog := commonImport.GetUserLog(c, userData.UserID, struct {
		UserID int64
	}{UserID: 0})
	userID, err := app.Operation.GetID(c, "id")
	if err != nil {
		app.Operation.Error(c, userLog, err)
		return
	}
	logContext = app.Operation.AddPayLoadToLog(logContext, userID)
	userLog.Data = struct {
		UserID int64
	}{UserID: userID}

	if userID == userData.UserID {
		app.Operation.Error(c, userLog, resources.ErrClient)
		return
	}
	err = app.UserService.DeleteUser(logContext, userID)
	if err != nil {
		app.Operation.Error(c, userLog, err)
		return
	}
	app.Operation.Success(c, userLog, http.StatusNoContent, nil)
}

// UpdateUser go doc
//
//	@Summary		User update  user
//	@Description	User update user
//	@Tags			UserManagementAPI
//	@Accept			application/json
//	@Produce		application/json
//	@Param			requestData	body		schema.UserUpdateRequest	true	"Request body in JSON format"
//	@Success		200			{object}	nil
//	@Failure		400			{object}	app.ResponseMessage
//	@Failure		401			{object}	app.ResponseMessage
//	@Failure		403			{object}	app.ResponseMessage
//	@Failure		500			{object}	app.ResponseMessage
//	@Router			/v1/user [put]
//
//	@Security		Bearer
func (UserManagementAPI) UpdateUser(c *gin.Context) {
	app := container.GetContainer()
	userData := app.Operation.GetUserData(c)
	logContext := app.Operation.GetLogDataContext(c)

	var requestData schema.UserUpdateRequest
	userLog := commonImport.GetUserLog(c, userData.UserID, requestData)
	err := app.Operation.BindAndValid(c, &requestData)
	if err != nil {
		app.Operation.Error(c, userLog, err)
		return
	}

	userLog.Data = requestData
	logContext = app.Operation.AddPayLoadToLog(logContext, requestData)

	if requestData.UserId == userData.UserID {
		app.Operation.Error(c, userLog, resources.ErrClient)
		return
	}
	err = app.UserService.UpdateUser(logContext, service.UserUpdateCmd{
		Email:    requestData.Email,
		Name:     requestData.Name,
		Password: requestData.Password,
		UserRole: resources.UserRoleType(requestData.UserRole),
		UserID:   requestData.UserId,
	})
	if err != nil {
		app.Operation.Error(c, userLog, err)
		return
	}
	app.Operation.Success(c, userLog, http.StatusOK, nil)
}

// ListUserLogs   godoc
//
//	@Summary		List User Information
//	@Description	List Users Information
//	@Tags			UserManagementAPI
//	@Accept			application/json
//	@Produce		application/json
//	@Param			userID	query		int64	false	"2"
//	@Param			limit	query		int32	false	"20"
//	@Param			offset	query		int32	false	"0"
//
//	@Success		200		{object}	schema.ListUserLogsResponse
//	@Failure		400		{object}	app.ResponseMessage
//	@Failure		500		{object}	app.ResponseMessage
//	@Router			/v1/user-logs [get]
//
//	@Security		Bearer
func (UserManagementAPI) ListUserLogs(c *gin.Context) {
	app := container.GetContainer()
	userData := app.Operation.GetUserData(c)
	logContext := app.Operation.GetLogDataContext(c)

	var requestData schema.UserLogsListRequest
	err := app.Operation.BindAndValid(c, &requestData)
	userLog := commonImport.GetUserLog(c, userData.UserID, requestData)
	if err != nil {
		app.Operation.Error(c, userLog, err)
		return
	}
	requestData.SqlFilter.Setup()
	userLog.Data = requestData
	logContext = app.Operation.AddPayLoadToLog(logContext, requestData)

	userLogs, count, err := app.UserService.ListUserLogs(logContext, domain.UserLogQuery{
		UserID: requestData.UserID,
		Limit:  requestData.Limit,
		Offset: requestData.Offset,
	})
	if err != nil {
		app.Operation.Error(c, userLog, err)
		return
	}
	app.Operation.Success(c, userLog, http.StatusOK, schema.ListUserLogsResponse{
		UserLogs: schema.MarshalUserLogsToSchema(userLogs),
		Total:    count,
	})
}
