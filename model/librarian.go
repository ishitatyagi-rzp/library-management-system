package model

// Librarian represents a librarian in the library management system
type Librarian struct {
	ID       int    `json:"id"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
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
func (l *Librarian) GetID() int {
	return l.ID
}

// SetID sets the librarian's ID
func (l *Librarian) SetID(id int) {
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
