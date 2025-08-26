package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// Student model
type Student struct {
	ID       int    `json:"id"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required,len=10,numeric"`
	Password string `json:"-"` // hashed - hide sensitive info
}

// Librarian model
type Librarian struct {
	ID       int    `json:"id"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"-"`
}

// Book model
type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title" binding:"required"`
	Author    string    `json:"author" binding:"required"`
	Category  string    `json:"category"`
	Quantity  int       `json:"quantity" binding:"gte=1"` //total copies of this book the library has
	Available int       `json:"available"`
	AddedAt   time.Time `json:"added_at"`
}

// BorrowedBook model
type BorrowedBook struct {
	ID         int        `json:"id"`
	StudentID  int        `json:"student_id"`
	BookID     int        `json:"book_id"`
	BorrowDate time.Time  `json:"borrow_date"`
	ReturnDate *time.Time `json:"return_date,omitempty"`
}

// -------- DATABASE CONNECTION --------
func ConnectDB() *sql.DB {
	dsn := "root:helloworld123@tcp(localhost:3306)/library_db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Connected to database")
	return db
}

// -------- PASSWORD HASHING --------
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// -------- REGISTER STUDENT --------
func RegisterStudent(c *gin.Context, db *sql.DB) {
	var s Student
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashed, err := HashPassword(s.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	_, err = db.Exec("INSERT INTO students (name, email, phone, password) VALUES (?, ?, ?, ?)",
		s.Name, s.Email, s.Phone, hashed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register student"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Student registered successfully"})
}

// -------- CHECK PASSWORD HASH --------
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// -------- LOGIN STUDENT --------
func LoginStudent(c *gin.Context, db *sql.DB) {
	var input Student
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var s Student

	err := db.QueryRow("SELECT id, name, email, phone, password FROM student WHERE email=?", input.Email).
		Scan(&s.ID, &s.Name, &s.Email, &s.Phone, &s.Password)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	if !CheckPasswordHash(input.Password, s.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "student": s})

}

// -------- MAIN FUNCTION --------
func main() {
	db := ConnectDB()
	defer db.Close()

	router := gin.Default()
	router.POST("/students/register", func(c *gin.Context) { RegisterStudent(c, db) })
	router.POST("/students/login", func(c *gin.Context) { LoginStudent(c, db) })
	router.Run(":8080")

}
