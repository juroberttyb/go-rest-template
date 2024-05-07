package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	// PostgreSQL driver.
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// Google Cloud SQL PostgreSQL driver.
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"

	_ "github.com/lib/pq"
	//"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/google_cloud_storage"

	"github.com/A-pen-app/kickstart/config"
	"github.com/A-pen-app/logging"
)

// postgresDB is the concrete PostgresSQL handle to a SQL database.
type postgresDB struct{ *gorm.DB }

// Connection pool configuration
var maxIdleConns int
var maxOpenConns int
var maxConnLifetime time.Duration

func init() {
	maxIdleConns = config.GetInt("DATABASE_MAX_IDLE_CONNECTIONS")
	maxOpenConns = config.GetInt("DATABASE_MAX_OPEN_CONNECTIONS")
	maxConnLifetime = config.GetMilliseconds("DATABASE_MAX_CONN_LIFETIME_MS")
}

// initialize initializes the PostgreSQL database handle.
func (db *postgresDB) initialize(ctx context.Context, cfg dbConfig) {
	// Assemble PostgreSQL database source and setup database handle.
	dbSource := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s
		sslmode=disable binary_parameters=yes`, cfg.Address, cfg.Port,
		cfg.Username, cfg.Password, cfg.DBName)

	// Connect to the PostgreSQL database.
	var err error
	db.DB, err = gorm.Open(cfg.Dialect, dbSource)
	if err != nil {
		panic(err)
	}

	// Configure connection pool.
	db.DB.DB().SetMaxIdleConns(maxIdleConns)
	db.DB.DB().SetMaxOpenConns(maxOpenConns)
	db.DB.DB().SetConnMaxLifetime(maxConnLifetime)

	// Load UUID extension if not loaded.
	stmt := fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	if err = db.Exec(stmt).Error; err != nil {
		panic(err)
	}

	// Load crosstab extension if not loaded.
	stmt = fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS \"tablefunc\"")
	if err = db.Exec(stmt).Error; err != nil {
		panic(err)
	}
	// do not run migrate when running tests, otherwise, the migration file path "file://database/migrations" will fail the integration tests
	db.migrate()
}

func (db *postgresDB) migrate() {
	driver, err := postgres.WithInstance(db.DB.DB(), &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres",
		driver)
	if err != nil {
		logging.Error(dbRootCtx, "Failed to create migration instance: %v", err)
		panic(err)
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
	defer m.Close()
}

// finalize finalizes the PostgreSQL database handle.
func (db *postgresDB) finalize() {
	// Close the PostgreSQL database handle.
	if err := db.Close(); err != nil {
		logging.Error(dbRootCtx, "Failed to close database handle: %v", err)
	}
}

// db returns the PostgreSQL GORM database handle.
func (db *postgresDB) db() interface{} {
	return db.DB
}
