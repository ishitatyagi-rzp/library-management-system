package controller

import (
	"net/http"

	"library-management-system/constants"
	"library-management-system/model"
	"library-management-system/service"

	"github.com/gin-gonic/gin"
)

// LibrarianController handles HTTP requests for librarian operations
type LibrarianController struct {
	librarianService service.LibrarianService
}

// NewLibrarianController creates a new LibrarianController instance
func NewLibrarianController(librarianService service.LibrarianService) *LibrarianController {
	return &LibrarianController{
		librarianService: librarianService,
	}
}

// RegisterLibrarian handles POST /librarians/register - librarian account registration
func (ctrl *LibrarianController) RegisterLibrarian(c *gin.Context) {
	var librarian model.Librarian
	if err := c.ShouldBindJSON(&librarian); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.librarianService.RegisterLibrarian(&librarian); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": constants.LibrarianRegisteredMessage})
}

// LoginLibrarian handles POST /librarians/login - librarian authentication
func (ctrl *LibrarianController) LoginLibrarian(c *gin.Context) {
	var input model.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ctrl.librarianService.LoginLibrarian(input.GetEmail(), input.GetPassword())
	if err != nil {
		if err.Error() == constants.InvalidCredentialsError {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   constants.LoginSuccessfulMessage,
		"librarian": response,
	})
}
