package model

import "time"

// BorrowedBook represents a borrowed book record in the library management system
type BorrowedBook struct {
	ID         int        `json:"id"`
	StudentID  int        `json:"student_id"`
	BookID     int        `json:"book_id"`
	BorrowDate time.Time  `json:"borrow_date"`
	ReturnDate *time.Time `json:"return_date,omitempty"`
}

// NewBorrowedBook creates a new BorrowedBook instance
func NewBorrowedBook(studentID, bookID int) *BorrowedBook {
	return &BorrowedBook{
		StudentID:  studentID,
		BookID:     bookID,
		BorrowDate: time.Now(),
	}
}

// GetID returns the borrowed book record's ID
func (bb *BorrowedBook) GetID() int {
	return bb.ID
}

// SetID sets the borrowed book record's ID
func (bb *BorrowedBook) SetID(id int) {
	bb.ID = id
}

// GetStudentID returns the student's ID
func (bb *BorrowedBook) GetStudentID() int {
	return bb.StudentID
}

// SetStudentID sets the student's ID
func (bb *BorrowedBook) SetStudentID(studentID int) {
	bb.StudentID = studentID
}

// GetBookID returns the book's ID
func (bb *BorrowedBook) GetBookID() int {
	return bb.BookID
}

// SetBookID sets the book's ID
func (bb *BorrowedBook) SetBookID(bookID int) {
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

// IsReturned checks if the book has been returned
func (bb *BorrowedBook) IsReturned() bool {
	return bb.ReturnDate != nil
}

// MarkAsReturned marks the book as returned with the current timestamp
func (bb *BorrowedBook) MarkAsReturned() {
	now := time.Now()
	bb.ReturnDate = &now
}
