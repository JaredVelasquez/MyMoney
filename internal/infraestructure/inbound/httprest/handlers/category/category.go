package api

import (
	"context"
	"net/http"

	"MyMoneyBackend/internal/application/category"
	"MyMoneyBackend/internal/domain"

	"github.com/gin-gonic/gin"
)

// CategoryHandler handles category-related HTTP requests
type CategoryHandler struct {
	categoryService *category.Service
}

// NewCategoryHandler creates a new CategoryHandler
func NewCategoryHandler(categoryService *category.Service) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory handles category creation
// @Summary Crear una nueva categoría
// @Description Crea una nueva categoría para el usuario autenticado
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param category body domain.CreateCategoryRequest true "Datos de la categoría a crear"
// @Success 201 {object} domain.Category
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req domain.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryService.CreateCategory(
		context.Background(),
		req.Name,
		req.Description,
		req.Icon,
		req.Color,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// GetUserCategories returns all categories for the current user
// @Summary Obtener todas las categorías del usuario
// @Description Retorna todas las categorías del usuario autenticado
// @Tags categories
// @Produce json
// @Security Bearer
// @Success 200 {array} domain.Category
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/categories [get]
func (h *CategoryHandler) GetUserCategories(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	categories, err := h.categoryService.GetCategoriesByUserID(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategory returns a specific category
// @Summary Obtener una categoría específica
// @Description Retorna una categoría específica por su ID
// @Tags categories
// @Produce json
// @Security Bearer
// @Param id path string true "ID de la categoría"
// @Success 200 {object} domain.Category
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	categoryID := c.Param("id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category ID is required"})
		return
	}

	category, err := h.categoryService.GetCategoryByID(context.Background(), categoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	// Check if the category belongs to the user
	if category.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// UpdateCategory updates a category
// @Summary Actualizar una categoría
// @Description Actualiza una categoría existente
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "ID de la categoría"
// @Param category body domain.UpdateCategoryRequest true "Datos de la categoría a actualizar"
// @Success 200 {object} domain.Category
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	categoryID := c.Param("id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category ID is required"})
		return
	}

	var req domain.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryService.UpdateCategory(
		context.Background(),
		categoryID,
		req.Name,
		req.Description,
		req.Icon,
		req.Color,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory deletes a category
// @Summary Eliminar una categoría
// @Description Elimina una categoría existente
// @Tags categories
// @Security Bearer
// @Param id path string true "ID de la categoría"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	categoryID := c.Param("id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category ID is required"})
		return
	}

	err := h.categoryService.DeleteCategory(context.Background(), categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}
