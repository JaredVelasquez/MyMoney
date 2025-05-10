package api

import (
	"context"
	"net/http"

	"mi-app-backend/internal/application/category"
	"mi-app-backend/internal/domain"

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
		req.Type,
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
		req.Type,
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

// GetCategoriesByType returns all categories of a specific type
func (h *CategoryHandler) GetCategoriesByType(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	typeStr := c.Param("type")
	if typeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category type is required"})
		return
	}

	categoryType := domain.CategoryType(typeStr)
	if categoryType != domain.CategoryTypeIncome && categoryType != domain.CategoryTypeExpense {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category type"})
		return
	}

	categories, err := h.categoryService.GetCategoriesByType(context.Background(), userID, categoryType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}
