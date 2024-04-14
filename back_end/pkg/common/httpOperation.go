package commonimport

import (
	"admin-panel/pkg/config"
	"admin-panel/pkg/resources"
	"context"
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	"strconv"
	"sync"
)

type HTTPOperation struct {
	Logger         *Logger
	Authenticator  Authenticator
	config         *config.Config
	UserLogManager UserLogRecorder
}

func (r HTTPOperation) Success(ctx *gin.Context, log UserLog, code int, data any) {
	log.Status = int32(code)
	go r.recordUserLog(ctx, log)
	ctx.Header("x-request-id", ctx.GetHeader("x-request-id"))
	ctx.JSON(code, data)
}

// ResponseMessage is needed to change then go also change app.ResponseMessage.
type ResponseMessage struct {
	Message string `json:"message"`
}

func (r HTTPOperation) Error(ctx *gin.Context, log UserLog, err error) {
	log.ErrorMessage = err.Error()
	r.Logger.LogError(r.GetLogDataContext(ctx), err)
	code, message := getErrorResponse(err)
	r.Success(ctx, log, code, ResponseMessage{Message: message})
}

func (r HTTPOperation) recordUserLog(ctx context.Context, log UserLog) {
	err := r.UserLogManager.Record(ctx, log)
	if err != nil {
		r.Logger.LogError(ctx, err)
	}
}

func getErrorResponse(err error) (int, string) {
	var holder ResponseError
	ok := errors.As(err, &holder)
	if ok {
		if holder.Code == http.StatusInternalServerError {
			return http.StatusInternalServerError, resources.ErrInternalServer.Error()
		}
		return holder.Code, holder.Message
	}

	code, message, ok := PostgresErrorTransform(err)
	if ok {
		return code, message
	}
	switch {
	case errors.Is(err, resources.ErrNoRecordFound), errors.Is(err, resources.ErrClient),
		errors.Is(err, resources.ErrUnmarshalData):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, resources.ErrUnAuthorized):
		return http.StatusUnauthorized, err.Error()
	case errors.Is(err, resources.ErrForbidden):
		return http.StatusForbidden, err.Error()
	case errors.Is(err, resources.ErrBadGateway):
		return http.StatusBadGateway, resources.ErrBadGateway.Error()
	case errors.Is(err, resources.ErrNotAcceptable):
		return http.StatusNotAcceptable, err.Error()
	default:
		return http.StatusInternalServerError, resources.ErrInternalServer.Error()
	}
}

func PostgresErrorTransform(err error) (int, string, bool) {
	if err == nil {
		return 0, "", false
	}
	pgErr := &pgconn.PgError{}
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return http.StatusConflict, resources.ErrDuplicateValue.Error(), true
		case pgerrcode.ForeignKeyViolation:
			return http.StatusBadRequest, resources.ErrNoRecordFound.Error(), true
		default:
			return http.StatusInternalServerError, resources.ErrInternalServer.Error(), true
		}
	}

	return 0, "", false
}

type Validator interface {
	Valid() error
}

// BindAndValid binds and validates data
func (r HTTPOperation) BindAndValid(ctx *gin.Context, form interface{}) error {
	err := ctx.Bind(form)
	if err != nil {
		ctx.Set(string(resources.LogData), r.AddPayLoadToLog(r.GetLogDataContext(ctx), form))
		return errors.Wrap(resources.ErrClient, err.Error())
	}
	buffer, ok := form.(Validator)
	if ok {
		err = buffer.Valid()
		if err != nil {
			ctx.Set(string(resources.LogData), r.AddPayLoadToLog(r.GetLogDataContext(ctx), form))
			return err
		}
	}
	return nil
}

func (r HTTPOperation) LimitBodySize() gin.HandlerFunc {
	return func(c *gin.Context) {
		var w http.ResponseWriter = c.Writer
		c.Request.Body = http.MaxBytesReader(w, c.Request.Body, 150*1024*1024)
		c.Next()
	}

}

