package resources

import "github.com/cockroachdb/errors"

type UserRoleType string

const (
	UserRoleTypeAdmin      UserRoleType = "admin"
	UserRoleTypeNormalUser UserRoleType = "normal user"
)

func UserRoleTypeValid(u string) error {
	if u == string(UserRoleTypeNormalUser) || u == string(UserRoleTypeAdmin) {
		return nil
	}
	return errors.Wrap(ErrClient, "wrong user role type")
}
