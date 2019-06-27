package dbqueries

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	db *sqlx.DB
)

func InitDB() {
	var err error
	db, err = sqlx.Connect("postgres", "postgres://docker:docker@127.0.0.1:5432/docker?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Panic(err)
	}

	doMigrate()
}
