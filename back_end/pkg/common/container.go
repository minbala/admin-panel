package commonimport

import (
	"admin-panel/pkg/config"
	"admin-panel/pkg/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
)

type Container struct {
	Logger          *Logger
	Operation       *HTTPOperation
	Config          *config.Config
	Authenticator   Authenticator
	DB              *pgxpool.Pool
	DBClient        *postgres.Queries
	PasswordManager PasswordManager
}

var container *Container
var setUpContainerOnce sync.Once

func ProvideContainer(logger *Logger, operation *HTTPOperation, config *config.Config, authenticator Authenticator,
	db *pgxpool.Pool, client *postgres.Queries, passwordManager PasswordManager) *Container {
	setUpContainerOnce.Do(func() {
		container = &Container{
			Logger:          logger,
			Operation:       operation,
			Config:          config,
			Authenticator:   authenticator,
			DB:              db,
			DBClient:        client,
			PasswordManager: passwordManager,
		}
	})
	return container
}
