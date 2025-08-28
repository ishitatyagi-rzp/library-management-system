package model

import (
	"time"

	"gorm.io/gorm"
)

// Book represents a book entity in the library management system.
// This struct maps to the 'books' table and includes inventory management fields
// to track both total quantity and currently available copies.
//
// Database Table: books
// Business Logic:
//   - Quantity: Total copies owned by the library (immutable after creation)
//   - Available: Current number of copies available for borrowing (changes with borrows/returns)
//   - Available should never exceed Quantity
//   - Available should never be negative
//
// GORM Tags: Used for database schema generation, constraints, and relationships
type Book struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`                                // Primary key, auto-generated
	Title     string         `json:"title" binding:"required" gorm:"type:varchar(200);not null;index"`  // Book title, indexed for searching
	Author    string         `json:"author" binding:"required" gorm:"type:varchar(100);not null;index"` // Book author, indexed for searching
	Category  string         `json:"category" gorm:"type:varchar(50);index"`                            // Book category/genre, indexed for filtering
	Quantity  int            `json:"quantity" binding:"gte=1" gorm:"not null;check:quantity >= 1"`      // Total copies owned (must be >= 1)
	Available int            `json:"available" gorm:"not null;check:available >= 0"`                    // Currently available copies for borrowing
	AddedAt   time.Time      `json:"added_at" gorm:"autoCreateTime"`                                    // Automatically set when record is created
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`                                  // Automatically updated when record changes
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`                                                    // For soft deletion support

	// Relationships
	BorrowedBooks []BorrowedBook `json:"-" gorm:"foreignKey:BookID"` // One-to-many relationship with borrowed_books
}

// NewBook creates a new Book instance with GORM-compatible fields.
// GORM will automatically set AddedAt, UpdatedAt timestamps.
func NewBook(title, author, category string, quantity int) *Book {
	return &Book{
		Title:     title,
		Author:    author,
		Category:  category,
		Quantity:  quantity,
		Available: quantity, // Initially all books are available
		// AddedAt and UpdatedAt will be set automatically by GORM
	}
}

// GetID returns the book's ID
func (b *Book) GetID() uint {
	return b.ID
}

// SetID sets the book's ID
func (b *Book) SetID(id uint) {
	b.ID = id
}

// GetTitle returns the book's title
func (b *Book) GetTitle() string {
	return b.Title
}

// SetTitle sets the book's title
func (b *Book) SetTitle(title string) {
	b.Title = title
}

// GetAuthor returns the book's author
func (b *Book) GetAuthor() string {
	return b.Author
}

// SetAuthor sets the book's author
func (b *Book) SetAuthor(author string) {
	b.Author = author
}

// GetCategory returns the book's category
func (b *Book) GetCategory() string {
	return b.Category
}

// SetCategory sets the book's category
func (b *Book) SetCategory(category string) {
	b.Category = category
}

// GetQuantity returns the book's total quantity
func (b *Book) GetQuantity() int {
	return b.Quantity
}

// SetQuantity sets the book's total quantity
func (b *Book) SetQuantity(quantity int) {
	b.Quantity = quantity
}

// GetAvailable returns the book's available quantity
func (b *Book) GetAvailable() int {
	return b.Available
}

// SetAvailable sets the book's available quantity
func (b *Book) SetAvailable(available int) {
	b.Available = available
}

// GetAddedAt returns when the book was added
func (b *Book) GetAddedAt() time.Time {
	return b.AddedAt
}

// SetAddedAt sets when the book was added
func (b *Book) SetAddedAt(addedAt time.Time) {
	b.AddedAt = addedAt
}

// IsAvailable checks if the book has copies available for borrowing.
// Returns true if at least one copy is available, false otherwise.
// This is a convenience method used before allowing borrow operations.
func (b *Book) IsAvailable() bool {
	return b.Available > 0
}

// BorrowBook decreases the available count by 1 when a book is borrowed.
// This method implements the core business logic for book borrowing inventory management.
//
// Returns:
//   - true: Successfully decremented available count (book was available)
//   - false: No copies available for borrowing (available count is 0)
//
// Business Rules:
//   - Only decrements if copies are available (Available > 0)
//   - Available count cannot go below 0
//   - This method should be called within a database transaction
func (b *Book) BorrowBook() bool {
	if b.Available > 0 {
		b.Available--
		return true
	}
	return false
}

// ReturnBook increases the available count by 1 when a book is returned.
// This method implements the core business logic for book return inventory management.
//
// Business Rules:
//   - Only increments if returned count won't exceed total quantity
//   - Available count cannot exceed Quantity (total owned copies)
//   - This method should be called within a database transaction
//
// Note: This method assumes the book was actually borrowed by the student.
// Validation should be done at the handler level before calling this method.
func (b *Book) ReturnBook() {
	if b.Available < b.Quantity {
		b.Available++
	}
}
