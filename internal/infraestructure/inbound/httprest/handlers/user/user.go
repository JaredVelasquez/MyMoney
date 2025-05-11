package api

import (
	"net/http"

	"MyMoneyBackend/internal/application/auth"
	services "MyMoneyBackend/internal/application/user"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService services.UserService
	authService auth.TokenService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService services.UserService, authService auth.TokenService) *UserHandler {
	return &UserHandler{
		userService: userService,
		authService: authService,
	}
}

// RegisterRequest represents the request for registering a new user
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents the request for logging in
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest represents the request for updating a user
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// ChangePasswordRequest represents the request for changing the password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// UserResponse represents the response with user data
type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

// RefreshTokenRequest represents the refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse represents the refresh token response
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// Register godoc
// @Summary Registrar un nuevo usuario
// @Description Registra un nuevo usuario en el sistema
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "Datos del usuario"
// @Success 201 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.RegisterUser(req.Email, req.Name, req.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token pair
	tokenPair, err := h.authService.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error generating token"})
		return
	}

	c.JSON(http.StatusCreated, LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		User: UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

// Login godoc
// @Summary Login
// @Description Authenticates a user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Access credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// Generate JWT token pair
	tokenPair, err := h.authService.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error generating token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		User: UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

// GetMe godoc
// @Summary Obtener usuario actual
// @Description Retorna los datos del usuario autenticado
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} UserResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}

// UpdateUser godoc
// @Summary Update user
// @Description Updates the authenticated user's information
// @Tags users
// @Accept json
// @Produce json
// @Param userData body UpdateUserRequest true "Updated user data"
// @Security Bearer
// @Success 200 {object} UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /users/update [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateUser(userID, req.Email, req.Name)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}

// ChangePassword godoc
// @Summary Change password
// @Description Changes the authenticated user's password
// @Tags users
// @Accept json
// @Produce json
// @Param passwordData body ChangePasswordRequest true "Password data"
// @Security Bearer
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/change-password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.ChangePassword(userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// RefreshToken godoc
// @Summary Refresh token
// @Description Refreshes an access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refreshToken body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} RefreshTokenResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/refresh-token [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate new access token
	newAccessToken, err := h.authService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, RefreshTokenResponse{
		AccessToken: newAccessToken,
	})
}
