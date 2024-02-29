package migrations

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

const (
	UP = iota
	DOWN
)

// Migrate sets up the db tables and functions
func Migrate(MigrationFilePath string, URL string, mode int) {

	// log start
	log.Println("Start migrations...")

	sourceUrl := fmt.Sprintf("file://%s", MigrationFilePath)

	m, err := migrate.New(sourceUrl, URL)
	if err != nil {
		panic(err)
	}

	if mode == UP {

		if err := m.Up(); err != nil {
			log.Println("error in running migrations:", err)
			return
		}
	} else {

		if err := m.Down(); err != nil {
			log.Println("error in resetting migrations:", err)
			return
		}
	}

	log.Println("Completed migrations")
}
