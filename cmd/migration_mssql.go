// +build mssql

package cmd

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/source/file"
	"github.com/trimeego/mg/cmd/sqlserver"
)

func CreateMigration(url string) (*migrate.Migrate, error) {
	// todo, get this from an argument, variable or config file
	db, err := sql.Open("mssql", url)
	if err != nil {
		return nil, err
	}

	driver, err := sqlserver.WithInstance(db, &sqlserver.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance("file://.", "mssql", driver)
	if err != nil {
		return nil, err
	}

	return m, nil
}
