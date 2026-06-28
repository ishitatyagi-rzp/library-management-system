package dao

import (
	"library-management-system/model"

	"gorm.io/gorm"
)

// LibrarianDAO interface defines all database operations for librarians
type LibrarianDAO interface {
	// Basic CRUD operations
	GetByID(id uint) (*model.Librarian, error)
	GetByEmail(email string) (*model.Librarian, error)
	Create(librarian *model.Librarian) error

	// Librarian-specific operations
	CheckExists(id uint) (bool, error)
}

// librarianDAOImpl implements LibrarianDAO interface
type librarianDAOImpl struct {
	db *gorm.DB
}

// NewLibrarianDAO creates a new LibrarianDAO instance
func NewLibrarianDAO(db *gorm.DB) LibrarianDAO {
	return &librarianDAOImpl{db: db}
}

// GetByID retrieves a librarian by ID
func (dao *librarianDAOImpl) GetByID(id uint) (*model.Librarian, error) {
	var librarian model.Librarian
	result := dao.db.First(&librarian, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &librarian, nil
}

// GetByEmail retrieves a librarian by email (for login)
func (dao *librarianDAOImpl) GetByEmail(email string) (*model.Librarian, error) {
	var librarian model.Librarian
	result := dao.db.Where("email = ?", email).First(&librarian)
	if result.Error != nil {
		return nil, result.Error
	}
	return &librarian, nil
}

// Create creates a new librarian record
func (dao *librarianDAOImpl) Create(librarian *model.Librarian) error {
	return dao.db.Create(librarian).Error
}

// CheckExists checks if a librarian exists by ID
func (dao *librarianDAOImpl) CheckExists(id uint) (bool, error) {
	var count int64
	err := dao.db.Model(&model.Librarian{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
