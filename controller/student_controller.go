package controller

import (
	"net/http"
	"strconv"

	"library-management-system/constants"
	"library-management-system/model"
	"library-management-system/service"

	"github.com/gin-gonic/gin"
)

// StudentController handles HTTP requests for student operations
type StudentController struct {
	studentService service.StudentService
}

// NewStudentController creates a new StudentController instance
func NewStudentController(studentService service.StudentService) *StudentController {
	return &StudentController{
		studentService: studentService,
	}
}

// RegisterStudent handles POST /students/register - student account registration
func (ctrl *StudentController) RegisterStudent(c *gin.Context) {
	var student model.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.studentService.RegisterStudent(&student); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": constants.StudentRegisteredMessage})
}

// LoginStudent handles POST /students/login - student authentication
func (ctrl *StudentController) LoginStudent(c *gin.Context) {
	var input model.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ctrl.studentService.LoginStudent(input.GetEmail(), input.GetPassword())
	if err != nil {
		if err.Error() == constants.InvalidCredentialsError {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": constants.LoginSuccessfulMessage,
		"student": response,
	})
}

// GetStudentBorrowHistory handles GET /students/:id/history - retrieves student's borrow history
func (ctrl *StudentController) GetStudentBorrowHistory(c *gin.Context) {
	studentIDStr := c.Param("id")
	studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	history, err := ctrl.studentService.GetStudentBorrowHistory(uint(studentID))
	if err != nil {
		if err.Error() == constants.StudentNotFoundError {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Borrow history retrieved successfully",
		"student_id":     studentID,
		"borrow_history": history,
		"total_records":  len(history),
	})
}

// GetAllStudents handles GET /librarians/students - retrieves all students (Librarian only)
func (ctrl *StudentController) GetAllStudents(c *gin.Context) {
	students, err := ctrl.studentService.GetAllStudents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Students retrieved successfully",
		"students":       students,
		"total_students": len(students),
	})
}

// GetStudentDetails handles GET /librarians/students/:id - retrieves detailed student info (Librarian only)
func (ctrl *StudentController) GetStudentDetails(c *gin.Context) {
	studentIDStr := c.Param("id")
	studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	details, err := ctrl.studentService.GetStudentDetails(uint(studentID))
	if err != nil {
		if err.Error() == constants.StudentNotFoundError {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":                  "Student details retrieved successfully",
		"student":                  details.Student,
		"borrow_history":           details.BorrowHistory,
		"total_books_borrowed":     details.TotalBooksBorrowed,
		"currently_borrowed_count": details.CurrentlyBorrowedCount,
	})
}

// UpdateStudent handles PUT /librarians/students/:id - updates student information (Librarian only)
func (ctrl *StudentController) UpdateStudent(c *gin.Context) {
	studentIDStr := c.Param("id")
	studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var updateData service.StudentUpdateRequest
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student, err := ctrl.studentService.UpdateStudent(uint(studentID), updateData)
	if err != nil {
		if err.Error() == constants.StudentNotFoundError {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Student updated successfully",
		"student": student,
	})
}

// DeleteStudent handles DELETE /librarians/students/:id - deletes student account (Librarian only)
func (ctrl *StudentController) DeleteStudent(c *gin.Context) {
	studentIDStr := c.Param("id")
	studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	err = ctrl.studentService.DeleteStudent(uint(studentID))
	if err != nil {
		if err.Error() == constants.StudentNotFoundError {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err.Error() == "Cannot delete student with unreturned books" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Student deleted successfully",
	})
}
