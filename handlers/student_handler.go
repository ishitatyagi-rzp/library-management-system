package handlers

import (
	"database/sql"
	"net/http"

	"library-management-system/constants"
	"library-management-system/model"
	"library-management-system/utils"

	"github.com/gin-gonic/gin"
)

// RegisterStudent handles student registration
func RegisterStudent(c *gin.Context, db *sql.DB) {
	var s model.Student
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, err := utils.HashPassword(s.GetPassword())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.HashPasswordError})
		return
	}

	_, err = db.Exec("INSERT INTO students (name, email, phone, password) VALUES (?, ?, ?, ?)",
		s.GetName(), s.GetEmail(), s.GetPhone(), hashed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.RegisterStudentError})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": constants.StudentRegisteredMessage})
}

// LoginStudent handles student login
func LoginStudent(c *gin.Context, db *sql.DB) {
	var input model.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var s model.Student
	err := db.QueryRow("SELECT id, name, email, phone, password FROM students WHERE email=?", input.GetEmail()).
		Scan(&s.ID, &s.Name, &s.Email, &s.Phone, &s.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.InvalidCredentialsError})
		return
	}

	if !utils.CheckPasswordHash(input.GetPassword(), s.GetPassword()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.InvalidCredentialsError})
		return
	}

	// Don't return password in response
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
