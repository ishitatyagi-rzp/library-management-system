package service

import (
	"errors"
	"time"

	"library-management-system/constants"
	"library-management-system/dao"
	"library-management-system/model"

	"gorm.io/gorm"
)

// BookService interface defines business logic operations for books
type BookService interface {
	// Book management
	GetAllBooks(limit, offset int) ([]model.Book, int, error)
	GetBookByID(id uint) (*model.Book, error)
	AddBook(book *model.Book) error
	UpdateBook(id uint, updateData map[string]interface{}) (*model.Book, error)
	DeleteBook(id uint) error

	// Borrowing operations
	BorrowBook(studentID, bookID uint) (*BorrowResponse, error)
	ReturnBook(studentID, bookID uint) (*ReturnResponse, error)
}

// BorrowResponse represents the response for a successful borrow operation
type BorrowResponse struct {
	StudentID  uint      `json:"student_id"`
	BookID     uint      `json:"book_id"`
	BookTitle  string    `json:"book_title"`
	BookAuthor string    `json:"book_author"`
	BorrowDate time.Time `json:"borrow_date"`
}

// ReturnResponse represents the response for a successful return operation
type ReturnResponse struct {
	StudentID  uint      `json:"student_id"`
	BookID     uint      `json:"book_id"`
	BookTitle  string    `json:"book_title"`
	BookAuthor string    `json:"book_author"`
	ReturnDate time.Time `json:"return_date"`
}

// bookServiceImpl implements BookService interface
type bookServiceImpl struct {
	db              *gorm.DB
	bookDAO         dao.BookDAO
	studentDAO      dao.StudentDAO
	borrowedBookDAO dao.BorrowedBookDAO
}

// NewBookService creates a new BookService instance
func NewBookService(db *gorm.DB, bookDAO dao.BookDAO, studentDAO dao.StudentDAO, borrowedBookDAO dao.BorrowedBookDAO) BookService {
	return &bookServiceImpl{
		db:              db,
		bookDAO:         bookDAO,
		studentDAO:      studentDAO,
		borrowedBookDAO: borrowedBookDAO,
	}
}

// GetAllBooks retrieves paginated list of books
func (s *bookServiceImpl) GetAllBooks(limit, offset int) ([]model.Book, int, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	books, err := s.bookDAO.GetAll(limit, offset)
	if err != nil {
		return nil, 0, errors.New(constants.FetchBooksError)
	}

	return books, len(books), nil
}

// GetBookByID retrieves a single book by ID
func (s *bookServiceImpl) GetBookByID(id uint) (*model.Book, error) {
	if id <= 0 {
		return nil, errors.New("invalid book id")
	}

	book, err := s.bookDAO.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(constants.BookNotFoundError)
		}
		return nil, errors.New(constants.FetchBookError)
	}

	return book, nil
}

// AddBook creates a new book
func (s *bookServiceImpl) AddBook(book *model.Book) error {
	// Business logic validation
	if book.Quantity <= 0 {
		return errors.New("quantity must be > 0")
	}

	// Initialize available count and timestamp
	book.Available = book.Quantity
	if book.AddedAt.IsZero() {
		book.AddedAt = time.Now()
	}

	if err := s.bookDAO.Create(book); err != nil {
		return errors.New(constants.AddBookError)
	}

	return nil
}

// UpdateBook updates an existing book
func (s *bookServiceImpl) UpdateBook(id uint, updateData map[string]interface{}) (*model.Book, error) {
	if id <= 0 {
		return nil, errors.New("invalid book id")
	}

	// Business logic validation for numeric fields
	if available, exists := updateData["available"]; exists {
		if availableInt, ok := available.(int); ok && availableInt < 0 {
			return nil, errors.New("available copies cannot be negative")
		}
	}

	if quantity, exists := updateData["quantity"]; exists {
		if quantityInt, ok := quantity.(int); ok && quantityInt < 0 {
			return nil, errors.New("quantity cannot be negative")
		}
	}

	// Cross-field validation
	if available, availableExists := updateData["available"]; availableExists {
		if quantity, quantityExists := updateData["quantity"]; quantityExists {
			if availableInt, ok1 := available.(int); ok1 {
				if quantityInt, ok2 := quantity.(int); ok2 && availableInt > quantityInt {
					return nil, errors.New("available copies cannot exceed total quantity")
				}
			}
		}
	}

	if len(updateData) == 0 {
		return nil, errors.New("no valid fields to update")
	}

	book, err := s.bookDAO.Update(id, updateData)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(constants.BookNotFoundError)
		}
		return nil, errors.New(constants.UpdateBookError)
	}

	return book, nil
}

