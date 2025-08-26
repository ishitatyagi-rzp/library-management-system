package handlers

import (
	"database/sql"
	"net/http"

	"library-management-system/constants"
	"library-management-system/model"
	"library-management-system/utils"

	"github.com/gin-gonic/gin"
)

// RegisterLibrarian handles librarian registration
func RegisterLibrarian(c *gin.Context, db *sql.DB) {
	var l model.Librarian
	if err := c.ShouldBindJSON(&l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, err := utils.HashPassword(l.GetPassword())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.HashPasswordError})
		return
	}

	_, err = db.Exec("INSERT INTO librarians (name, email, password) VALUES (?, ?, ?)",
		l.GetName(), l.GetEmail(), hashed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.RegisterLibrarianError})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": constants.LibrarianRegisteredMessage})
}

// LoginLibrarian handles librarian login
func LoginLibrarian(c *gin.Context, db *sql.DB) {
	var input model.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var l model.Librarian
	err := db.QueryRow("SELECT id, name, email, password FROM librarians WHERE email=?", input.GetEmail()).
		Scan(&l.ID, &l.Name, &l.Email, &l.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.InvalidCredentialsError})
		return
	}

	if !utils.CheckPasswordHash(input.GetPassword(), l.GetPassword()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.InvalidCredentialsError})
		return
	}

	// Don't return password in response
	c.JSON(http.StatusOK, gin.H{
		"message": constants.LoginSuccessfulMessage,
		"librarian": gin.H{
			"id":    l.GetID(),
			"name":  l.GetName(),
			"email": l.GetEmail(),
		},
	})
}
