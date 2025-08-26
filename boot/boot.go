package boot

import (
	"library-management-system/constants"
	"library-management-system/handlers"
	"library-management-system/utils"

	"github.com/gin-gonic/gin"
)

// StartServer initializes and starts the server
func StartServer() {
	// Connect to database
	db := utils.ConnectDB()
	defer db.Close()

	// Initialize Gin router
	router := gin.Default()

	// Student routes
	router.POST("/students/register", func(c *gin.Context) {
		handlers.RegisterStudent(c, db)
	})
	router.POST("/students/login", func(c *gin.Context) {
		handlers.LoginStudent(c, db)
	})

	// Librarian routes
	router.POST("/librarians/register", func(c *gin.Context) {
		handlers.RegisterLibrarian(c, db)
	})
	router.POST("/librarians/login", func(c *gin.Context) {
		handlers.LoginLibrarian(c, db)
	})

	// Start the server
	router.Run(constants.ServerPort)
}
