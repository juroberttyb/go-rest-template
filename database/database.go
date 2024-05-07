package database

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/A-pen-app/kickstart/config"
	"github.com/A-pen-app/logging"
)

// DB is the interface handle to a SQL database.
type DB interface {
	initialize(ctx context.Context, cfg dbConfig)
	finalize()
	db() interface{}
}

// dbConfig is an alias for DBContext.
type dbConfig DBContext

// DBContext contains the necessary information to connect to a SQL database.
type DBContext struct {
	// The dialect of the SQL database.
	Dialect string

	// The username used to login to the database.
	Username string

	// The password used to login to the database.
	Password string

	// The address of the database service to connect to.
	Address string

	// The port of the database service to connect to.
	Port string

	// The name of the database to connect to.
	DBName string
}

// Global database interface.
var dbIntf DB

// Database root context.
var dbRootCtx context.Context

// Initialize initializes the database module and database handles.
func Initialize(ctx context.Context) {
	// Initialize PostgreSQL.
	initializePostgres(ctx)

	// Save database root context.
	dbRootCtx = ctx

	// Create database according to dialect.
	dialect := config.GetString("DATABASE_DIALECT")
	switch dialect {
	case "postgres", "cloudsqlpostgres":
		dbIntf = &postgresDB{}
	default:
		panic("invalid dialect")
	}

	// Get database configuration from environment variables.
	cfg := dbConfig{
		Dialect:  config.GetString("DATABASE_DIALECT"),
		Username: config.GetString("DATABASE_USERNAME"),
		Password: config.GetString("DATABASE_PASSWORD"),
		Address:  config.GetString("DATABASE_HOST"),
		Port:     config.GetString("DATABASE_PORT"),
		DBName:   config.GetString("DATABASE_NAME"),
	}

	// Initialize the database context.
	dbIntf.initialize(ctx, cfg)
}

// Finalize finalizes the database module and closes the database handles.
func Finalize() {
	// Finalize PostgreSQL.
	finalizePostgres()

	// Make sure database instance has been initialized.
	if dbIntf == nil {
		panic("database has not been initialized")
	}

	// Finalize database instance.
	dbIntf.finalize()
}

// GetDB returns the database instance.
func GetDB() interface{} {
	return dbIntf.db()
}

// GetSQL returns the SQL database instance.
func GetSQL() *gorm.DB {
	return GetDB().(*gorm.DB)
}

// GetMongoDB returns a MongoDB database instance.
func GetMongoDB() *mongo.Client {
	return GetDB().(*mongo.Client)
}

// DBTransactionFunc is the function pointer type to pass to database
// transaction executor functions.
type DBTransactionFunc func(tx *sqlx.Tx) error

// Transaction executes the provided function as a transaction,
// and automatically performs commit / rollback accordingly.
func Transaction(db *sqlx.DB, txFunc DBTransactionFunc) (err error) {
	// Obtain transaction handle.
	var tx *sqlx.Tx
	if tx, err = db.Beginx(); err != nil {
		logging.Error(dbRootCtx, "Failed to begin transaction: %v", err)
		return err
	}

	// Defer commit / rollback before we execute transaction.
	defer func() {
		// Recover from panic.
		var recovered interface{}
		if recovered = recover(); recovered != nil {
			// Assemble log string.
			message := fmt.Sprintf("\x1b[31m%v\n[Stack Trace]\n%s\x1b[m",
				recovered, debug.Stack())

			// Record the stack trace to logging service.
			logging.Error(dbRootCtx, message)
		}

		// Perform rollback if panic or  is encountered.
		if recovered != nil || err != nil {
			if rerr := tx.Rollback(); rerr != nil {
				logging.Error(dbRootCtx, "Failed to rollback transaction: %v", rerr)
			}
		}
	}()

	// Execute transaction.
	if err = txFunc(tx); err != nil {
		logging.Error(dbRootCtx, "Failed to execute transaction: %v", err)
		return err
	}

	// Commit transaction.
	if err = tx.Commit(); err != nil {
		logging.Error(dbRootCtx, "Failed to commit transaction: %v", err)
		return err
	}

	return nil
}
