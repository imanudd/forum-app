package database

import (
	"database/sql"
	"fmt"

	"github.com/imanudd/forum-app/config"
	migrate "github.com/rubenv/sql-migrate"
)

const (
	MIGRATION_TYPE_UP    = "up"
	MIGRATION_TYPE_DOWN  = "down"
	MIGRATION_TYPE_FRESH = "fresh"
)

type Migration struct {
	cfg *config.MainConfig
	db  *sql.DB
}

func New(cfg *config.MainConfig, db *sql.DB) *Migration {
	return &Migration{
		cfg: cfg,
		db:  db,
	}
}

func (m *Migration) Start(migrationType string) {
	migrations := &migrate.FileMigrationSource{Dir: "./database/migration"}

	var direction migrate.MigrationDirection

	switch migrationType {
	case MIGRATION_TYPE_UP:
		direction = migrate.Up
	case MIGRATION_TYPE_DOWN:
		direction = migrate.Down
	case MIGRATION_TYPE_FRESH:
		if m.cfg.Environment == "production" {
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

	count, err := migrate.Exec(m.db, "postgres", migrations, direction)
	if err != nil {
		panic(err)
	}
	fmt.Printf("applied %d migrations to database\n", count)
}
