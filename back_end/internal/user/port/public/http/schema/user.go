package schema

import (
	"admin-panel/internal/user/domain"
	commonImport "admin-panel/pkg/common"
	"admin-panel/pkg/resources"
	"fmt"
	"github.com/cockroachdb/errors"
	"net/mail"
	"strconv"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l LoginRequest) Valid() error {
	err := ValidateEmail(l.Email)
	if err != nil {
		return errors.Wrap(resources.ErrClient, err.Error())
	}
	if !commonImport.IsStringValid(l.Password) {
		return fmt.Errorf("password is required. %w", resources.ErrClient)
	}
	return nil
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type UserCreateRequest struct {
	Email    string ` json:"email"`
	Name     string ` json:"name"`
	Password string ` json:"password"`
	UserRole string `json:"userRole"`
}

func (l UserCreateRequest) Valid() error {
	err := ValidateEmail(l.Email)
	if err != nil {
		return errors.Wrap(resources.ErrClient, err.Error())
	}
	if !commonImport.IsStringValid(l.Name) {
		return fmt.Errorf("name is required. %w", resources.ErrClient)
	}
	if !commonImport.IsStringValid(l.Password) {
		return fmt.Errorf("password is required. %w", resources.ErrClient)
	}
	err = resources.UserRoleTypeValid(l.UserRole)
	if err != nil {
		return err
	}
	return nil
}

type UsersListRequest struct {
	Name     string `json:"name" form:"name"`
	UserRole string `json:"userRole" form:"userRole"`
	commonImport.SqlFilter
}
type User struct {
	Id        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	UserRole  string    `json:"userRole"`
	CreatedAt time.Time `json:"createdAt"`
}

type ListUserResponse struct {
	Users []User `json:"users"`
	Total int64  `json:"total"`
}

func MarshalUsersToSchema(input []domain.User) []User {
	response := make([]User, len(input))
	for i, user := range input {
		response[i] = User{
			Id:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			UserRole:  string(user.UserRole),
			CreatedAt: user.CreatedAt,
		}
	}
	return response
}

func MarshalUserToSchema(user domain.User) User {
	return User{
		Id:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		UserRole:  string(user.UserRole),
		CreatedAt: user.CreatedAt,
	}
}

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

type UserUpdateRequest struct {
	UserId   int64  `json:"userId"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	UserRole string `json:"userRole"`
}

func (l UserUpdateRequest) Valid() error {
	err := ValidateEmail(l.Email)
	if err != nil {
		return errors.Wrap(resources.ErrClient, err.Error())
	}
	if !commonImport.IsIDValid(l.UserId) {
		return fmt.Errorf("user id is required. %w", resources.ErrClient)
	}
	if !commonImport.IsStringValid(l.Name) {
		return fmt.Errorf("name is required. %w", resources.ErrClient)
	}
	err = resources.UserRoleTypeValid(l.UserRole)
	if err != nil {
		return err
	}
	return nil
}

type UserLog struct {
	Id           int64     `json:"id"`
	UserId       int64     `json:"userId"`
	Event        string    `json:"event"`
	RequestUrl   string    `json:"requestUrl"`
	Data         string    `json:"data"`
	Status       string    `json:"status"`
	ErrorMessage string    `json:"errorMessage"`
	CreatedAt    time.Time `json:"createdAt"`
}

type ListUserLogsResponse struct {
	UserLogs []UserLog `json:"userLogs"`
	Total    int64     `json:"total"`
}

func MarshalUserLogsToSchema(input []domain.UserLog) []UserLog {
	response := make([]UserLog, len(input))
	for i, log := range input {
		response[i] = UserLog{
			Id:           log.ID,
			UserId:       log.UserID,
			Event:        log.Event,
			RequestUrl:   log.RequestUrl,
			Data:         string(log.Data),
			Status:       strconv.Itoa(int(log.Status)),
			ErrorMessage: log.ErrorMessage,
			CreatedAt:    log.CreatedAt,
		}
	}
	return response
}

type UserLogsListRequest struct {
	UserID int64 `json:"userID" form:"userID"`
	commonImport.SqlFilter
}
