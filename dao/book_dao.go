package dao

import (
	"library-management-system/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// BookDAO interface defines all database operations for books
type BookDAO interface {
	// Basic CRUD operations
	GetAll(limit, offset int) ([]model.Book, error)
	GetByID(id uint) (*model.Book, error)
	Create(book *model.Book) error
	Update(id uint, updateData map[string]interface{}) (*model.Book, error)
	Delete(id uint) error

	// Book-specific operations
	GetBookForUpdate(tx *gorm.DB, bookID uint) (*model.Book, error)
	DecrementAvailable(tx *gorm.DB, bookID uint) error
	IncrementAvailable(tx *gorm.DB, bookID uint) error
	CheckIfCurrentlyBorrowed(bookID uint) (bool, error)
}

// bookDAOImpl implements BookDAO interface
type bookDAOImpl struct {
	db *gorm.DB
}

// NewBookDAO creates a new BookDAO instance
func NewBookDAO(db *gorm.DB) BookDAO {
	return &bookDAOImpl{db: db}
}

// GetAll retrieves books with pagination
func (dao *bookDAOImpl) GetAll(limit, offset int) ([]model.Book, error) {
	var books []model.Book
	result := dao.db.Order("id DESC").
		Limit(limit).
		Offset(offset).
		Find(&books)
	return books, result.Error
}

// GetByID retrieves a book by its ID
func (dao *bookDAOImpl) GetByID(id uint) (*model.Book, error) {
	var book model.Book
	result := dao.db.First(&book, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

// Create creates a new book record
func (dao *bookDAOImpl) Create(book *model.Book) error {
	return dao.db.Create(book).Error
}

// Update updates a book with the provided data
func (dao *bookDAOImpl) Update(id uint, updateData map[string]interface{}) (*model.Book, error) {
	var book model.Book

	// First check if book exists
	if err := dao.db.First(&book, id).Error; err != nil {
		return nil, err
	}

	// Update the book
	if err := dao.db.Model(&book).Updates(updateData).Error; err != nil {
		return nil, err
	}

	// Return updated book
	if err := dao.db.First(&book, id).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

// Delete removes a book from the database
func (dao *bookDAOImpl) Delete(id uint) error {
	result := dao.db.Delete(&model.Book{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetBookForUpdate retrieves a book with row-level locking (FOR UPDATE)
func (dao *bookDAOImpl) GetBookForUpdate(tx *gorm.DB, bookID uint) (*model.Book, error) {
	var book model.Book
	result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&book, bookID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

// DecrementAvailable decreases available count atomically
func (dao *bookDAOImpl) DecrementAvailable(tx *gorm.DB, bookID uint) error {
	result := tx.Model(&model.Book{}).
		Where("id = ? AND available > 0", bookID).
		UpdateColumn("available", gorm.Expr("available - ?", 1))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // No rows affected means book not available
	}
	return nil
}

// IncrementAvailable increases available count atomically
func (dao *bookDAOImpl) IncrementAvailable(tx *gorm.DB, bookID uint) error {
	result := tx.Model(&model.Book{}).
		Where("id = ? AND available < quantity", bookID).
		UpdateColumn("available", gorm.Expr("available + ?", 1))

	return result.Error // It's OK if RowsAffected == 0 (available == quantity)
}

// CheckIfCurrentlyBorrowed checks if a book has any unreturned borrows
func (dao *bookDAOImpl) CheckIfCurrentlyBorrowed(bookID uint) (bool, error) {
	var count int64
	err := dao.db.Model(&model.BorrowedBook{}).
		Where("book_id = ? AND return_date IS NULL", bookID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
