package constants

const (
	// Database Configuration
	DatabaseDSN = "root:helloworld123@tcp(localhost:3306)/library_db"

	// Server Configuration
	ServerPort = ":8080"

	// Response Messages
	StudentRegisteredMessage   = "Student registered successfully"
	LibrarianRegisteredMessage = "Librarian registered successfully"
	LoginSuccessfulMessage     = "Login successful"

	// Error Messages
	HashPasswordError       = "Failed to hash password"
	RegisterStudentError    = "Failed to register student"
	RegisterLibrarianError  = "Failed to register librarian"
	InvalidCredentialsError = "Invalid credentials"
	DatabaseConnectionError = "Failed to connect to database"
	DatabaseOpenError       = "Failed to open database"
)
