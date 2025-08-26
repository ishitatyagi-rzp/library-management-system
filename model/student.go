package model

// Student represents a student in the library management system
type Student struct {
	ID       int    `json:"id"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required,len=10,numeric"`
	Password string `json:"password" binding:"required"`
}

// NewStudent creates a new Student instance
func NewStudent(name, email, phone, password string) *Student {
	return &Student{
		Name:     name,
		Email:    email,
		Phone:    phone,
		Password: password,
	}
}

// GetID returns the student's ID
func (s *Student) GetID() int {
	return s.ID
}

// SetID sets the student's ID
func (s *Student) SetID(id int) {
	s.ID = id
}

// GetName returns the student's name
func (s *Student) GetName() string {
	return s.Name
}

// SetName sets the student's name
func (s *Student) SetName(name string) {
	s.Name = name
}

// GetEmail returns the student's email
func (s *Student) GetEmail() string {
	return s.Email
}

// SetEmail sets the student's email
func (s *Student) SetEmail(email string) {
	s.Email = email
}

// GetPhone returns the student's phone
func (s *Student) GetPhone() string {
	return s.Phone
}

// SetPhone sets the student's phone
func (s *Student) SetPhone(phone string) {
	s.Phone = phone
}

// GetPassword returns the student's password
func (s *Student) GetPassword() string {
	return s.Password
}

// SetPassword sets the student's password
func (s *Student) SetPassword(password string) {
	s.Password = password
}
