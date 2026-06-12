package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	models "augment/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/shopspring/decimal"
)

var database *sql.DB

// Error for duplicate fund name so the handler can return the proper status code
var ErrDuplicate = errors.New("duplicate object found")
var ErrDuplicateEmail = errors.New("duplicate email found")
var ErrDuplicateUsername = errors.New("duplicate username found")

type DBArgs struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

const (
	maxReties  = 5
	retryDelay = 2 * time.Second
	timeout    = 10 * time.Second
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

func InsertResource(resource any) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Println("Starting InsertResource transaction")
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return 0, err
	}
	defer tx.Rollback()

	switch r := resource.(type) {
	case *models.Fund:
		log.Println("Starting InsertFund transaction")
		return InsertFundTx(ctx, tx, r)
	case *models.Investor:
		log.Println("Starting InsertInvestor transaction")
		return InsertInvestorTx(ctx, tx, r)
	case *models.Cap:
		log.Println("Starting InsertCap transaction")
		return InsertCapTx(ctx, tx, r)
	default:
		log.Printf("Invalid resource type: %T", resource)
		return 0, fmt.Errorf("Invalid resource type: %T", resource)
	}
}

func InsertFundTx(ctx context.Context, tx *sql.Tx, fund *models.Fund) (int64, error) {
	// Check if a fund with the same name already exists and return an error if so
	var existingID int64
	err := tx.QueryRowContext(ctx, "SELECT id FROM funds WHERE name = ?", fund.Name).Scan(&existingID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Failed to check for existing fund: %v", err)
		return 0, err
	}

	// Existing fund found with same name found, return error
	if err == nil {
		// Custom error to check in fund handler to return proper status and message.
		log.Printf("Duplicate fund name detected: %s", fund.Name)
		return 0, ErrDuplicate
	}

	query := "INSERT INTO funds (name, units) VALUES (?, ?)"

	res, err := tx.ExecContext(ctx, query, fund.Name, *fund.Units)
	if err != nil {
		log.Printf("Failed to execute fund insert: %v", err)
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

func InsertInvestorTx(ctx context.Context, tx *sql.Tx, investor *models.Investor) (int64, error) {
	// TODO: With more time, I would have moved the following checks into the base investor.go file
	// Check if an investor with the same email already exists and return an error if so
	var existingEmail string
	var existingUsername string
	err := tx.QueryRowContext(ctx, "SELECT email, username FROM investors WHERE email = ? OR username = ?", investor.Email, investor.Username).Scan(&existingEmail, &existingUsername)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Failed to check for existing investor: %v", err)
		return 0, err
	}

	// Existing investor found with same email or username found, return error
	if err == nil {
		if existingEmail == investor.Email {
			log.Printf("Duplicate investor email detected: %s", investor.Email)
			return 0, ErrDuplicateEmail
		} else if existingUsername == investor.Username {
			log.Printf("Duplicate investor username detected: %s", investor.Username)
			// Custom error to check in investor handler to return proper status and message.
			return 0, ErrDuplicateUsername
		}
		// Should not get here, but error out just in case since we found something in the DB
		log.Printf("Duplicate investor detected with email: %s or username: %s", investor.Email, investor.Username)
		return 0, ErrDuplicate
	}

	query := "INSERT INTO investors (username, email, first_name, last_name) VALUES (?, ?, ?, ?)"

	res, err := tx.ExecContext(ctx, query, investor.Username, investor.Email, investor.FirstName, investor.LastName)
	if err != nil {
		log.Printf("Failed to execute investor insert: %v", err)
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

func InsertCapTx(ctx context.Context, tx *sql.Tx, cap *models.Cap) (int64, error) {

	query := "INSERT INTO cap (fund, buyer, seller, units) VALUES (?, ?, ?, ?)"

	res, err := tx.ExecContext(ctx, query, cap.Fund, cap.Buyer, cap.Seller, cap.Units)
	if err != nil {
		log.Printf("Failed to execute cap insert: %v", err)
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
		log.Printf("Querying funds table for ID %d", id)
		// Limit to fields we want in real world, but * works for this simple app
		query := "SELECT * FROM funds WHERE id = ?"
		row := database.QueryRowContext(ctx, query, id)

		f := &models.Fund{}

		// Avoid nil pointer on units which is a pointer in the model
		f.Units = new(int)

		if err := row.Scan(&f.ID, &f.Name, f.Units, &f.CreatedOn, &f.LastModified); err != nil {
			log.Printf("Failed to scan row for ID %d: %v", id, err)
			return nil, err
		}

		return f, nil
	case "investors":
		log.Printf("Querying investors table for ID %d", id)
		// Limit to fields we want in real world, but * works for this simple app
		query := "SELECT * FROM investors WHERE id = ?"
		row := database.QueryRowContext(ctx, query, id)

		i := &models.Investor{}

		if err := row.Scan(&i.ID, &i.Username, &i.Email, &i.FirstName, &i.LastName, &i.CreatedOn, &i.LastModified); err != nil {
			log.Printf("Failed to scan row for ID %d: %v", id, err)
			return nil, err
		}

		return i, nil
	case "cap":
		log.Printf("Querying cap table for ID %d", id)
		// Limit to fields we want in real world, but * works for this simple app
		query := "SELECT * FROM cap WHERE id = ?"
		row := database.QueryRowContext(ctx, query, id)

		c := &models.Cap{}

		// Avoid nil pointer errors on pointers in the model
		c.Fund = new(int)
		c.Buyer = new(int)
		c.Seller = new(int)
		c.Units = new(decimal.Decimal)

		if err := row.Scan(&c.ID, c.Fund, c.Buyer, c.Seller, c.Units, &c.CreatedOn, &c.LastModified); err != nil {
			log.Printf("Failed to scan row for ID %d: %v", id, err)
			return nil, err
		}

		return c, nil
	default:
		log.Printf("Invalid table name: %s", table)
		return nil, fmt.Errorf("Invalid table name: %s", table)
	}
}

func GetFundUnitsForInvestor(fund int64, investorID int64) (decimal.Decimal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Aggregate in sql
	query := `
		SELECT COALESCE(SUM(CASE WHEN buyer = ? THEN units ELSE -units END), 0)
		FROM cap
		WHERE fund = ? AND ? IN (buyer, seller)`

	row := database.QueryRowContext(ctx, query, investorID, fund, investorID)

	var netUnits decimal.Decimal

	err := row.Scan(&netUnits)
	if err != nil {
		log.Printf("Failed to get net units for investor %d in fund %d: %v", investorID, fund, err)
		return decimal.Zero, err
	}

	return netUnits, nil
}
