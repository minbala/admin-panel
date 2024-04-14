package service

import (
	"admin-panel/internal/user/domain"
	"admin-panel/pkg/resources"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Get(ctx context.Context, query domain.UserQuery) ([]domain.User, error) {
	args := m.Called(ctx, query)
	var output []domain.User
	if args.Get(0) != nil {
		if m, ok := args.Get(0).([]domain.User); ok {
			output = m
		}
	}
	return output, args.Error(1)
}

func (m *MockUserRepo) Count(ctx context.Context, query domain.UserQuery) (int64, error) {
	args := m.Called(ctx, query)
	var output int64
	if args.Get(0) != nil {
		if m, ok := args.Get(0).(int64); ok {
			output = m
		}
	}
	return output, args.Error(1)
}

func (m *MockUserRepo) Create(ctx context.Context, user domain.User) (int64, error) {
	args := m.Called(ctx, user)
	var output int64
	if args.Get(0) != nil {
		if m, ok := args.Get(0).(int64); ok {
			output = m
		}
	}
	return output, args.Error(1)
}

func (m *MockUserRepo) Delete(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepo) Update(ctx context.Context, user domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

type MockPasswordManager struct {
	mock.Mock
}

func (m *MockPasswordManager) CheckPassword(password string, hashPassword string) error {
	args := m.Called(password, hashPassword)
	return args.Error(0)
}
func (m *MockPasswordManager) Hash(plain string) (string, error) {
	args := m.Called(plain)
	var output string
	if args.Get(0) != nil {
		if m, ok := args.Get(0).(string); ok {
			output = m
		}
	}
	return output, args.Error(1)
}

// need to add mock jwt manager, session manager, userLog manager
func TestCreateUser(t *testing.T) {
	mockUserRepo := new(MockUserRepo)
	mockPasswordManager := new(MockPasswordManager)
	ctx := context.Background()
	mockPasswordManager.On("Hash", "minbala33").Return("minbala3333", nil)
	mockUserRepo.On("Create", context.Background(), mock.Anything).Return(int64(1), nil)
	createdAt := time.Now()

	userService := ProvideUserService(mockUserRepo, mockPasswordManager, nil, nil, nil, nil)
	userCreateCmd := UserCreateCmd{
		Email:    "minbala33@gmail.com",
		Name:     "minbala",
		Password: "minbala33",
		UserRole: "admin",
	}
	userID, err := userService.CreateUser(ctx, userCreateCmd)
	assert.NoError(t, err)

	mockUserCall := mockUserRepo.On("Get", mock.Anything, mock.Anything).Return([]domain.User{{Name: "minbala",
		Email: "minbala33@gmail.com", ID: 1, CreatedAt: createdAt,
		UserRole: "admin", Password: "minbala3333"}}, nil)
	assert.NoError(t, err)

	user, err := userService.GetUser(ctx, domain.UserQuery{ID: userID})
	assert.NoError(t, err)

	hashPassword, err := userService.passwordManager.Hash(userCreateCmd.Password)
	assert.NoError(t, err)
	assert.Equal(t, domain.User{
		ID:        userID,
		Name:      userCreateCmd.Name,
		Email:     userCreateCmd.Email,
		Password:  hashPassword,
		UserRole:  userCreateCmd.UserRole,
		CreatedAt: createdAt,
	}, user)

	mockUserRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	err = userService.UpdateUser(ctx, UserUpdateCmd{UserID: userID, Name: "minbala", Password: "",
		UserRole: "user", Email: "minbala22@gmail.com"})
	assert.NoError(t, err)

	mockUserCall.Unset()
	mockUserCall = mockUserRepo.On("Get", mock.Anything, mock.Anything).Return([]domain.User{{Name: "minbala", Email: "minbala22@gmail.com", ID: userID, CreatedAt: createdAt,
		UserRole: "user", Password: "minbala3333"}}, nil)

	user, err = userService.GetUser(ctx, domain.UserQuery{ID: userID})
	assert.NoError(t, err)

	hashPassword, err = userService.passwordManager.Hash(userCreateCmd.Password)
	assert.NoError(t, err)
	assert.Equal(t, domain.User{
		ID:        userID,
		Name:      "minbala",
		Email:     "minbala22@gmail.com",
		Password:  hashPassword,
		UserRole:  "user",
		CreatedAt: createdAt,
	}, user)

	mockUserRepo.On("Delete", mock.Anything, mock.Anything).Return(nil)
	err = userService.DeleteUser(ctx, userID)
	assert.NoError(t, err)
	mockUserCall.Unset()
	mockUserRepo.On("Get", mock.Anything, mock.Anything).Return(nil, resources.ErrNoRecordFound)
	err = userService.DeleteUser(ctx, userID)
	assert.Error(t, err)
}
