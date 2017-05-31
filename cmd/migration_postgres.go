// +build postgres

package cmd

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

func CreateMigration(url string) (*migrate.Migrate, error) {
	// todo, get this from an argument, variable or config file
	db, err := sql.Open("postgres", url)
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://.",
		"postgres", driver)
	if err != nil {
		return nil, err
	}

	return m, nil
}
