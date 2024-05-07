package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	// PostgreSQL driver.
	_ "github.com/lib/pq"
	// Google Cloud SQL PostgreSQL driver.
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"

	"github.com/A-pen-app/kickstart/config"
	"github.com/A-pen-app/logging"
)

// Context for PostgreSQL.
var sqlxCtx *DBContext
var sqlxDB *sqlx.DB

// init loads the PostgreSQL database configurations.
func init() {
	// Setup context & load configurations for PostgreSQL.
	// TODO: use different environment names for different databases.
	sqlxCtx = &DBContext{
		Dialect:  config.GetString("DATABASE_DIALECT"),
		Username: config.GetString("DATABASE_USERNAME"),
		Password: config.GetString("DATABASE_PASSWORD"),
		Address:  config.GetString("DATABASE_HOST"),
		Port:     config.GetString("DATABASE_PORT"),
		DBName:   config.GetString("DATABASE_NAME"),
	}
}

// GetPostgres returns the PostgreSQL database handle.
func GetPostgres() *sqlx.DB {
	return sqlxDB
}

// Initialize the PostgreSQL database handle.
func initializePostgres(ctx context.Context) {
	// Make sure the PostgreSQL configurations have been loaded.
	if sqlxCtx == nil {
		panic("database configuration not loaded")
	}

	// Assemble PostgreSQL database source and setup database handle.
	dbSource := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s
		sslmode=disable`, sqlxCtx.Address, sqlxCtx.Port,
		sqlxCtx.Username, sqlxCtx.Password, sqlxCtx.DBName)
	sqlxDB = sqlx.MustConnect(sqlxCtx.Dialect, dbSource)
	sqlxDB.DB.SetMaxOpenConns(maxOpenConns)       // The default is 0 (unlimited)
	sqlxDB.DB.SetMaxIdleConns(maxIdleConns)       // defaultMaxIdleConns = 2
	sqlxDB.DB.SetConnMaxLifetime(maxConnLifetime) // 0, connections are reused forever.

	// Load UUID extension if not loaded.
	stmt := fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	_, err := sqlxDB.Exec(stmt)
	if err != nil {
		panic(err)
	}
}

// Finalize the PostgreSQL database handle.
func finalizePostgres() {
	// Check to see if the PostgreSQL database handle has been initialized.
	if sqlxDB == nil {
		logging.Error(dbRootCtx, "Database handle not initialized")
		return
	}

	// Close the PostgreSQL database handle.
	if err := sqlxDB.Close(); err != nil {
		logging.Error(dbRootCtx, "Failed to close database handle: %v", err)
	}
}
