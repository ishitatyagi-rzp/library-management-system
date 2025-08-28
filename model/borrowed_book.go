package model

import "time"

// BorrowedBook represents a transaction record when a student borrows a book.
// This struct maps to the 'borrowed_books' table and tracks the complete
// lifecycle of a book loan from borrowing to returning.
//
// Database Table: borrowed_books
// Relationships:
//   - student_id: Foreign key to students.id
//   - book_id: Foreign key to books.id
//
// Business Logic:
//   - ReturnDate is NULL when book is currently borrowed
//   - ReturnDate is set to current timestamp when book is returned
//   - Each active borrow (ReturnDate = NULL) decreases book.available by 1
//   - Each return sets ReturnDate and increases book.available by 1
//
// GORM Tags: Used for database schema generation, foreign keys, and relationships
type BorrowedBook struct {
	ID         uint       `json:"id" gorm:"primaryKey;autoIncrement"` // Primary key, auto-generated
	StudentID  uint       `json:"student_id" gorm:"not null;index"`   // Foreign key to students table
	BookID     uint       `json:"book_id" gorm:"not null;index"`      // Foreign key to books table
	BorrowDate time.Time  `json:"borrow_date" gorm:"autoCreateTime"`  // Automatically set when record is created
	ReturnDate *time.Time `json:"return_date,omitempty"`              // When book was returned (NULL if not returned)

	// GORM Relationships - enables automatic joins and eager loading
	Student Student `json:"-" gorm:"foreignKey:StudentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Belongs to Student
	Book    Book    `json:"-" gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`    // Belongs to Book
}

// NewBorrowedBook creates a new BorrowedBook instance with GORM-compatible fields.
// GORM will automatically set BorrowDate timestamp.
func NewBorrowedBook(studentID, bookID uint) *BorrowedBook {
	return &BorrowedBook{
		StudentID: studentID,
		BookID:    bookID,
		// BorrowDate will be set automatically by GORM
	}
}

// GetID returns the borrowed book record's ID
func (bb *BorrowedBook) GetID() uint {
	return bb.ID
}

// SetID sets the borrowed book record's ID
func (bb *BorrowedBook) SetID(id uint) {
	bb.ID = id
}

// GetStudentID returns the student's ID
func (bb *BorrowedBook) GetStudentID() uint {
	return bb.StudentID
}

// SetStudentID sets the student's ID
func (bb *BorrowedBook) SetStudentID(studentID uint) {
	bb.StudentID = studentID
}

// GetBookID returns the book's ID
func (bb *BorrowedBook) GetBookID() uint {
	return bb.BookID
}

// SetBookID sets the book's ID
func (bb *BorrowedBook) SetBookID(bookID uint) {
	bb.BookID = bookID
}

// GetBorrowDate returns the borrow date
func (bb *BorrowedBook) GetBorrowDate() time.Time {
	return bb.BorrowDate
}

// SetBorrowDate sets the borrow date
func (bb *BorrowedBook) SetBorrowDate(borrowDate time.Time) {
	bb.BorrowDate = borrowDate
}

// GetReturnDate returns the return date (can be nil)
func (bb *BorrowedBook) GetReturnDate() *time.Time {
	return bb.ReturnDate
}

// SetReturnDate sets the return date
func (bb *BorrowedBook) SetReturnDate(returnDate *time.Time) {
	bb.ReturnDate = returnDate
}

// IsReturned checks if the book has been returned by examining the ReturnDate field.
// This method is used to determine the status of a borrowing transaction.
//
// Returns:
//   - true: Book has been returned (ReturnDate is not NULL)
//   - false: Book is still borrowed (ReturnDate is NULL)
func (bb *BorrowedBook) IsReturned() bool {
	return bb.ReturnDate != nil
}

// MarkAsReturned sets the return date to the current timestamp, indicating the book has been returned.
// This method implements the business logic for processing book returns.
//
// Business Logic:
//   - Sets ReturnDate to current time, marking the transaction as complete
//   - Should be called when processing a return request
//   - After calling this method, the corresponding book.available should be incremented
//   - This operation should be performed within a database transaction
func (bb *BorrowedBook) MarkAsReturned() {
	now := time.Now()
	bb.ReturnDate = &now
}
