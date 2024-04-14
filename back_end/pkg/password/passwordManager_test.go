package passwordpackage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPasswordHash(t *testing.T) {
	passwordManager := ProvidePasswordManager()
	password := "minbala33"
	hashPassword, err := passwordManager.Hash(password)
	assert.NoError(t, err)
	err = passwordManager.CheckPassword(password, hashPassword)
	assert.NoError(t, err)

	wrongPasword := "minbala333"
	err = passwordManager.CheckPassword(wrongPasword, hashPassword)
	assert.Error(t, err)
}