func (r HTTPOperation) AddDataToLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		dataBag := map[string]interface{}{}
		dataBag["method"] = c.Request.Method
		dataBag["url"] = c.Request.URL.String()
		requestID := c.Request.Header.Get("x-request-id")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		dataBag["requestId"] = requestID
		ctx := context.WithValue(context.Background(), resources.Data, dataBag)
		c.Set(string(resources.LogData), ctx)
		c.Next()
	}
}

func (r HTTPOperation) GetLogDataContext(c *gin.Context) context.Context {
	if c == nil {
		return nil
	}
	buffer, exist := c.Get(string(resources.LogData))
	if exist {
		ctx, ok := buffer.(context.Context)
		if ok {
			return ctx
		}
	}
	return nil
}

func (r HTTPOperation) BadRequestError(message string) error {
	return ResponseError{Code: http.StatusBadRequest, Message: message}
}

func (r HTTPOperation) UnauthorizedError(message string) error {
	return ResponseError{Code: http.StatusUnauthorized, Message: message}
}

func (r HTTPOperation) ForbiddenError(message string) error {
	return ResponseError{Code: http.StatusForbidden, Message: message}
}

func (r HTTPOperation) InternalServerError(message string) error {
	return ResponseError{Code: http.StatusInternalServerError, Message: message}
}

func (r HTTPOperation) GetID(ctx *gin.Context, key string) (int64, error) {
	if ctx == nil {
		return 0, resources.ErrInternalServer
	}
	idStr := ctx.Param(key)
	if idStr == "" || idStr == "0" {
		return 0, ResponseError{Code: http.StatusBadRequest, Message: "invalid request data"}
	}
	return strconv.ParseInt(idStr, 10, 64)
}

func (r HTTPOperation) ValidateAdminUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		userData, err := r.Authenticator.AuthenticateAdmin(c.Request.Context(), c.GetHeader("Authorization"))
		if err != nil {
			userLog := GetUserLog(c, userData.UserID, nil)
			ctx := r.GetLogDataContext(c)
			r.Logger.LogError(ctx, err)
			r.Error(c, userLog, err)
			c.Abort()
			return
		}
		ctx := r.GetLogDataContext(c)
		c.Set(string(resources.LogData), r.Logger.AddDataToLog(ctx, map[string]interface{}{"user": userData}))
		c.Set(string(resources.UserData), userData)
		c.Next()
	}
}

func (r HTTPOperation) ValidateNormalUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		userData, err := r.Authenticator.Authenticate(c.Request.Context(), c.GetHeader("Authorization"))
		if err != nil {
			userLog := GetUserLog(c, userData.UserID, nil)
			ctx := r.GetLogDataContext(c)
			r.Logger.LogError(ctx, err)
			r.Error(c, userLog, err)
			c.Abort()
			return
		}
		ctx := r.GetLogDataContext(c)
		c.Set(string(resources.LogData), r.Logger.AddDataToLog(ctx, map[string]interface{}{"user": userData}))
		c.Set(string(resources.UserData), userData)
		c.Next()
	}
}

func (r HTTPOperation) AddPayLoadToLog(ctx context.Context, data any) context.Context {
	return r.Logger.AddDataToLog(ctx, map[string]interface{}{string(resources.PayLoad): data})
}

func (r HTTPOperation) GetUserData(c *gin.Context) UserData {
	buffer, exist := c.Get(string(resources.UserData))
	if exist {
		employeeInformation, ok := buffer.(UserData)
		if ok {
			return employeeInformation
		}
	}
	return UserData{}
}

var operation *HTTPOperation
var setUpOperationOnce sync.Once

func ProvideHTTPOperation(logger *Logger, authenticator Authenticator, config *config.Config, userLogManager UserLogRecorder) *HTTPOperation {
	setUpOperationOnce.Do(func() {
		operation = &HTTPOperation{
			Logger:         logger,
			Authenticator:  authenticator,
			UserLogManager: userLogManager,
			config:         config}
	})
	return operation
}
