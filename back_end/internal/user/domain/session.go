package domain

import (
	"context"
	"time"
)

type SessionRepositoryInterface interface {
	Create(ctx context.Context, session Session) error
	Delete(ctx context.Context, session Session) error
	Get(ctx context.Context, userID int64, accessToken string) (Session, error)
}

type Session struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	AccessToken string    `json:"accessToken"`
	UserID      int64     `json:"userId"`
}
