package domain

import (
	"context"
	"time"
)

type UserLogRepositoryInterface interface {
	Create(ctx context.Context, log UserLog) error
	Get(ctx context.Context, query UserLogQuery) ([]UserLog, error)
	Count(ctx context.Context, query UserLogQuery) (int64, error)
}

type UserLog struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	UserID       int64     `json:"userId"`
	Event        string    `json:"event"`
	RequestUrl   string    `json:"requestUrl"`
	Data         []byte    `json:"data"`
	Status       int32     `json:"status"`
	ErrorMessage string    `json:"errorMessage"`
}

type UserLogQuery struct {
	UserID int64 `json:"userId"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}
