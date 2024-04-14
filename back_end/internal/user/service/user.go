package service

import (
	"admin-panel/internal/user/domain"
	commonImport "admin-panel/pkg/common"
	"admin-panel/pkg/config"
	"admin-panel/pkg/resources"
	"context"
	"github.com/bytedance/sonic"
	"github.com/cockroachdb/errors"
	"sync"
	"time"
)

type UserService struct {
	userRepo        domain.UserRepositoryInterface
	passwordManager commonImport.PasswordManager
	tokenManager    commonImport.JwtManager
	userSession     domain.SessionRepositoryInterface
	config          *config.Config
	userLogRepo     domain.UserLogRepositoryInterface
}

var userService *UserService
var initializeUserServiceOnce sync.Once

func ProvideUserService(userRepo domain.UserRepositoryInterface,
	passwordManager commonImport.PasswordManager,
	tokenManager commonImport.JwtManager, sessionRepo domain.SessionRepositoryInterface,
	config *config.Config, userLogRepo domain.UserLogRepositoryInterface) *UserService {
	initializeUserServiceOnce.Do(func() {
		userService = &UserService{
			userRepo:        userRepo,
			passwordManager: passwordManager,
			tokenManager:    tokenManager,
			userSession:     sessionRepo,
			config:          config,
			userLogRepo:     userLogRepo,
		}
	})
	return userService
}

func ProvideUserManager(userRepo domain.UserRepositoryInterface,
	passwordManager commonImport.PasswordManager,
	tokenManager commonImport.JwtManager, sessionRepo domain.SessionRepositoryInterface,
	config *config.Config, userLogRepo domain.UserLogRepositoryInterface) UserManager {
	initializeUserServiceOnce.Do(func() {
		userService = &UserService{
			userRepo:        userRepo,
			passwordManager: passwordManager,
			tokenManager:    tokenManager,
			userSession:     sessionRepo,
			config:          config,
			userLogRepo:     userLogRepo,
		}
	})
	return userService

}

func ProvideUserLogRecorder(userRepo domain.UserRepositoryInterface,
	passwordManager commonImport.PasswordManager,
	tokenManager commonImport.JwtManager, sessionRepo domain.SessionRepositoryInterface,
	config *config.Config, userLogRepo domain.UserLogRepositoryInterface) commonImport.UserLogRecorder {
	initializeUserServiceOnce.Do(func() {
		userService = &UserService{
			userRepo:        userRepo,
			passwordManager: passwordManager,
			tokenManager:    tokenManager,
			userSession:     sessionRepo,
			config:          config,
			userLogRepo:     userLogRepo,
		}
	})
	return userService

}

type LoginCmd struct {
	Email    string
	Password string
}

type SessionCreateCmd struct {
	EmployeeId int64
}

type UserPasswordCheckCmd struct {
	Email    string
	Password string `json:"password"`
	UserID   int64
}

func (c UserService) LoginUser(ctx context.Context, cmd LoginCmd) (string, int64, error) {
	user, err := c.CheckPassword(ctx, UserPasswordCheckCmd{
		Email:    cmd.Email,
		Password: cmd.Password,
	})
	if err != nil {
		return "", 0, err
	}
	accessToken, err := c.createSession(ctx, user.ID)
	if err != nil {
		return "", 0, err
	}
	return accessToken, user.ID, nil
}

// CheckPassword check password if success return user data
func (c UserService) CheckPassword(ctx context.Context, input UserPasswordCheckCmd) (domain.User, error) {
	user, err := c.GetUser(ctx, domain.UserQuery{
		ID:    input.UserID,
		Email: input.Email,
	})
	if err != nil {
		return user, err
	}
	err = c.passwordManager.CheckPassword(input.Password, user.Password)
	if err != nil {
		return user, errors.Wrap(resources.ErrClient, "email or password is incorrect")
	}

	return user, nil
}

func (c UserService) GetUser(ctx context.Context, input domain.UserQuery) (domain.User, error) {
	input.Limit = 1
	userData, err := c.userRepo.Get(ctx, input)
	if err != nil {
		return domain.User{}, resources.ErrNoRecordFound
	}
	if len(userData) == 0 {
		return domain.User{}, resources.ErrNoRecordFound
	}
	return userData[0], nil
}

