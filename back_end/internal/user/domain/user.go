package domain

import (
	"admin-panel/pkg/resources"
	"context"
	"time"
)

type UserRepositoryInterface interface {
	Get(ctx context.Context, query UserQuery) ([]User, error)
	Count(ctx context.Context, query UserQuery) (int64, error)
	Create(ctx context.Context, user User) (int64, error)
	Delete(ctx context.Context, userID int64) error
	Update(ctx context.Context, user User) error
}

type User struct {
	ID        int64                  `json:"id"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
	Email     string                 `json:"email"`
	Name      string                 `json:"name"`
	Password  string                 `json:"password"`
	UserRole  resources.UserRoleType `json:"userRole"`
}

type UserQuery struct {
	ID       int64                  `json:"id"`
	Email    string                 `json:"email"`
	UserRole resources.UserRoleType `json:"userRole"`
	Name     string                 `json:"name"`
	Limit    int32                  `json:"limit"`
	Offset   int32                  `json:"offset"`
}
