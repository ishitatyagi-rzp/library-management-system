package dao

import (
	"library-management-system/model"

	"gorm.io/gorm"
)

// StudentDAO interface defines all database operations for students
type StudentDAO interface {
	// Basic CRUD operations
	GetAll() ([]model.Student, error)
	GetByID(id uint) (*model.Student, error)
	GetByEmail(email string) (*model.Student, error)
	Create(student *model.Student) error
	Update(id uint, updateData map[string]interface{}) (*model.Student, error)
	Delete(id uint) error

	// Student-specific operations
	CheckExists(id uint) (bool, error)
	GetBorrowHistory(studentID uint) ([]model.BorrowedBook, error)
	HasUnreturnedBooks(studentID uint) (bool, error)
}

// studentDAOImpl implements StudentDAO interface
type studentDAOImpl struct {
	db *gorm.DB
}

// NewStudentDAO creates a new StudentDAO instance
func NewStudentDAO(db *gorm.DB) StudentDAO {
	return &studentDAOImpl{db: db}
}

// GetAll retrieves all students (without passwords)
func (dao *studentDAOImpl) GetAll() ([]model.Student, error) {
	var students []model.Student
	result := dao.db.Select("id, name, email, phone").Order("name").Find(&students)
	return students, result.Error
}

// GetByID retrieves a student by ID
func (dao *studentDAOImpl) GetByID(id uint) (*model.Student, error) {
	var student model.Student
	result := dao.db.First(&student, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &student, nil
}

// GetByEmail retrieves a student by email (for login)
func (dao *studentDAOImpl) GetByEmail(email string) (*model.Student, error) {
	var student model.Student
	result := dao.db.Where("email = ?", email).First(&student)
	if result.Error != nil {
		return nil, result.Error
	}
	return &student, nil
}

// Create creates a new student record
func (dao *studentDAOImpl) Create(student *model.Student) error {
	return dao.db.Create(student).Error
}

// Update updates a student with the provided data
func (dao *studentDAOImpl) Update(id uint, updateData map[string]interface{}) (*model.Student, error) {
	var student model.Student

	// First check if student exists
	if err := dao.db.First(&student, id).Error; err != nil {
		return nil, err
	}

	// Update the student
	if err := dao.db.Model(&student).Updates(updateData).Error; err != nil {
		return nil, err
	}

	// Return updated student (without password)
	if err := dao.db.Select("id, name, email, phone").First(&student, id).Error; err != nil {
		return nil, err
	}

	return &student, nil
}

// Delete removes a student from the database
func (dao *studentDAOImpl) Delete(id uint) error {
	result := dao.db.Delete(&model.Student{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// CheckExists checks if a student exists by ID
func (dao *studentDAOImpl) CheckExists(id uint) (bool, error) {
	var count int64
	err := dao.db.Model(&model.Student{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetBorrowHistory retrieves complete borrowing history for a student
func (dao *studentDAOImpl) GetBorrowHistory(studentID uint) ([]model.BorrowedBook, error) {
	var borrowedBooks []model.BorrowedBook
	result := dao.db.Preload("Book").Where("student_id = ?", studentID).Order("borrow_date DESC").Find(&borrowedBooks)
	return borrowedBooks, result.Error
}

// HasUnreturnedBooks checks if a student has any unreturned books
func (dao *studentDAOImpl) HasUnreturnedBooks(studentID uint) (bool, error) {
	var count int64
	err := dao.db.Model(&model.BorrowedBook{}).
		Where("student_id = ? AND return_date IS NULL", studentID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
