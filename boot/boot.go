// Package boot handles server initialization and routing configuration.
// This package is responsible for setting up the HTTP server, establishing
// database connections, configuring routes, and starting the application.
package boot

import (
	"library-management-system/constants"
	"library-management-system/controller"
	"library-management-system/dao"
	"library-management-system/service"
	"library-management-system/utils"

	"github.com/gin-gonic/gin"
)

// StartServer initializes and starts the HTTP server with all configured routes.
// This function performs the following operations:
//  1. Establishes database connection
//  2. Creates DAO, Service, and Controller instances
//  3. Creates Gin router instance
//  4. Registers all API routes with controllers
//  5. Starts the HTTP server on the configured port
//
// The server will run indefinitely until terminated or an error occurs.
// All routes are organized by functionality following Controller-Service-DAO pattern.
func StartServer() {
	// Establish GORM database connection - will terminate app if connection fails
	db := utils.ConnectDB()

	// Get underlying SQL DB for proper cleanup
	sqlDB, err := db.DB()
	if err == nil {
		// Ensure database connection is closed when function exits to prevent leaks
		defer sqlDB.Close()
	}

	// ========== INITIALIZE LAYERS ==========
	// Create DAO instances for data access
	bookDAO := dao.NewBookDAO(db)
	studentDAO := dao.NewStudentDAO(db)
	librarianDAO := dao.NewLibrarianDAO(db)
	borrowedBookDAO := dao.NewBorrowedBookDAO(db)

	// Create Service instances for business logic
	bookService := service.NewBookService(db, bookDAO, studentDAO, borrowedBookDAO)
	studentService := service.NewStudentService(studentDAO, borrowedBookDAO)
	librarianService := service.NewLibrarianService(librarianDAO)

	// Create Controller instances for HTTP handling
	bookController := controller.NewBookController(bookService)
	studentController := controller.NewStudentController(studentService)
	librarianController := controller.NewLibrarianController(librarianService)

	// Initialize Gin router with default middleware (logger and recovery)
	router := gin.Default()

	// ========== STUDENT ROUTES ==========
	// Routes accessible to students for self-service operations

	// Student Authentication Routes
	// These routes handle student registration and login without requiring authentication

	// POST /students/register - Creates a new student account
	router.POST("/students/register", studentController.RegisterStudent)

	// POST /students/login - Authenticates student credentials and returns user info
	router.POST("/students/login", studentController.LoginStudent)

	// Student Data Access Routes
	// These routes allow students to view their own borrowing information

	// GET /students/:id/history - Retrieves complete borrowing history for a student
	// Parameter :id should match the authenticated student's ID in production
	router.GET("/students/:id/history", studentController.GetStudentBorrowHistory)

	// ========== LIBRARIAN ROUTES ==========
	// Routes for librarians to manage the library system and oversee operations

	// Librarian Authentication Routes
	// These routes handle librarian registration and login

	// POST /librarians/register - Creates a new librarian account
	router.POST("/librarians/register", librarianController.RegisterLibrarian)

	// POST /librarians/login - Authenticates librarian credentials
	router.POST("/librarians/login", librarianController.LoginLibrarian)

	// Student Management Routes (Librarian Access Only)
	// These routes allow librarians to view and manage student accounts
	// In production, these should be protected by librarian authentication middleware

	// GET /librarians/students - Retrieves list of all registered students
	router.GET("/librarians/students", studentController.GetAllStudents)

	// GET /librarians/students/:id - Retrieves detailed student information with borrow history
	router.GET("/librarians/students/:id", studentController.GetStudentDetails)

	// PUT /librarians/students/:id - Updates student account information
	router.PUT("/librarians/students/:id", studentController.UpdateStudent)

	// DELETE /librarians/students/:id - Deletes student account (only if no active borrows)
	router.DELETE("/librarians/students/:id", studentController.DeleteStudent)

	// ========== BOOK ROUTES ==========
	// Routes for book-related operations accessible to both students and librarians

	// Public Book Information Routes
	// These routes allow anyone to view book information and availability

	// GET /books - Retrieves list of all books with availability status
	router.GET("/books", bookController.GetAllBooks)

	// GET /books/:id - Retrieves detailed information about a specific book
	router.GET("/books/:id", bookController.GetBookByID)

	// Book Transaction Routes
	// These routes handle borrowing and returning of books
	// In production, should require student authentication

	// POST /books/borrow - Allows students to borrow available books
	// Requires student_id and book_id in request body
	router.POST("/books/borrow", bookController.BorrowBook)

	// POST /books/return - Allows students to return borrowed books
	// Requires student_id and book_id in request body
	router.POST("/books/return", bookController.ReturnBook)

	// Book Management Routes (Librarian Access Only)
	// These routes allow librarians to manage the book inventory
	// In production, should require librarian authentication

	// POST /books - Adds new books to the library collection
	router.POST("/books", bookController.AddBook)

	// PUT /books/:id - Updates existing book information and quantities
	router.PUT("/books/:id", bookController.UpdateBook)

	// DELETE /books/:id - Removes books from collection (only if not currently borrowed)
	router.DELETE("/books/:id", bookController.DeleteBook)

	// Start the HTTP server and listen for incoming requests
	// This is a blocking call - the function will not return until server stops
	router.Run(constants.ServerPort)
}
