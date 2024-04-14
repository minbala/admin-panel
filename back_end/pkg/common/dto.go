package commonimport

import (
	"github.com/golang-jwt/jwt/v5"
)

type SqlFilter struct {
	Limit  int32 `json:"limit" example:"20" form:"limit"`
	Offset int32 `json:"offset" example:"0" form:"offset"`
}

func (s *SqlFilter) Setup() {
	if s.Limit == 0 {
		s.Limit = 20
	}
	if s.Limit == -1 || s.Limit > 1000 {
		s.Limit = 1000
	}
}

type AuthJwtTokenInput struct {
	UserID int64 `json:"userId"`
}

type AuthJwtClaims struct {
	AuthJwtTokenInput
	jwt.RegisteredClaims
}

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e ResponseError) Error() string {
	return e.Message
}

type UserData struct {
	UserID   int64
	UserRole string
}

type UserLog struct {
	UserID       int64  `json:"userId"`
	Event        string `json:"event"`
	RequestUrl   string `json:"requestUrl"`
	Data         any    `json:"data"`
	Status       int32  `json:"status"`
	ErrorMessage string `json:"errorMessage"`
}
