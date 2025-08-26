package utils

import (
	"database/sql"
	"fmt"
	"log"

	"library-management-system/constants"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain text password using bcrypt
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CheckPasswordHash compares a password with its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ConnectDB establishes a connection to the database
func ConnectDB() *sql.DB {
	db, err := sql.Open("mysql", constants.DatabaseDSN)
	if err != nil {
		log.Fatal(constants.DatabaseOpenError, err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(constants.DatabaseConnectionError, err)
	}
	fmt.Println("Connected to database")
	return db
}
