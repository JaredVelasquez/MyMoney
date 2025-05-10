package repository

import (
	"database/sql"
	"errors"
	"time"

	"mi-app-backend/internal/domain"

	"github.com/google/uuid"
)

// UserRepository implementa la interfaz app.UserRepository
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository crea un nuevo repositorio de usuarios
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create crea un nuevo usuario en la base de datos
func (r *UserRepository) Create(user *domain.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	user.UpdatedAt = now

	query := `
		INSERT INTO users (id, email, name, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(
		query,
		user.ID,
		user.Email,
		user.Name,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

// GetByID obtiene un usuario por su ID
func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	query := `
		SELECT id, email, name, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail obtiene un usuario por su email
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, email, name, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user domain.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// Update actualiza la información de un usuario
func (r *UserRepository) Update(user *domain.User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET email = $1, name = $2, password = $3, updated_at = $4
		WHERE id = $5
	`

	result, err := r.db.Exec(
		query,
		user.Email,
		user.Name,
		user.Password,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// Delete elimina un usuario
func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
