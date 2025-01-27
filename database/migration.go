package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/imanudd/forum-app/config"
	migrate "github.com/rubenv/sql-migrate"
)

const (
	MIGRATION_TYPE_CREATE = "create"
	MIGRATION_TYPE_UP     = "up"
	MIGRATION_TYPE_DOWN   = "down"
	MIGRATION_TYPE_FRESH  = "fresh"
)

type Migration struct {
	cfg *config.Config
	db  *sql.DB
}

func New(cfg *config.Config, db *sql.DB) *Migration {
	return &Migration{
		cfg: cfg,
		db:  db,
	}
}

func (m *Migration) Start(migrationType string) {
	migrations := &migrate.FileMigrationSource{Dir: m.cfg.Service.MigrationDir}

	var direction migrate.MigrationDirection

	switch migrationType {
	case MIGRATION_TYPE_UP:
		direction = migrate.Up
	case MIGRATION_TYPE_DOWN:
		direction = migrate.Down
	case MIGRATION_TYPE_FRESH:
		if m.cfg.Service.Environment == "production" {
			fmt.Print("cannot migrate fresh in production")
			return
		}

		fmt.Println("drop schema !!!")
		_, err := m.db.Exec("drop schema public cascade; create schema public;")
		if err != nil {
			panic(err)
		}

		direction = migrate.Up
	}

	count, err := migrate.Exec(m.db, "mysql", migrations, direction)
	if err != nil {
		panic(err)
	}

	fmt.Printf("applied %d migrations to database\n", count)
}

func (m *Migration) CreateMigrationFile(name string) error {
	migrationsDir := m.cfg.Service.MigrationDir

	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(migrationsDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create migrations directory: %w", err)
		}
	}

	timestamp := time.Now().Format("20060102150405") // Format: YYYYMMDDHHMMSS
	filePath := fmt.Sprintf("%s/%s-%s.sql", migrationsDir, timestamp, name)

	migrationTemplate :=
		`-- +migrate Up
'SQL QUERY' 

-- +migrate Down
'SQL QUERY'
		`

	if err := os.WriteFile(filePath, []byte(migrationTemplate), 0644); err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}

	return nil
}
