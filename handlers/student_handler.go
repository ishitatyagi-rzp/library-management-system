// Package handlers contains HTTP request handlers for the Library Management System.
// This package specifically handles student-related operations including registration,
// authentication, and borrowing history retrieval. All handlers follow the same pattern:
// 1. Parse and validate request data
// 2. Execute business logic and database operations
// 3. Return appropriate HTTP responses
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"library-management-system/constants"
	"library-management-system/model"
	"library-management-system/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterStudent handles HTTP POST requests for student account registration.
// This endpoint creates a new student account with secure password hashing.
//
// Request Body (JSON):
//   - name: Student's full name (required)
//   - email: Valid email address (required, must be unique)
//   - phone: 10-digit phone number (required, numeric only)
//   - password: Plain text password (required, will be hashed)
//
// Response:
//   - 200: Registration successful with success message
//   - 400: Invalid request data or validation failure
//   - 500: Internal server error (password hashing or database failure)
//
// Example:
//
//	POST /students/register
//	{"name": "John Doe", "email": "john@example.com", "phone": "1234567890", "password": "secret123"}
func RegisterStudent(c *gin.Context, db *gorm.DB) {
	// Parse JSON request body into Student model with validation
	var s model.Student
	if err := c.ShouldBindJSON(&s); err != nil {
		// Gin's binding includes validation - return detailed error message
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password before storing in database for security
	// Never store plain text passwords
	hashed, err := utils.HashPassword(s.GetPassword())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.HashPasswordError})
		return
	}

	// Set the hashed password
	s.SetPassword(hashed)

	// Insert new student record using GORM - automatically prevents SQL injection
	result := db.Create(&s)
	if result.Error != nil {
		// Database error could be due to duplicate email or other constraints
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.RegisterStudentError})
		return
	}

	// Registration successful - return success message
	c.JSON(http.StatusOK, gin.H{"message": constants.StudentRegisteredMessage})
}

// LoginStudent handles HTTP POST requests for student authentication.
// This endpoint verifies student credentials and returns user information on success.
//
// Request Body (JSON):
//   - email: Student's registered email address (required)
//   - password: Student's plain text password (required)
//
// Response:
//   - 200: Login successful with student information (password excluded)
//   - 400: Invalid request format
//   - 401: Invalid credentials (wrong email or password)
//   - 500: Internal server error
//
// Example:
//
//	POST /students/login
//	{"email": "john@example.com", "password": "secret123"}
func LoginStudent(c *gin.Context, db *gorm.DB) {
	// Parse login credentials from request body
	var input model.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve student record by email using GORM
	var s model.Student
	result := db.Where("email = ?", input.GetEmail()).First(&s)
	if result.Error != nil {
		// Student not found or database error - return generic invalid credentials message
		// Don't reveal whether email exists for security reasons
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.InvalidCredentialsError})
		return
	}

	// Verify the provided password against the stored hash
	if !utils.CheckPasswordHash(input.GetPassword(), s.GetPassword()) {
		// Password doesn't match - return same generic error message
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.InvalidCredentialsError})
		return
	}

	// Authentication successful - return user info without password
	// Never include password or password hash in API responses
	c.JSON(http.StatusOK, gin.H{
		"message": constants.LoginSuccessfulMessage,
		"student": gin.H{
			"id":    s.GetID(),
			"name":  s.GetName(),
			"email": s.GetEmail(),
			"phone": s.GetPhone(),
		},
	})
}

// GetStudentBorrowHistory handles HTTP GET requests to retrieve a student's complete borrowing history.
// This endpoint returns both current and past borrowing records with book details.
//
// URL Parameter:
//   - id: Student ID (path parameter)
//
// Response:
//   - 200: History retrieved successfully with array of borrow records
//   - 404: Student not found
//   - 500: Internal server error
//
// Each history record includes:
//   - Book information (title, author, category)
//   - Borrow date
//   - Return date (null if still borrowed)
//   - Status ("Currently Borrowed" or "Returned")
//
// Example:
//
//	GET /students/123/history
func GetStudentBorrowHistory(c *gin.Context, db *gorm.DB) {
	// Extract student ID from URL path parameter and convert to uint
	studentIDStr := c.Param("id")
	studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	// Verify that the student exists in the database using GORM
	// This prevents revealing borrowing information for non-existent students
	var student model.Student
	result := db.First(&student, studentID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": constants.StudentNotFoundError})
		return
	}

	// Retrieve all borrowing records for the student with book details using GORM
	// Using Preload to eagerly load the associated Book data
	// Ordered by borrow_date DESC to show most recent borrows first
	var borrowedBooks []model.BorrowedBook
	result = db.Preload("Book").Where("student_id = ?", uint(studentID)).Order("borrow_date DESC").Find(&borrowedBooks)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve borrow history"})
		return
	}

	// BorrowHistoryItem represents a single borrowing record with book details
	// This struct combines data from both borrowed_books and books tables
	type BorrowHistoryItem struct {
		ID         uint       `json:"id"`                    // Borrow record ID
		BookID     uint       `json:"book_id"`               // Book ID
		Title      string     `json:"title"`                 // Book title
		Author     string     `json:"author"`                // Book author
		Category   string     `json:"category"`              // Book category
		BorrowDate time.Time  `json:"borrow_date"`           // When book was borrowed
		ReturnDate *time.Time `json:"return_date,omitempty"` // When book was returned (null if not returned)
		Status     string     `json:"status"`                // Computed status based on return date
	}

	// Convert GORM results to response format
	var history []BorrowHistoryItem
	for _, borrowed := range borrowedBooks {
		item := BorrowHistoryItem{
			ID:         borrowed.GetID(),
			BookID:     borrowed.GetBookID(),
			Title:      borrowed.Book.GetTitle(),
			Author:     borrowed.Book.GetAuthor(),
			Category:   borrowed.Book.GetCategory(),
			BorrowDate: borrowed.GetBorrowDate(),
			ReturnDate: borrowed.GetReturnDate(),
		}

		// Determine status based on whether book has been returned
		// ReturnDate is null for books that haven't been returned
		if item.ReturnDate != nil {
			item.Status = "Returned"
		} else {
			item.Status = "Currently Borrowed"
		}

		history = append(history, item)
	}

	// Return successful response with complete borrowing history
	c.JSON(http.StatusOK, gin.H{
		"message":        "Borrow history retrieved successfully",
		"student_id":     studentID,
		"borrow_history": history,
		"total_records":  len(history), // Convenient count for frontend pagination
	})
}
