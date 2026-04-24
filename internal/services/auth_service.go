package services

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/coolmate/ecommerce-backend/internal/models"
	"github.com/coolmate/ecommerce-backend/internal/repositories"
	"github.com/coolmate/ecommerce-backend/pkg/auth"
	"github.com/coolmate/ecommerce-backend/pkg/cache"
	"github.com/coolmate/ecommerce-backend/internal/utils"
)

type IAuthService interface {
	Register(req *RegisterRequest) (*AuthResponse, error)
	Login(req *LoginRequest) (*AuthResponse, error)
	Refresh(refreshToken string) (*AuthResponse, error)
	Logout(refreshToken string) error
}

type AuthService struct {
	userRepo   repositories.IUserRepository
	jwtManager *auth.JWTManager
	cache      *cache.CacheManager
}

func NewAuthService(
	userRepo repositories.IUserRepository,
	jwtManager *auth.JWTManager,
	cache *cache.CacheManager,
) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		cache:      cache,
	}
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	User         *UserDTO    `json:"user"`
}

type UserDTO struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
}

func (as *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Check if user exists
	existingUser, _ := as.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user already exists")
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		Role:         models.RoleCustomer,
		Status:       models.UserActive,
	}

	if err := as.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	accessToken, err := as.jwtManager.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := as.jwtManager.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token hash
	tokenHash := hashToken(refreshToken)
	rt := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().AddDate(0, 0, 7),
	}
	if err := as.userRepo.SaveRefreshToken(rt); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &UserDTO{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      string(user.Role),
		},
	}, nil
}

func (as *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	user, err := as.userRepo.GetByEmail(req.Email)
	if err != nil || user == nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	if !utils.VerifyPassword(user.PasswordHash, req.Password) {
		return nil, fmt.Errorf("invalid email or password")
	}

	if user.Status != models.UserActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Generate tokens
	accessToken, err := as.jwtManager.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := as.jwtManager.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token hash
	tokenHash := hashToken(refreshToken)
	rt := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().AddDate(0, 0, 7),
	}
	if err := as.userRepo.SaveRefreshToken(rt); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &UserDTO{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      string(user.Role),
		},
	}, nil
}

func (as *AuthService) Refresh(refreshToken string) (*AuthResponse, error) {
	userID, err := as.jwtManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get user
	user, err := as.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Generate new access token
	accessToken, err := as.jwtManager.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate new refresh token
	newRefreshToken, err := as.jwtManager.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Save new refresh token
	tokenHash := hashToken(newRefreshToken)
	rt := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().AddDate(0, 0, 7),
	}
	if err := as.userRepo.SaveRefreshToken(rt); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	// Revoke old refresh token
	oldTokenHash := hashToken(refreshToken)
	as.userRepo.RevokeRefreshToken(oldTokenHash)

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User: &UserDTO{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      string(user.Role),
		},
	}, nil
}

func (as *AuthService) Logout(refreshToken string) error {
	tokenHash := hashToken(refreshToken)
	return as.userRepo.RevokeRefreshToken(tokenHash)
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}
