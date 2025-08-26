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
	Password string `json:"password" binding:"required"` // Allow binding for registration
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Librarian model
type Librarian struct {
	ID       int    `json:"id"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"` // Allow binding for registration
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

// -------- REGISTER LIBRARIAN --------

func RegisterLibrarian(c *gin.Context, db *sql.DB) {
	var l Librarian
	if err := c.ShouldBindJSON(&l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashed, err := HashPassword(l.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	_, err = db.Exec("INSERT INTO librarians (name, email, password) VALUES (?, ?, ?)", l.Name, l.Email, hashed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register librarian"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Librarian registered successfully"})
}

// -------- CHECK PASSWORD HASH --------
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// -------- LOGIN STUDENT --------
func LoginStudent(c *gin.Context, db *sql.DB) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var s Student
	err := db.QueryRow("SELECT id, name, email, phone, password FROM students WHERE email=?", input.Email).
		Scan(&s.ID, &s.Name, &s.Email, &s.Phone, &s.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !CheckPasswordHash(input.Password, s.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Don't return password in response
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"student": gin.H{
			"id":    s.ID,
			"name":  s.Name,
			"email": s.Email,
			"phone": s.Phone,
		},
	})
}

// -------- LOGIN LIBRARIAN --------
func LoginLibrarian(c *gin.Context, db *sql.DB) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var l Librarian
	err := db.QueryRow("SELECT id, name, email, password FROM librarians WHERE email=?", input.Email).
		Scan(&l.ID, &l.Name, &l.Email, &l.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !CheckPasswordHash(input.Password, l.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Don't return password in response
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"librarian": gin.H{
			"id":    l.ID,
			"name":  l.Name,
			"email": l.Email,
		},
	})
}

// -------- MAIN FUNCTION --------
func main() {
	db := ConnectDB()
	defer db.Close()

	router := gin.Default()
	router.POST("/students/register", func(c *gin.Context) { RegisterStudent(c, db) })
	router.POST("/students/login", func(c *gin.Context) { LoginStudent(c, db) })
	router.POST("/librarians/register", func(c *gin.Context) { RegisterLibrarian(c, db) })
	router.POST("/librarians/login", func(c *gin.Context) { LoginLibrarian(c, db) })
	router.Run(":8080")

}
