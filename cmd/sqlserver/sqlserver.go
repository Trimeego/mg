package sqlserver

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	nurl "net/url"

	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database"
)

func init() {
	db := SqlServer{}
	database.Register("mssql", &db)
}

var DefaultMigrationsTable = "SYC_SCHEMA_MIGRATIONS"

var (
	ErrNilConfig      = fmt.Errorf("no config")
	ErrNoDatabaseName = fmt.Errorf("no database name")
	ErrDatabaseDirty  = fmt.Errorf("database is dirty")
)

type Config struct {
	MigrationsTable string
	DatabaseName    string
}

type SqlServer struct {
	db       *sql.DB
	isLocked bool

	// Open and WithInstance need to garantuee that config is never nil
	config *Config
}

func WithInstance(instance *sql.DB, config *Config) (database.Driver, error) {
	if config == nil {
		return nil, ErrNilConfig
	}

	if err := instance.Ping(); err != nil {
		return nil, err
	}

	query := `select DB_NAME()`
	var databaseName string
	if err := instance.QueryRow(query).Scan(&databaseName); err != nil {
		return nil, &database.Error{OrigErr: err, Query: []byte(query)}
	}

	if len(databaseName) == 0 {
		return nil, ErrNoDatabaseName
	}

	config.DatabaseName = databaseName

	if len(config.MigrationsTable) == 0 {
		config.MigrationsTable = DefaultMigrationsTable
	}

	px := &SqlServer{
		db:     instance,
		config: config,
	}

	if err := px.ensureVersionTable(); err != nil {
		return nil, err
	}

	return px, nil
}

func (s *SqlServer) Open(url string) (database.Driver, error) {
	purl, err := nurl.Parse(url)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("mssql", migrate.FilterCustomQuery(purl).String())
	if err != nil {
		return nil, err
	}

	migrationsTable := purl.Query().Get("x-migrations-table")
	if len(migrationsTable) == 0 {
		migrationsTable = DefaultMigrationsTable
	}

	px, err := WithInstance(db, &Config{
		DatabaseName:    purl.Path,
		MigrationsTable: migrationsTable,
	})
	if err != nil {
		return nil, err
	}

	return px, nil
}

func (s *SqlServer) Close() error {
	return s.db.Close()
}

// https://www.postgresql.org/docs/9.6/static/explicit-locking.html#ADVISORY-LOCKS
func (s *SqlServer) Lock() error {

	if s.isLocked {
		return database.ErrLocked
	}

	s.isLocked = true

	return nil

}

func (s *SqlServer) Unlock() error {

	if !s.isLocked {
		return nil
	}

	s.isLocked = false
	return nil
}

func (s *SqlServer) Run(migration io.Reader) error {
	migr, err := ioutil.ReadAll(migration)
	if err != nil {
		return err
	}

	// run migration
	query := string(migr[:])
	if _, err := s.db.Exec(query); err != nil {
		// TODO: cast to postgress error and get line number
		return database.Error{OrigErr: err, Err: "migration failed", Query: migr}
	}

	return nil
}

func (s *SqlServer) SetVersion(version int, dirty bool) error {
	tx, err := s.db.Begin()
	if err != nil {
		return &database.Error{OrigErr: err, Err: "transaction start failed"}
	}

	query := `TRUNCATE TABLE ` + s.config.MigrationsTable
	if _, err := s.db.Exec(query); err != nil {
		tx.Rollback()
		return &database.Error{OrigErr: err, Query: []byte(query)}
	}

	if version >= 0 {
		query = `INSERT INTO "` + s.config.MigrationsTable + `" (version, dirty) VALUES ($1, $2)`
		if _, err := s.db.Exec(query, version, dirty); err != nil {
			tx.Rollback()
			return &database.Error{OrigErr: err, Query: []byte(query)}
		}
	}

	if err := tx.Commit(); err != nil {
		return &database.Error{OrigErr: err, Err: "transaction commit failed"}
	}

	return nil
}

func (s *SqlServer) Version() (version int, dirty bool, err error) {
	query := `SELECT top 1 version, dirty FROM ` + s.config.MigrationsTable
	err = s.db.QueryRow(query).Scan(&version, &dirty)
	if err != nil {
		return database.NilVersion, false, nil
	}

	return version, dirty, nil
}

func (s *SqlServer) Drop() error {
	return fmt.Errorf("Drop is not supported for Loves")
}

func (s *SqlServer) ensureVersionTable() error {
	// check if migration table exists
	var count int
	query := `Select count(*) from sys.tables WHERE name = $1`
	if err := s.db.QueryRow(query, s.config.MigrationsTable).Scan(&count); err != nil {
		return &database.Error{OrigErr: err, Query: []byte(query)}
	}
	if count == 1 {
		return nil
	}

	// if not, create the empty migration table
	query = `CREATE TABLE "` + s.config.MigrationsTable + `" (version bigint not null primary key, dirty bit not null)`
	if _, err := s.db.Exec(query); err != nil {
		return &database.Error{OrigErr: err, Query: []byte(query)}
	}
	return nil
}
