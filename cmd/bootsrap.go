package cmd

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/imanudd/forum-app/config"
	_ "github.com/lib/pq"
)

func NewMysql(cfg *config.Config) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.Db.Username, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.DbName))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
