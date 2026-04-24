package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/coolmate/ecommerce-backend/internal/models"
)

type JWTManager struct {
	secret                   string
	accessTokenExpireMinutes int
	refreshTokenExpireDays   int
}

type Claims struct {
	UserID   uint
	Email    string
	Role     string
	jwt.RegisteredClaims
}

func NewJWTManager(secret string, accessExpireMinutes, refreshExpireDays int) *JWTManager {
	return &JWTManager{
		secret:                   secret,
		accessTokenExpireMinutes: accessExpireMinutes,
		refreshTokenExpireDays:   refreshExpireDays,
	}
}

func (jm *JWTManager) GenerateAccessToken(user *models.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jm.accessTokenExpireMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jm.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (jm *JWTManager) GenerateRefreshToken(user *models.User) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", user.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, jm.refreshTokenExpireDays)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jm.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

func (jm *JWTManager) VerifyAccessToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jm.secret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}

func (jm *JWTManager) VerifyRefreshToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jm.secret), nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid refresh token: %w", err)
	}

	claims := token.Claims.(*jwt.RegisteredClaims)
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return 0, fmt.Errorf("refresh token expired")
	}

	var userID uint
	_, err = fmt.Sscanf(claims.Subject, "%d", &userID)
	if err != nil {
		return 0, fmt.Errorf("invalid user id in token: %w", err)
	}

	return userID, nil
}
