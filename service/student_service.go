package service

import (
	"errors"
	"time"

	"library-management-system/constants"
	"library-management-system/dao"
	"library-management-system/model"
	"library-management-system/utils"

	"gorm.io/gorm"
)

// StudentService interface defines business logic operations for students
type StudentService interface {
	// Authentication
	RegisterStudent(student *model.Student) error
	LoginStudent(email, password string) (*StudentLoginResponse, error)

	// Student management (for librarians)
	GetAllStudents() ([]StudentInfo, error)
	GetStudentDetails(studentID uint) (*StudentDetails, error)
	UpdateStudent(studentID uint, updateData StudentUpdateRequest) (*StudentInfo, error)
	DeleteStudent(studentID uint) error

	// Student self-service
	GetStudentBorrowHistory(studentID uint) ([]BorrowHistoryItem, error)
}

// StudentLoginResponse represents the response for successful student login
type StudentLoginResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// StudentInfo represents basic student information (without sensitive data)
type StudentInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// StudentDetails represents detailed student information with borrow history
type StudentDetails struct {
	Student                StudentInfo         `json:"student"`
	BorrowHistory          []BorrowHistoryItem `json:"borrow_history"`
	TotalBooksBorrowed     int                 `json:"total_books_borrowed"`
	CurrentlyBorrowedCount int                 `json:"currently_borrowed_count"`
}

// BorrowHistoryItem represents a single borrowing record with book details
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

// StudentUpdateRequest represents data for updating student information
type StudentUpdateRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone" binding:"required,len=10,numeric"`
}

// studentServiceImpl implements StudentService interface
type studentServiceImpl struct {
	studentDAO      dao.StudentDAO
	borrowedBookDAO dao.BorrowedBookDAO
}

// NewStudentService creates a new StudentService instance
func NewStudentService(studentDAO dao.StudentDAO, borrowedBookDAO dao.BorrowedBookDAO) StudentService {
	return &studentServiceImpl{
		studentDAO:      studentDAO,
		borrowedBookDAO: borrowedBookDAO,
	}
}

// RegisterStudent handles student registration with password hashing
func (s *studentServiceImpl) RegisterStudent(student *model.Student) error {
	// Hash the password
	hashedPassword, err := utils.HashPassword(student.GetPassword())
	if err != nil {
		return errors.New(constants.HashPasswordError)
	}

	student.SetPassword(hashedPassword)

	if err := s.studentDAO.Create(student); err != nil {
		return errors.New(constants.RegisterStudentError)
	}

	return nil
}

// LoginStudent handles student authentication
func (s *studentServiceImpl) LoginStudent(email, password string) (*StudentLoginResponse, error) {
	student, err := s.studentDAO.GetByEmail(email)
	if err != nil {
		// Don't reveal whether email exists
		return nil, errors.New(constants.InvalidCredentialsError)
	}

	// Verify password
	if !utils.CheckPasswordHash(password, student.GetPassword()) {
		return nil, errors.New(constants.InvalidCredentialsError)
	}

	return &StudentLoginResponse{
		ID:    student.GetID(),
		Name:  student.GetName(),
		Email: student.GetEmail(),
		Phone: student.GetPhone(),
	}, nil
}

// GetAllStudents retrieves all students (for librarians)
func (s *studentServiceImpl) GetAllStudents() ([]StudentInfo, error) {
	students, err := s.studentDAO.GetAll()
	if err != nil {
		return nil, errors.New("Failed to retrieve students")
	}

	var studentInfos []StudentInfo
	for _, student := range students {
		studentInfos = append(studentInfos, StudentInfo{
			ID:    student.GetID(),
			Name:  student.GetName(),
			Email: student.GetEmail(),
			Phone: student.GetPhone(),
		})
	}

	return studentInfos, nil
}

// GetStudentDetails retrieves detailed student information with borrow history
func (s *studentServiceImpl) GetStudentDetails(studentID uint) (*StudentDetails, error) {
	if studentID <= 0 {
		return nil, errors.New("Invalid student ID")
	}

	// Get student basic info
	student, err := s.studentDAO.GetByID(studentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(constants.StudentNotFoundError)
		}
		return nil, errors.New("Failed to retrieve student details")
	}

	// Get borrow history
	borrowedBooks, err := s.borrowedBookDAO.GetBorrowHistory(studentID)
	if err != nil {
		return nil, errors.New("Failed to retrieve borrow history")
	}

	// Convert to response format
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

	return &StudentDetails{
		Student: StudentInfo{
			ID:    student.GetID(),
			Name:  student.GetName(),
			Email: student.GetEmail(),
			Phone: student.GetPhone(),
		},
		BorrowHistory:          history,
		TotalBooksBorrowed:     len(history),
		CurrentlyBorrowedCount: currentlyBorrowedCount,
	}, nil
}

// UpdateStudent updates student information (for librarians)
func (s *studentServiceImpl) UpdateStudent(studentID uint, updateData StudentUpdateRequest) (*StudentInfo, error) {
	if studentID <= 0 {
		return nil, errors.New("Invalid student ID")
	}

	// Convert to map for DAO
	updateMap := map[string]interface{}{
		"name":  updateData.Name,
		"email": updateData.Email,
		"phone": updateData.Phone,
	}

	student, err := s.studentDAO.Update(studentID, updateMap)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(constants.StudentNotFoundError)
		}
		return nil, errors.New("Failed to update student")
	}

	return &StudentInfo{
		ID:    student.GetID(),
		Name:  student.GetName(),
		Email: student.GetEmail(),
		Phone: student.GetPhone(),
	}, nil
}

// DeleteStudent removes a student account (only if no active borrows)
func (s *studentServiceImpl) DeleteStudent(studentID uint) error {
	if studentID <= 0 {
		return errors.New("Invalid student ID")
	}

	// Check if student has unreturned books
	hasUnreturned, err := s.studentDAO.HasUnreturnedBooks(studentID)
	if err != nil {
		return errors.New("Failed to check student books")
	}

	if hasUnreturned {
		return errors.New("Cannot delete student with unreturned books")
	}

	if err := s.studentDAO.Delete(studentID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(constants.StudentNotFoundError)
		}
		return errors.New("Failed to delete student")
	}

	return nil
}

// GetStudentBorrowHistory retrieves complete borrowing history for a student
func (s *studentServiceImpl) GetStudentBorrowHistory(studentID uint) ([]BorrowHistoryItem, error) {
	if studentID <= 0 {
		return nil, errors.New("Invalid student ID")
	}

	// Verify student exists
	studentExists, err := s.studentDAO.CheckExists(studentID)
	if err != nil || !studentExists {
		return nil, errors.New(constants.StudentNotFoundError)
	}

	// Get borrow history
	borrowedBooks, err := s.borrowedBookDAO.GetBorrowHistory(studentID)
	if err != nil {
		return nil, errors.New("Failed to retrieve borrow history")
	}

	// Convert to response format
	var history []BorrowHistoryItem
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

		// Determine status based on return date
		if item.ReturnDate != nil {
			item.Status = "Returned"
		} else {
			item.Status = "Currently Borrowed"
		}

		history = append(history, item)
	}

	return history, nil
}
