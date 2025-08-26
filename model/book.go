package model

import "time"

// Book represents a book in the library management system
type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title" binding:"required"`
	Author    string    `json:"author" binding:"required"`
	Category  string    `json:"category"`
	Quantity  int       `json:"quantity" binding:"gte=1"` // total copies of this book the library has
	Available int       `json:"available"`
	AddedAt   time.Time `json:"added_at"`
}

// NewBook creates a new Book instance
func NewBook(title, author, category string, quantity int) *Book {
	return &Book{
		Title:     title,
		Author:    author,
		Category:  category,
		Quantity:  quantity,
		Available: quantity, // Initially all books are available
		AddedAt:   time.Now(),
	}
}

// GetID returns the book's ID
func (b *Book) GetID() int {
	return b.ID
}

// SetID sets the book's ID
func (b *Book) SetID(id int) {
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

// IsAvailable checks if the book is available for borrowing
func (b *Book) IsAvailable() bool {
	return b.Available > 0
}

// BorrowBook decreases the available count when a book is borrowed
func (b *Book) BorrowBook() bool {
	if b.Available > 0 {
		b.Available--
		return true
	}
	return false
}

// ReturnBook increases the available count when a book is returned
func (b *Book) ReturnBook() {
	if b.Available < b.Quantity {
		b.Available++
	}
}
