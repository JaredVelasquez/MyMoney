package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserClaims represents JWT claims with user information
type UserClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// RefreshTokenClaims represents claims for refresh tokens
type RefreshTokenClaims struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	TokenID string `json:"token_id"` // Para identificar de forma única el refresh token
	jwt.RegisteredClaims
}

// TokenService handles token generation and validation
type TokenService struct{}

// TokenPair represents an access token and refresh token pair
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// NewTokenService creates a new TokenService
func NewTokenService() *TokenService {
	return &TokenService{}
}

// ValidateToken validates a JWT token
func (s *TokenService) ValidateToken(tokenString string) (*UserClaims, error) {
	// Get JWT secret from environment
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET not set")
	}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Check if token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token
func (s *TokenService) ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	// Get JWT refresh secret from environment
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		secret = os.Getenv("JWT_SECRET") // Fallback to regular secret if refresh secret not set
		if secret == "" {
			return nil, fmt.Errorf("JWT_SECRET not set")
		}
	}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Extract claims
	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok {
		return nil, fmt.Errorf("invalid refresh token claims")
	}

	return claims, nil
}

// GenerateToken generates a JWT token for a user
func (s *TokenService) GenerateToken(userID, email string) (string, error) {
	// Get JWT secret from environment
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET not set")
	}

	// Create claims with expiration
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
	claims := &UserClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "MyMoneyBackend",
			Subject:   userID,
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateTokenPair generates an access token and refresh token pair
func (s *TokenService) GenerateTokenPair(userID, email string) (*TokenPair, error) {
	// Generate access token
	accessToken, err := s.GenerateToken(userID, email)
	if err != nil {
		return nil, err
	}

	// Get JWT refresh secret from environment
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		refreshSecret = os.Getenv("JWT_SECRET") // Fallback to regular secret if refresh secret not set
		if refreshSecret == "" {
			return nil, fmt.Errorf("JWT_SECRET not set")
		}
	}

	// Create refresh token with longer expiration
	refreshExpirationTime := time.Now().Add(7 * 24 * time.Hour) // Refresh token valid for 7 days
	tokenID := fmt.Sprintf("%s-%d", userID, time.Now().Unix())  // Unique token ID

	refreshClaims := &RefreshTokenClaims{
		UserID:  userID,
		Email:   email,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "MyMoneyBackend",
			Subject:   userID,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(refreshSecret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
	}, nil
}

// RefreshAccessToken generates a new access token from a valid refresh token
func (s *TokenService) RefreshAccessToken(refreshTokenString string) (string, error) {
	// Validate refresh token
	claims, err := s.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", err
	}

	// Generate new access token
	return s.GenerateToken(claims.UserID, claims.Email)
}
