// Package utils provides utility functions for the Library Management System.
// This package contains reusable functions for password hashing, database
// connections, and other common operations that are used across multiple
// parts of the application.
package utils

import (
	"fmt"
	"log"

	"library-management-system/constants"
	"library-management-system/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// HashPassword securely hashes a plain text password using bcrypt algorithm.
// bcrypt is a preferred password hashing function as it includes a salt
// and is computationally expensive to make brute force attacks difficult.
//
// Parameters:
//   - password: The plain text password to hash
//
// Returns:
//   - string: The hashed password as a string
//   - error: Any error that occurred during hashing
//
// Usage:
//
//	hashedPwd, err := HashPassword("mypassword123")
func HashPassword(password string) (string, error) {
	// Use bcrypt's default cost (currently 10) for good security/performance balance
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CheckPasswordHash compares a plain text password with its bcrypt hash.
// This function is used during login to verify if the provided password
// matches the stored hash.
//
// Parameters:
//   - password: The plain text password to verify
//   - hash: The bcrypt hash to compare against
//
// Returns:
//   - bool: true if password matches hash, false otherwise
//
// Usage:
//
//	isValid := CheckPasswordHash("mypassword123", storedHash)
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// If no error, passwords match; if error, they don't match
	return err == nil
}

// ConnectDB establishes and returns a GORM database connection to MySQL.
// This function handles the initial connection setup, validates connectivity,
// and performs automatic database migration for all models.
// GORM provides built-in protection against SQL injection attacks.
//
// Returns:
//   - *gorm.DB: A pointer to the GORM database connection
//
// Note: GORM manages connection pooling automatically. The connection
// should be closed when the application shuts down.
//
// Usage:
//
//	db := ConnectDB()
//	// GORM handles connection management automatically
func ConnectDB() *gorm.DB {
	// Open GORM connection using MySQL driver with DSN from constants
	db, err := gorm.Open(mysql.Open(constants.DatabaseDSN), &gorm.Config{})
	if err != nil {
		// Fatal error will terminate the application - appropriate for startup failures
		log.Fatal(constants.DatabaseOpenError, err)
	}

	// Get the underlying SQL DB to perform ping test
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(constants.DatabaseConnectionError, err)
	}

	// Ping verifies the connection is actually established
	if err := sqlDB.Ping(); err != nil {
		log.Fatal(constants.DatabaseConnectionError, err)
	}

	// Perform automatic migration for all models
	// This creates/updates tables based on struct definitions
	err = db.AutoMigrate(
		&model.Student{},
		&model.Librarian{},
		&model.Book{},
		&model.BorrowedBook{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database tables:", err)
	}

	// Connection successful - log for monitoring/debugging purposes
	fmt.Println("Connected to database with GORM")
	fmt.Println("Database migration completed successfully")
	return db
}
