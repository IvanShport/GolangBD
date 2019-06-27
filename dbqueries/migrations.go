package dbqueries

import (
	"log"

	"github.com/rubenv/sql-migrate"
)

func doMigrate() {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Println(err)
	} else if n != 0 {
		log.Printf("Applied %d migrations!\n", n)
	}
}
