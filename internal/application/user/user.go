package services

import (
	"errors"
	"log"
	"time"

	"MyMoneyBackend/internal/domain"
	"MyMoneyBackend/internal/domain/ports/app"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

// UserService handles user business logic
type UserService struct {
	userRepo app.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo app.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// RegisterUser registers a new user
func (s *UserService) RegisterUser(email, name, password string) (*domain.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, errors.New("error creating user")
	}

	// Create user
	user := &domain.User{
		ID:        uuid.New().String(),
		Email:     email,
		Name:      name,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(user); err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, errors.New("error creating user")
	}

	// Don't return the password
	user.Password = ""
	return user, nil
}

// GetUserByID gets a user by ID
func (s *UserService) GetUserByID(id string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Don't return the password
	user.Password = ""
	return user, nil
}

// GetUserByEmail gets a user by email
func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	// Don't return the password
	user.Password = ""
	return user, nil
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(id, email, name string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Only update fields that were provided
	if email != "" {
		user.Email = email
	}
	if name != "" {
		user.Name = name
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Don't return the password
	user.Password = ""
	return user, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(id, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error updating password")
	}

	user.Password = string(hashedPassword)

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	return nil
}

// AuthenticateUser authenticates a user
func (s *UserService) AuthenticateUser(email, password string) (*domain.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Don't return the password
	user.Password = ""
	return user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id string) error {
	return s.userRepo.Delete(id)
}
