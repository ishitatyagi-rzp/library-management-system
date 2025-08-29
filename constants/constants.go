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

	HashPasswordError       = "failed to hash password"
	RegisterStudentError    = "failed to register student"
	RegisterLibrarianError  = "failed to register librarian"
	InvalidCredentialsError = "invalid credentials"
	StudentNotFoundError    = "student not found"

	// Database Connection Error Messages
	// These messages handle database connectivity issues

	DatabaseConnectionError = "failed to connect to database"
	DatabaseOpenError       = "failed to open database"

	// Book Operation Error Messages
	// These messages handle book-related business rule violations and errors

	BookNotFoundError          = "book not found"
	BookNotAvailableError      = "book is not available for borrowing"
	BorrowBookError            = "failed to borrow book"
	BookAlreadyBorrowedError   = "student has already borrowed this book"
	AddBookError               = "failed to add book"
	UpdateBookError            = "failed to update book"
	DeleteBookError            = "failed to delete book"
	BookCurrentlyBorrowedError = "cannot delete book that is currently borrowed"
	ReturnBookError            = "failed to return book"
	BookNotBorrowedError       = "book was not borrowed by this student"

	// Book Operation Success Messages
	// These messages confirm successful book-related operations

	BookBorrowedMessage = "Book borrowed successfully"
	BookAddedMessage    = "Book added successfully"
	BookUpdatedMessage  = "Book updated successfully"
	BookDeletedMessage  = "Book deleted successfully"
	BookReturnedMessage = "Book returned successfully"

	// Additional error constants used in handlers
	FetchBooksError = "failed to retrieve books"
	FetchBookError  = "failed to retrieve book"
)
