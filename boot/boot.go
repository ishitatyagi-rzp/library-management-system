// Package boot handles server initialization and routing configuration.
// This package is responsible for setting up the HTTP server, establishing
// database connections, configuring routes, and starting the application.
// It acts as the central orchestrator that brings together all other packages.
package boot

import (
	"library-management-system/constants"
	"library-management-system/handlers"
	"library-management-system/utils"

	"github.com/gin-gonic/gin"
)

// StartServer initializes and starts the HTTP server with all configured routes.
// This function performs the following operations:
//  1. Establishes database connection
//  2. Creates Gin router instance
//  3. Registers all API routes
//  4. Starts the HTTP server on the configured port
//
// The server will run indefinitely until terminated or an error occurs.
// All routes are organized by functionality (student, librarian, book operations).
func StartServer() {
	// Establish GORM database connection - will terminate app if connection fails
	db := utils.ConnectDB()

	// Get underlying SQL DB for proper cleanup
	sqlDB, err := db.DB()
	if err == nil {
		// Ensure database connection is closed when function exits to prevent leaks
		defer sqlDB.Close()
	}

	// Initialize Gin router with default middleware (logger and recovery)
	router := gin.Default()

	// ========== STUDENT ROUTES ==========
	// Routes accessible to students for self-service operations

	// Student Authentication Routes
	// These routes handle student registration and login without requiring authentication

	// POST /students/register - Creates a new student account
	router.POST("/students/register", func(c *gin.Context) {
		handlers.RegisterStudent(c, db)
	})

	// POST /students/login - Authenticates student credentials and returns user info
	router.POST("/students/login", func(c *gin.Context) {
		handlers.LoginStudent(c, db)
	})

	// Student Data Access Routes
	// These routes allow students to view their own borrowing information

	// GET /students/:id/history - Retrieves complete borrowing history for a student
	// Parameter :id should match the authenticated student's ID in production
	router.GET("/students/:id/history", func(c *gin.Context) {
		handlers.GetStudentBorrowHistory(c, db)
	})

	// ========== LIBRARIAN ROUTES ==========
	// Routes for librarians to manage the library system and oversee operations

	// Librarian Authentication Routes
	// These routes handle librarian registration and login

	// POST /librarians/register - Creates a new librarian account
	router.POST("/librarians/register", func(c *gin.Context) {
		handlers.RegisterLibrarian(c, db)
	})

	// POST /librarians/login - Authenticates librarian credentials
	router.POST("/librarians/login", func(c *gin.Context) {
		handlers.LoginLibrarian(c, db)
	})

	// Student Management Routes (Librarian Access Only)
	// These routes allow librarians to view and manage student accounts
	// In production, these should be protected by librarian authentication middleware

	// GET /librarians/students - Retrieves list of all registered students
	router.GET("/librarians/students", func(c *gin.Context) {
		handlers.GetAllStudents(c, db)
	})

	// GET /librarians/students/:id - Retrieves detailed student information with borrow history
	router.GET("/librarians/students/:id", func(c *gin.Context) {
		handlers.GetStudentDetails(c, db)
	})

	// PUT /librarians/students/:id - Updates student account information
	router.PUT("/librarians/students/:id", func(c *gin.Context) {
		handlers.UpdateStudent(c, db)
	})

	// DELETE /librarians/students/:id - Deletes student account (only if no active borrows)
	router.DELETE("/librarians/students/:id", func(c *gin.Context) {
		handlers.DeleteStudent(c, db)
	})

	// ========== BOOK ROUTES ==========
	// Routes for book-related operations accessible to both students and librarians

	// Public Book Information Routes
	// These routes allow anyone to view book information and availability

	// GET /books - Retrieves list of all books with availability status
	router.GET("/books", func(c *gin.Context) {
		handlers.GetAllBooks(c, db)
	})

	// GET /books/:id - Retrieves detailed information about a specific book
	router.GET("/books/:id", func(c *gin.Context) {
		handlers.GetBookByID(c, db)
	})

	// Book Transaction Routes
	// These routes handle borrowing and returning of books
	// In production, should require student authentication

	// POST /books/borrow - Allows students to borrow available books
	// Requires student_id and book_id in request body
	router.POST("/books/borrow", func(c *gin.Context) {
		handlers.BorrowBook(c, db)
	})

	// POST /books/return - Allows students to return borrowed books
	// Requires student_id and book_id in request body
	router.POST("/books/return", func(c *gin.Context) {
		handlers.ReturnBook(c, db)
	})

	// Book Management Routes (Librarian Access Only)
	// These routes allow librarians to manage the book inventory
	// In production, should require librarian authentication

	// POST /books - Adds new books to the library collection
	router.POST("/books", func(c *gin.Context) {
		handlers.AddBook(c, db)
	})

	// PUT /books/:id - Updates existing book information and quantities
	router.PUT("/books/:id", func(c *gin.Context) {
		handlers.UpdateBook(c, db)
	})

	// DELETE /books/:id - Removes books from collection (only if not currently borrowed)
	router.DELETE("/books/:id", func(c *gin.Context) {
		handlers.DeleteBook(c, db)
	})

	// Start the HTTP server and listen for incoming requests
	// This is a blocking call - the function will not return until server stops
	router.Run(constants.ServerPort)
}
