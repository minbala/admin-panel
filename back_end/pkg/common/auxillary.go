package commonimport

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func IsStringValid(s string) bool {
	if m := strings.TrimSpace(s); len(m) == 0 || len(s) == 0 {
		return false
	}
	return true
}

func IsIDValid(s int64) bool {
	return s > 0
}

func GetUserLog(c *gin.Context, userID int64, data any) UserLog {
	return UserLog{
		UserID:       userID,
		Event:        c.Request.Method,
		RequestUrl:   c.FullPath(),
		Data:         data,
		Status:       0,
		ErrorMessage: "",
	}
}
