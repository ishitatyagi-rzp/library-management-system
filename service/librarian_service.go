package service

import (
	"errors"

	"library-management-system/constants"
	"library-management-system/dao"
	"library-management-system/model"
	"library-management-system/utils"
)

// LibrarianService interface defines business logic operations for librarians
type LibrarianService interface {
	// Authentication
	RegisterLibrarian(librarian *model.Librarian) error
	LoginLibrarian(email, password string) (*LibrarianLoginResponse, error)
}

// LibrarianLoginResponse represents the response for successful librarian login
type LibrarianLoginResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// librarianServiceImpl implements LibrarianService interface
type librarianServiceImpl struct {
	librarianDAO dao.LibrarianDAO
}

// NewLibrarianService creates a new LibrarianService instance
func NewLibrarianService(librarianDAO dao.LibrarianDAO) LibrarianService {
	return &librarianServiceImpl{
		librarianDAO: librarianDAO,
	}
}

// RegisterLibrarian handles librarian registration with password hashing
func (s *librarianServiceImpl) RegisterLibrarian(librarian *model.Librarian) error {
	// Hash the password
	hashedPassword, err := utils.HashPassword(librarian.GetPassword())
	if err != nil {
		return errors.New(constants.HashPasswordError)
	}

	librarian.SetPassword(hashedPassword)

	if err := s.librarianDAO.Create(librarian); err != nil {
		return errors.New(constants.RegisterLibrarianError)
	}

	return nil
}

// LoginLibrarian handles librarian authentication
func (s *librarianServiceImpl) LoginLibrarian(email, password string) (*LibrarianLoginResponse, error) {
	librarian, err := s.librarianDAO.GetByEmail(email)
	if err != nil {
		// Don't reveal whether email exists for security
		return nil, errors.New(constants.InvalidCredentialsError)
	}

	// Verify password
	if !utils.CheckPasswordHash(password, librarian.GetPassword()) {
		return nil, errors.New(constants.InvalidCredentialsError)
	}

	return &LibrarianLoginResponse{
		ID:    librarian.GetID(),
		Name:  librarian.GetName(),
		Email: librarian.GetEmail(),
	}, nil
}
