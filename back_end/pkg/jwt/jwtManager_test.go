package jwt

import (
	commonImport "admin-panel/pkg/common"
	"admin-panel/pkg/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseToken(t *testing.T) {
	configuration := &config.Config{
		App: config.App{
			JWTSECRET:              "minbala33",
			AccessTokenExpiredTime: 1,
		},
		Database: config.Database{},
	}
	var userID int64 = 1
	parser := ProvideTokenParser(configuration)
	token, err := parser.CreateAuthToken(commonImport.AuthJwtTokenInput{
		UserID: userID,
	})
	assert.NoError(t, err)

	data, err := parser.ParseAuthJwtToken(token, configuration.App.JWTSECRET)
	assert.NoError(t, err)
	assert.Equal(t, data.UserID, userID)

	_, err = parser.ParseAuthJwtToken("token", configuration.App.JWTSECRET)
	assert.Error(t, err)
}