// DeleteBook removes a book
func (s *bookServiceImpl) DeleteBook(id uint) error {
	if id <= 0 {
		return errors.New("invalid book id")
	}

	// Check if book is currently borrowed
	borrowed, err := s.bookDAO.CheckIfCurrentlyBorrowed(id)
	if err != nil {
		return errors.New(constants.DeleteBookError)
	}

	if borrowed {
		return errors.New(constants.BookCurrentlyBorrowedError)
	}

	if err := s.bookDAO.Delete(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(constants.BookNotFoundError)
		}
		return errors.New(constants.DeleteBookError)
	}

	return nil
}

// BorrowBook handles the complete book borrowing process
func (s *bookServiceImpl) BorrowBook(studentID, bookID uint) (*BorrowResponse, error) {
	// Validate input
	if studentID <= 0 || bookID <= 0 {
		return nil, errors.New("student_id and book_id must be positive")
	}

	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, errors.New(constants.BorrowBookError)
	}
	defer tx.Rollback() // Safe to call even after commit

	// Check if student exists
	studentExists, err := s.studentDAO.CheckExists(studentID)
	if err != nil || !studentExists {
		return nil, errors.New(constants.StudentNotFoundError)
	}

	// Lock and get book details
	book, err := s.bookDAO.GetBookForUpdate(tx, bookID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(constants.BookNotFoundError)
		}
		return nil, errors.New(constants.BorrowBookError)
	}

	// Check availability
	if book.Available <= 0 {
		return nil, errors.New(constants.BookNotAvailableError)
	}

	// Check if already borrowed
	alreadyBorrowed, err := s.borrowedBookDAO.HasActiveBorrow(studentID, bookID)
	if err != nil {
		return nil, errors.New(constants.BorrowBookError)
	}
	if alreadyBorrowed {
		return nil, errors.New(constants.BookAlreadyBorrowedError)
	}

	// Decrement available count atomically
	if err := s.bookDAO.DecrementAvailable(tx, bookID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(constants.BookNotAvailableError)
		}
		return nil, errors.New(constants.BorrowBookError)
	}

	// Create borrow record
	now := time.Now()
	borrowRecord := &model.BorrowedBook{
		StudentID:  studentID,
		BookID:     bookID,
		BorrowDate: now,
	}

	if err := s.borrowedBookDAO.Create(tx, borrowRecord); err != nil {
		return nil, errors.New(constants.BorrowBookError)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New(constants.BorrowBookError)
	}

	// Return success response
	return &BorrowResponse{
		StudentID:  studentID,
		BookID:     bookID,
		BookTitle:  book.Title,
		BookAuthor: book.Author,
		BorrowDate: now,
	}, nil
}

// ReturnBook handles the complete book returning process
func (s *bookServiceImpl) ReturnBook(studentID, bookID uint) (*ReturnResponse, error) {
	// Validate input
	if studentID <= 0 || bookID <= 0 {
		return nil, errors.New("student_id and book_id must be positive")
	}

	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, errors.New(constants.ReturnBookError)
	}
	defer tx.Rollback()

	// Find active borrow record
	borrowRecord, err := s.borrowedBookDAO.FindActiveBorrow(tx, studentID, bookID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(constants.BookNotBorrowedError)
		}
		return nil, errors.New(constants.ReturnBookError)
	}

	// Set return date
	now := time.Now()
	if err := s.borrowedBookDAO.UpdateReturnDate(tx, borrowRecord.ID, now); err != nil {
		return nil, errors.New(constants.ReturnBookError)
	}

	// Increment available count
	if err := s.bookDAO.IncrementAvailable(tx, bookID); err != nil {
		return nil, errors.New(constants.ReturnBookError)
	}

	// Get book details for response
	book, err := s.bookDAO.GetByID(bookID)
	if err != nil {
		return nil, errors.New(constants.ReturnBookError)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New(constants.ReturnBookError)
	}

	// Return success response
	return &ReturnResponse{
		StudentID:  studentID,
		BookID:     bookID,
		BookTitle:  book.Title,
		BookAuthor: book.Author,
		ReturnDate: now,
	}, nil
}
