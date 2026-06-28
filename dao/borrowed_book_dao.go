package dao

import (
	"library-management-system/model"
	"time"

	"gorm.io/gorm"
)

// BorrowedBookDAO interface defines all database operations for borrowed books
type BorrowedBookDAO interface {
	// Borrowing operations
	Create(tx *gorm.DB, borrowRecord *model.BorrowedBook) error
	FindActiveBorrow(tx *gorm.DB, studentID, bookID uint) (*model.BorrowedBook, error)
	UpdateReturnDate(tx *gorm.DB, borrowRecordID uint, returnDate time.Time) error

	// Query operations
	HasActiveBorrow(studentID, bookID uint) (bool, error)
	GetBorrowHistory(studentID uint) ([]model.BorrowedBook, error)
	CountActiveBorrows(bookID uint) (int64, error)
	CountStudentActiveBorrows(studentID uint) (int64, error)
}

// borrowedBookDAOImpl implements BorrowedBookDAO interface
type borrowedBookDAOImpl struct {
	db *gorm.DB
}

// NewBorrowedBookDAO creates a new BorrowedBookDAO instance
func NewBorrowedBookDAO(db *gorm.DB) BorrowedBookDAO {
	return &borrowedBookDAOImpl{db: db}
}

// Create creates a new borrow record within a transaction
func (dao *borrowedBookDAOImpl) Create(tx *gorm.DB, borrowRecord *model.BorrowedBook) error {
	return tx.Create(borrowRecord).Error
}

// FindActiveBorrow finds an active borrow record for a student and book
func (dao *borrowedBookDAOImpl) FindActiveBorrow(tx *gorm.DB, studentID, bookID uint) (*model.BorrowedBook, error) {
	var borrowRecord model.BorrowedBook
	result := tx.Where("student_id = ? AND book_id = ? AND return_date IS NULL", studentID, bookID).First(&borrowRecord)
	if result.Error != nil {
		return nil, result.Error
	}
	return &borrowRecord, nil
}

// UpdateReturnDate sets the return date for a borrow record
func (dao *borrowedBookDAOImpl) UpdateReturnDate(tx *gorm.DB, borrowRecordID uint, returnDate time.Time) error {
	return tx.Model(&model.BorrowedBook{}).Where("id = ?", borrowRecordID).Update("return_date", returnDate).Error
}

// HasActiveBorrow checks if a student has an active borrow for a specific book
func (dao *borrowedBookDAOImpl) HasActiveBorrow(studentID, bookID uint) (bool, error) {
	var count int64
	err := dao.db.Model(&model.BorrowedBook{}).
		Where("student_id = ? AND book_id = ? AND return_date IS NULL", studentID, bookID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetBorrowHistory retrieves borrow history for a student with book details
func (dao *borrowedBookDAOImpl) GetBorrowHistory(studentID uint) ([]model.BorrowedBook, error) {
	var borrowedBooks []model.BorrowedBook
	result := dao.db.Preload("Book").Where("student_id = ?", studentID).Order("borrow_date DESC").Find(&borrowedBooks)
	return borrowedBooks, result.Error
}

// CountActiveBorrows counts active borrows for a specific book
func (dao *borrowedBookDAOImpl) CountActiveBorrows(bookID uint) (int64, error) {
	var count int64
	err := dao.db.Model(&model.BorrowedBook{}).
		Where("book_id = ? AND return_date IS NULL", bookID).
		Count(&count).Error
	return count, err
}

// CountStudentActiveBorrows counts active borrows for a specific student
func (dao *borrowedBookDAOImpl) CountStudentActiveBorrows(studentID uint) (int64, error) {
	var count int64
	err := dao.db.Model(&model.BorrowedBook{}).
		Where("student_id = ? AND return_date IS NULL", studentID).
		Count(&count).Error
	return count, err
}
