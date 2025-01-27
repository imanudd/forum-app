package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "clean-arch-pattern",
	Short: "Clean Arch Pattern",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use serve to start a server")
		fmt.Println("Use -h to see the list of command")
	},
}

func Run() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateFreshCmd)
	migrateCmd.AddCommand(migrateCreateCmd)

	rootCommand.AddCommand(migrateCmd)
	rootCommand.AddCommand(restCommand)

	if err := rootCommand.Execute(); err != nil {
		log.Fatal(err)
	}

}
