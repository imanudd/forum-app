package cmd

import (
	"log"

	"github.com/imanudd/forum-app/config"
	migration "github.com/imanudd/forum-app/database"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use: "migration",
	Run: func(_ *cobra.Command, _ []string) {
		log.Println("use -h to show available commands")
	},
}

var migrateUpCmd = &cobra.Command{
	Use: "up",
	Run: func(_ *cobra.Command, _ []string) {
		startMigrate("up")
	},
}

var migrateDownCmd = &cobra.Command{
	Use: "down",
	Run: func(_ *cobra.Command, _ []string) {
		startMigrate("down")
	},
}

var migrateFreshCmd = &cobra.Command{
	Use: "fresh",
	Run: func(_ *cobra.Command, _ []string) {
		startMigrate("fresh")
	},
}

var migrateCreateCmd = &cobra.Command{
	Use:  "create [filename]",
	Args: cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		startMigrate("create", args[0])
	},
}

func startMigrate(migrationType string, fileName ...string) {
	if err := config.InitConfig(); err != nil {
		log.Fatalln(" failed to initialize config", err.Error())
	}

	cfg := config.Get()

	db := NewMysql(cfg)

	m := migration.New(cfg, db)

	if migrationType == "create" {
		m.CreateMigrationFile(fileName[0])
		return
	}

	m.Start(migrationType)
}
