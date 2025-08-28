package model

// BorrowRequest represents the input for borrowing a book.
// This struct is used for API requests and does not map to a database table.
// It contains validation tags for request validation.
type BorrowRequest struct {
	StudentID uint `json:"student_id" binding:"required,gt=0"` // Must be a valid positive student ID
	BookID    uint `json:"book_id" binding:"required,gt=0"`    // Must be a valid positive book ID
}

// NewBorrowRequest creates a new BorrowRequest instance
func NewBorrowRequest(studentID, bookID uint) *BorrowRequest {
	return &BorrowRequest{
		StudentID: studentID,
		BookID:    bookID,
	}
}

// GetStudentID returns the student ID
func (br *BorrowRequest) GetStudentID() uint {
	return br.StudentID
}

// SetStudentID sets the student ID
func (br *BorrowRequest) SetStudentID(studentID uint) {
	br.StudentID = studentID
}

// GetBookID returns the book ID
func (br *BorrowRequest) GetBookID() uint {
	return br.BookID
}

// SetBookID sets the book ID
func (br *BorrowRequest) SetBookID(bookID uint) {
	br.BookID = bookID
}
