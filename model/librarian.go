package model

// Librarian represents a librarian entity in the library management system.
// Librarians have administrative privileges to manage books, students, and library operations.
//
// GORM Tags: Used for database schema generation and constraints
type Librarian struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`                                           // Primary key, auto-generated
	Name     string `json:"name" binding:"required" gorm:"type:varchar(100);not null"`                    // Full name, cannot be empty
	Email    string `json:"email" binding:"required,email" gorm:"type:varchar(100);uniqueIndex;not null"` // Must be valid and unique
	Password string `json:"password" binding:"required" gorm:"type:varchar(255);not null"`                // Hashed password for storage
}

// NewLibrarian creates a new Librarian instance
func NewLibrarian(name, email, password string) *Librarian {
	return &Librarian{
		Name:     name,
		Email:    email,
		Password: password,
	}
}

// GetID returns the librarian's ID
func (l *Librarian) GetID() uint {
	return l.ID
}

// SetID sets the librarian's ID
func (l *Librarian) SetID(id uint) {
	l.ID = id
}

// GetName returns the librarian's name
func (l *Librarian) GetName() string {
	return l.Name
}

// SetName sets the librarian's name
func (l *Librarian) SetName(name string) {
	l.Name = name
}

// GetEmail returns the librarian's email
func (l *Librarian) GetEmail() string {
	return l.Email
}

// SetEmail sets the librarian's email
func (l *Librarian) SetEmail(email string) {
	l.Email = email
}

// GetPassword returns the librarian's password
func (l *Librarian) GetPassword() string {
	return l.Password
}

// SetPassword sets the librarian's password
func (l *Librarian) SetPassword(password string) {
	l.Password = password
}
