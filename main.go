// Package main serves as the entry point for the Library Management System.
// This is a RESTful API built with Go and Gin framework that provides
// functionality for managing library operations including book borrowing,
// student management, and librarian administrative tasks.
package main

import "library-management-system/boot"

// main is the application entry point.
// It initializes and starts the HTTP server with all configured routes,
// database connections, and middleware.
func main() {
	// Delegate server startup to the boot package for better organization
	boot.StartServer()
}
