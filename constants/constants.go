// Package constants provides centralized configuration and message constants
// for the Library Management System. This package ensures consistency across
// the application and makes it easy to modify messages and configuration
// values without searching through multiple files.
package constants

const (
	// Database Configuration Constants
	// These values should be moved to environment variables in production

	// DatabaseDSN is the MySQL database connection string
	// Format: username:password@protocol(address)/dbname
	DatabaseDSN = "root:helloworld123@tcp(localhost:3306)/library_db"

	// Server Configuration Constants

	// ServerPort defines the port on which the HTTP server will listen
	// The colon prefix is required by Gin framework
	ServerPort = ":8080"

	// Success Response Messages
	// These messages are returned when operations complete successfully

	StudentRegisteredMessage   = "Student registered successfully"
	LibrarianRegisteredMessage = "Librarian registered successfully"
	LoginSuccessfulMessage     = "Login successful"

	// Authentication and User Error Messages
	// These messages handle user-related errors and authentication failures

	HashPasswordError       = "Failed to hash password"
	RegisterStudentError    = "Failed to register student"
	RegisterLibrarianError  = "Failed to register librarian"
	InvalidCredentialsError = "Invalid credentials"
	StudentNotFoundError    = "Student not found"

	// Database Connection Error Messages
	// These messages handle database connectivity issues

	DatabaseConnectionError = "Failed to connect to database"
	DatabaseOpenError       = "Failed to open database"

	// Book Operation Error Messages
	// These messages handle book-related business rule violations and errors

	BookNotFoundError          = "Book not found"
	BookNotAvailableError      = "Book is not available for borrowing"
	BorrowBookError            = "Failed to borrow book"
	BookAlreadyBorrowedError   = "Student has already borrowed this book"
	AddBookError               = "Failed to add book"
	UpdateBookError            = "Failed to update book"
	DeleteBookError            = "Failed to delete book"
	BookCurrentlyBorrowedError = "Cannot delete book that is currently borrowed"
	ReturnBookError            = "Failed to return book"
	BookNotBorrowedError       = "Book was not borrowed by this student"

	// Book Operation Success Messages
	// These messages confirm successful book-related operations

	BookBorrowedMessage = "Book borrowed successfully"
	BookAddedMessage    = "Book added successfully"
	BookUpdatedMessage  = "Book updated successfully"
	BookDeletedMessage  = "Book deleted successfully"
	BookReturnedMessage = "Book returned successfully"

	// Additional error constants used in handlers
	FetchBooksError = "Failed to retrieve books"
	FetchBookError  = "Failed to retrieve book"
)
