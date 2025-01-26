package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/imanudd/forum-app/config"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitElastic(cfg *config.MainConfig) *elasticsearch.Client {

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			cfg.ElasticHost,
		},
		Username:               cfg.ElasticUsername,
		Password:               cfg.ElasticPassword,
		CertificateFingerprint: cfg.ElasticCAFingerprint,
	})

	if err != nil {
		panic(err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer res.Body.Close()
	log.Println(res)

	log.Println("succesfully connected to elastic")

	return es

}

func NewPostgres(cfg *config.MainConfig) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUsername,
		cfg.PostgresPassword,
		cfg.DBName,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

// init postgress with gorm
func InitPostgreSQL(cfg *config.MainConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Jakarta", cfg.PostgresHost, cfg.PostgresUsername, cfg.PostgresPassword, cfg.DBName, cfg.PostgresPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatalf(err.Error())
		panic(err)
	}

	log.Printf("Successfully connected to database server")

	rdb, err := db.DB()
	if err != nil {
		log.Fatalf(err.Error())
		panic(err)
	}

	rdb.SetMaxIdleConns(cfg.MaxIdleConns)
	rdb.SetMaxOpenConns(cfg.MaxOpenConns)
	rdb.SetConnMaxLifetime(time.Duration(int(time.Minute) * cfg.ConnMaxLifetime))

	return db
}
