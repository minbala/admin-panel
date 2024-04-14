package postgres

import (
	"github.com/jackc/pgx/v5/pgtype"
	"strings"
)

func IsStringValid(s string) bool {
	if m := strings.TrimSpace(s); len(m) == 0 || len(s) == 0 {
		return false
	}
	return true
}

// StringToNullString converts *string to sql.NullString.
func StringToNullString(s string) pgtype.Text {
	if IsStringValid(s) {
		return pgtype.Text{String: s, Valid: true}
	}
	return pgtype.Text{}
}

func Int64ToPgxInt(s int64) pgtype.Int8 {
	if s != 0 {
		return pgtype.Int8{Int64: s, Valid: true}
	}
	return pgtype.Int8{}
}

func ContainQuery(query string) string {
	return "%" + query + "%"
}
