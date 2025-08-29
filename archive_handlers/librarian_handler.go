package handlers

import (
	"net/http"
	"strconv"
	"time"

	"library-management-system/constants"
	"library-management-system/model"
	"library-management-system/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterLibrarian handles librarian registration using GORM
func RegisterLibrarian(c *gin.Context, db *gorm.DB) {
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

	// Set the hashed password
	l.SetPassword(hashed)

	result := db.Create(&l)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.RegisterLibrarianError})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": constants.LibrarianRegisteredMessage})
}

// LoginLibrarian handles librarian login using GORM
func LoginLibrarian(c *gin.Context, db *gorm.DB) {
	var input model.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var l model.Librarian
	result := db.Where("email = ?", input.GetEmail()).First(&l)
	if result.Error != nil {
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

// GetAllStudents retrieves all students using GORM (Librarian only)
func GetAllStudents(c *gin.Context, db *gorm.DB) {
	var students []model.Student
	result := db.Select("id, name, email, phone").Order("name").Find(&students)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve students"})
		return
	}

	// Convert to response format without password
	var responseStudents []gin.H
	for _, student := range students {
		responseStudents = append(responseStudents, gin.H{
			"id":    student.GetID(),
			"name":  student.GetName(),
			"email": student.GetEmail(),
			"phone": student.GetPhone(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Students retrieved successfully",
		"students":       responseStudents,
		"total_students": len(students),
	})
}

// GetStudentDetails retrieves detailed information about a student including borrow history using GORM (Librarian only)
func GetStudentDetails(c *gin.Context, db *gorm.DB) {
	studentIDStr := c.Param("id")
	studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	// Get student basic info using GORM
	var student model.Student
	result := db.Select("id, name, email, phone").First(&student, uint(studentID))
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": constants.StudentNotFoundError})
		return
	}

	// Get borrow history using GORM with Preload for Book details
	var borrowedBooks []model.BorrowedBook
	result = db.Preload("Book").Where("student_id = ?", uint(studentID)).Order("borrow_date DESC").Find(&borrowedBooks)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve borrow history"})
		return
	}

	type BorrowHistoryItem struct {
		ID         uint       `json:"id"`
		BookID     uint       `json:"book_id"`
		Title      string     `json:"title"`
		Author     string     `json:"author"`
		Category   string     `json:"category"`
		BorrowDate time.Time  `json:"borrow_date"`
		ReturnDate *time.Time `json:"return_date,omitempty"`
		Status     string     `json:"status"`
	}

	var history []BorrowHistoryItem
	var currentlyBorrowedCount int

	for _, borrowed := range borrowedBooks {
		item := BorrowHistoryItem{
			ID:         borrowed.GetID(),
			BookID:     borrowed.GetBookID(),
			Title:      borrowed.Book.GetTitle(),
			Author:     borrowed.Book.GetAuthor(),
			Category:   borrowed.Book.GetCategory(),
			BorrowDate: borrowed.GetBorrowDate(),
			ReturnDate: borrowed.GetReturnDate(),
		}

		// Set status and count currently borrowed books
		if item.ReturnDate != nil {
			item.Status = "Returned"
		} else {
			item.Status = "Currently Borrowed"
			currentlyBorrowedCount++
		}

		history = append(history, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Student details retrieved successfully",
		"student": gin.H{
			"id":    student.GetID(),
			"name":  student.GetName(),
			"email": student.GetEmail(),
			"phone": student.GetPhone(),
		},
		"borrow_history":           history,
		"total_books_borrowed":     len(history),
		"currently_borrowed_count": currentlyBorrowedCount,
	})
}

// UpdateStudent allows librarians to update student information using GORM
func UpdateStudent(c *gin.Context, db *gorm.DB) {
	studentIDStr := c.Param("id")
	studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	// Struct to store updated student data
	var updateData struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
		Phone string `json:"phone" binding:"required,len=10,numeric"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if student exists and update using GORM
	var student model.Student
	result := db.First(&student, uint(studentID))
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": constants.StudentNotFoundError})
		return
	}

	// Update student information using GORM
	result = db.Model(&student).Updates(model.Student{
		Name:  updateData.Name,
		Email: updateData.Email,
		Phone: updateData.Phone,
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Student updated successfully",
		"student": gin.H{
			"id":    studentID,
			"name":  updateData.Name,
			"email": updateData.Email,
			"phone": updateData.Phone,
		},
	})
}

// DeleteStudent allows librarians to delete student accounts using GORM (only if no active borrows)
func DeleteStudent(c *gin.Context, db *gorm.DB) {
	studentIDStr := c.Param("id")
	studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	// Check if student has any unreturned books using GORM
	var unreturnedCount int64
	result := db.Model(&model.BorrowedBook{}).Where("student_id = ? AND return_date IS NULL", uint(studentID)).Count(&unreturnedCount)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check student books"})
		return
	}

	if unreturnedCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete student with unreturned books"})
		return
	}

	// Delete student using GORM
	result = db.Delete(&model.Student{}, uint(studentID))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete student"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": constants.StudentNotFoundError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Student deleted successfully",
	})
}
