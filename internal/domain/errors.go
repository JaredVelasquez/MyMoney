package domain

import "errors"

// Errores comunes para entidades
var (
	ErrEmptyID          = errors.New("el ID no puede estar vacío")
	ErrEmptyName        = errors.New("el nombre no puede estar vacío")
	ErrEmptyUserID      = errors.New("el ID de usuario no puede estar vacío")
	ErrEmptyEmail       = errors.New("el email no puede estar vacío")
	ErrInvalidEmail     = errors.New("el email no es válido")
	ErrEmptyPassword    = errors.New("la contraseña no puede estar vacía")
	ErrPasswordTooShort = errors.New("la contraseña debe tener al menos 6 caracteres")
	ErrInvalidAmount    = errors.New("el monto debe ser mayor que cero")
	ErrEmptyCategoryID  = errors.New("el ID de categoría no puede estar vacío")
)
