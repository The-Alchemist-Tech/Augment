package db

import (
	"database/sql"
	"fmt"
	"os"
    "path/filepath"
	"time"
	"context"
	"log"

	fundModel "augment/models"
	
	"github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/mysql"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/go-sql-driver/mysql"
)

var database *sql.DB

type DBArgs struct {
	Host string
	Port string
	Name string
	User string
	Password string
}

const (
	maxReties = 5
	retryDelay = 2 * time.Second
	timeout = 10 * time.Second
)

func CreateDatabase() (*sql.DB, error) {
	DBargs := getDatabaseArgs()

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true",
		DBargs.User, 
		DBargs.Password, 
		DBargs.Host, 
		DBargs.Port, 
		DBargs.Name,
	)

	// Open DB connection
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	// Check DB is ready before running migrations
	for i := 0; i < maxReties; i++ {
		if err = db.Ping(); err == nil {
			// Connected successfully, do migrations
			if migrationErr := migrateDatabase(db); migrationErr != nil {
				return nil, migrationErr
			}

			// Set global database variable and return connection
			database = db
			return db, nil
		}

		// Sleep before retrying connection
		if i < maxReties-1 {
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("Failed to connect to database after %d attempts: %v", maxReties, err)
}

func getDatabaseArgs() (args DBArgs) {
	args.Host = os.Getenv("DB_HOST")
    if args.Host == "" {
        args.Host = "db"
    }

    args.Port = os.Getenv("DB_PORT")
    if args.Port == "" {
        args.Port = "3306"
    }

    args.User = os.Getenv("DB_USER")
    if args.User == "" {
        args.User = "user"
    }

    args.Password = os.Getenv("DB_PASSWORD")
    if args.Password == "" {
        args.Password = "userpassword"
    }

    args.Name = os.Getenv("DB_NAME")
    if args.Name == "" {
        args.Name = "augment"
    }

	return args
}

func migrateDatabase(db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	migrationPath := fmt.Sprintf("file://%s", filepath.Join(dir, "database", "migrations"))

	// Create migration runner
	runner, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"mysql",
		driver,
	)

	if err != nil {
		return err
	}

	// Do migration, ignore "no change" error
	if err := runner.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }

	return nil
}

func InsertFund(name string, units int) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Println("Starting InsertFund transaction")
	tx, err := database.BeginTx(ctx, nil)
    if err != nil { 
		log.Printf("Failed to begin transaction: %v", err)
		return 0,err 
	}
    defer tx.Rollback()

	// TODO: Check if a fund with the same name already exists and return an error if so

	query := "INSERT INTO funds (name, units) VALUES (?, ?)"

	res, err := tx.ExecContext(ctx, query, name, units)
	if err != nil {
		log.Printf("Failed to execute insert query: %v", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Failed to get last inserted ID: %v", err)
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return 0, err
	}

	return id, nil
}

func GetResourceByID(table string, id int64) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Printf("Retrieving resource from table %s with ID: %d", table, id)

	switch table {
	case "funds":
		log.Printf("SWITCH: Querying funds table for ID %d", id)
		// Limit to fields we want in real world, but * works for this simple app
		query := "SELECT * FROM funds WHERE id = ?"
		row := database.QueryRowContext(ctx, query, id)

		f := &fundModel.Fund{}

		if err := row.Scan(&f.ID, &f.Name, &f.Units, &f.CreatedOn, &f.LastModified); err != nil {
			log.Printf("Failed to scan row for ID %d: %v", id, err)
			return nil, err
		}

		return f, nil
	default:
		log.Printf("Invalid table name: %s", table)
		return nil, fmt.Errorf("Invalid table name: %s", table)
	}

	log.Printf("Resource not found in table %s with ID %d", table, id)
	return nil, fmt.Errorf("Resource not found in table %s with ID %d", table, id)
}
