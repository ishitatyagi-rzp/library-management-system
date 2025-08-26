package model

// LoginInput represents the input for login requests
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// NewLoginInput creates a new LoginInput instance
func NewLoginInput(email, password string) *LoginInput {
	return &LoginInput{
		Email:    email,
		Password: password,
	}
}

// GetEmail returns the email
func (l *LoginInput) GetEmail() string {
	return l.Email
}

// SetEmail sets the email
func (l *LoginInput) SetEmail(email string) {
	l.Email = email
}

// GetPassword returns the password
func (l *LoginInput) GetPassword() string {
	return l.Password
}

// SetPassword sets the password
func (l *LoginInput) SetPassword(password string) {
	l.Password = password
}