func (c UserService) createSession(ctx context.Context, userID int64) (string, error) {
	accessToken, err := c.tokenManager.CreateAuthToken(commonImport.AuthJwtTokenInput{
		UserID: userID,
	})
	if err != nil {
		return "", err
	}
	err = c.userSession.Create(ctx, domain.Session{UserID: userID, AccessToken: accessToken})
	return accessToken, err
}

func (c UserService) GetUserData(ctx context.Context, userID int64) (commonImport.UserData, error) {
	userData, err := c.GetUser(ctx, domain.UserQuery{
		ID:    userID,
		Limit: 1,
	})
	if err != nil {
		return commonImport.UserData{}, err
	}
	return convertToCommonModel(userData), nil
}

func convertToCommonModel(input domain.User) commonImport.UserData {
	return commonImport.UserData{
		UserID:   input.ID,
		UserRole: string(input.UserRole),
	}
}

func (c UserService) LogoutUser(ctx context.Context, authToken string, userID int64) error {
	err := c.userSession.Delete(ctx, domain.Session{UserID: userID, AccessToken: authToken})
	if err != nil {
		return err
	}
	return nil
}

func (c UserService) parseToken(token string) (authData *commonImport.AuthJwtClaims, err error) {
	return c.tokenManager.ParseAuthJwtToken(token, c.config.App.JWTSECRET)
}

type UserCreateCmd struct {
	Email    string
	Name     string
	Password string
	UserRole resources.UserRoleType
}

func (c UserService) CreateUser(ctx context.Context, cmd UserCreateCmd) (int64, error) {
	var err error
	user := domain.User{
		Email:    cmd.Email,
		Name:     cmd.Name,
		Password: cmd.Password,
		UserRole: cmd.UserRole,
	}
	user.Password, err = c.passwordManager.Hash(cmd.Password)
	if err != nil {
		return 0, err
	}
	user.ID, err = c.userRepo.Create(ctx, user)
	return user.ID, err
}

func (c UserService) ListUsers(ctx context.Context, input domain.UserQuery) ([]domain.User,
	int64, error) {
	var count int64
	stream := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		var err error
		count, err = c.userRepo.Count(ctx, input)
		stream <- err
	}()
	users, err := c.userRepo.Get(ctx, input)
	if err != nil {
		return nil, 0, err
	}
	if err = <-stream; err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (c UserService) DeleteUser(ctx context.Context, userID int64) error {
	_, err := c.GetUser(ctx, domain.UserQuery{ID: userID})
	if err != nil {
		return err
	}
	return c.userRepo.Delete(ctx, userID)
}

type UserUpdateCmd struct {
	Email    string
	Name     string
	Password string
	UserRole resources.UserRoleType
	UserID   int64
}

func (c UserService) UpdateUser(ctx context.Context, cmd UserUpdateCmd) error {
	user, err := c.GetUser(ctx, domain.UserQuery{ID: cmd.UserID})
	if err != nil {
		return err
	}
	if isStringValid(cmd.Password) {
		user.Password, err = c.passwordManager.Hash(cmd.Password)
		if err != nil {
			return err
		}
	}
	user.UserRole = cmd.UserRole
	user.Email = cmd.Email
	user.Name = cmd.Name
	return c.userRepo.Update(ctx, user)
}

func (c UserService) ListUserLogs(ctx context.Context, input domain.UserLogQuery) ([]domain.UserLog,
	int64, error) {
	var count int64
	stream := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		var err error
		count, err = c.userLogRepo.Count(ctx, input)
		stream <- err
	}()
	userLogs, err := c.userLogRepo.Get(ctx, input)
	if err != nil {
		return nil, 0, err
	}
	if err = <-stream; err != nil {
		return nil, 0, err
	}
	return userLogs, count, nil
}

func (c UserService) Record(ctx context.Context, log commonImport.UserLog) error {
	data, err := sonic.Marshal(log.Data)
	if err != nil {
		return err
	}
	return c.userLogRepo.Create(ctx, domain.UserLog{
		ID:           0,
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
		UserID:       log.UserID,
		Event:        log.Event,
		RequestUrl:   log.RequestUrl,
		Data:         data,
		Status:       log.Status,
		ErrorMessage: log.ErrorMessage,
	})
}
