package cmd

import (
	"admin-panel/cmd/commands"
	"admin-panel/cmd/commands/migrate"
	commonImport "admin-panel/pkg/common"
	"github.com/spf13/cobra"
)

var Commands = []func(container *commonImport.Container) *cobra.Command{
	(&commands.ServeCommand{}).NewCommand,
	(&migrate.MigrateCommand{}).NewCommand,
}

var rootCmd = &cobra.Command{}

func Run(container *commonImport.Container) error {
	for _, command := range Commands {
		rootCmd.AddCommand(command(container))
	}
	return rootCmd.Execute()
}
