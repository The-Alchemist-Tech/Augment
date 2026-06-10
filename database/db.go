package db

import (
	"database/sql"
	"fmt"
	"os"
    "path/filepath"
	"time"
	
	"github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/mysql"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/go-sql-driver/mysql"
)

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
)

func CreateDatabase() (*sql.DB, error) {
	DBargs := getDatabaseArgs()

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true",
		DBargs.User, 
		DBargs.Password, 
		DBargs.Host, 
		DBargs.Port, 
		DBargs.Name
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