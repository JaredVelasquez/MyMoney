package domain

import (
	"errors"
	"time"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Password  string    `json:"-"` // Password is not exposed in JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate validates the user data
func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

// registerRequest represents the user registration request
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// loginRequest represents the user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// updateUserRequest represents the update user request
type UpdateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// changePasswordRequest represents the change password request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}
