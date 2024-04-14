// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CreateSession(ctx context.Context, arg CreateSessionParams) error
	CreateUser(ctx context.Context, arg CreateUserParams) (int64, error)
	CreateUserLog(ctx context.Context, arg CreateUserLogParams) error
	DeleteSession(ctx context.Context, arg DeleteSessionParams) error
	DeleteUser(ctx context.Context, id int64) error
	GetSession(ctx context.Context, arg GetSessionParams) (Sessions, error)
	GetUserLogs(ctx context.Context, arg GetUserLogsParams) ([]UserLogs, error)
	GetUserLogsCount(ctx context.Context, userID pgtype.Int8) (int64, error)
	GetUsers(ctx context.Context, arg GetUsersParams) ([]Users, error)
	GetUsersCount(ctx context.Context, arg GetUsersCountParams) (int64, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) error
}

var _ Querier = (*Queries)(nil)
