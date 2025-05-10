package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"mi-app-backend/internal/domain"

	"github.com/google/uuid"
)

// CategoryRepository implementa la interfaz repositories.CategoryRepository
type CategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository crea un nuevo repositorio de categorías
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

// Create crea una nueva categoría en la base de datos
func (r *CategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	if category.ID == "" {
		category.ID = uuid.New().String()
	}

	now := time.Now()
	if category.CreatedAt.IsZero() {
		category.CreatedAt = now
	}
	category.UpdatedAt = now

	query := `
		INSERT INTO categories (id, name, description, type, color, icon, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		category.ID,
		category.Name,
		category.Description,
		category.Type,
		category.Color,
		category.Icon,
		category.UserID,
		category.CreatedAt,
		category.UpdatedAt,
	)

	return err
}

// GetByID obtiene una categoría por su ID
func (r *CategoryRepository) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	query := `
		SELECT id, name, description, type, color, icon, user_id, created_at, updated_at
		FROM categories
		WHERE id = $1
	`

	var category domain.Category
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.Type,
		&category.Color,
		&category.Icon,
		&category.UserID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No se encontró la categoría
		}
		return nil, err
	}

	return &category, nil
}

// GetByUserID obtiene todas las categorías de un usuario
func (r *CategoryRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Category, error) {
	query := `
		SELECT id, name, description, type, color, icon, user_id, created_at, updated_at
		FROM categories
		WHERE user_id = $1
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var category domain.Category
		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.Type,
			&category.Color,
			&category.Icon,
			&category.UserID,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// GetByType obtiene todas las categorías de un usuario por tipo
func (r *CategoryRepository) GetByType(ctx context.Context, userID string, categoryType domain.CategoryType) ([]*domain.Category, error) {
	query := `
		SELECT id, name, description, type, color, icon, user_id, created_at, updated_at
		FROM categories
		WHERE user_id = $1 AND type = $2
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, categoryType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var category domain.Category
		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.Type,
			&category.Color,
			&category.Icon,
			&category.UserID,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// Update actualiza una categoría existente
func (r *CategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	category.UpdatedAt = time.Now()

	query := `
		UPDATE categories
		SET name = $1, description = $2, type = $3, color = $4, icon = $5, updated_at = $6
		WHERE id = $7
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		category.Name,
		category.Description,
		category.Type,
		category.Color,
		category.Icon,
		category.UpdatedAt,
		category.ID,
	)

	return err
}

// Delete elimina una categoría
func (r *CategoryRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
