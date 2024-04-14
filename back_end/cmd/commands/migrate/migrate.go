package migrate

import (
	commonImport "admin-panel/pkg/common"
	"admin-panel/pkg/postgres"
	"github.com/spf13/cobra"
	"log"
)

type MigrateCommand struct {
	container *commonImport.Container
}

func (m *MigrateCommand) RunE(cmd *cobra.Command, args []string) error {
	postgres.ProvideDBWithInitialData(m.container.DBClient, m.container.DB, m.container.PasswordManager.Hash)
	log.Println("data filling has been completed")
	return nil
}

func (m *MigrateCommand) NewCommand(container *commonImport.Container) *cobra.Command {
	m.container = container
	migrateCommand := cobra.Command{
		Use:   "migrate",
		Short: "Migrate",
		Long:  "Migrate database with initial data",
		RunE:  m.RunE,
	}

	return &migrateCommand
}
