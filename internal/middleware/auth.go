package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/coolmate/ecommerce-backend/internal/utils"
	"github.com/coolmate/ecommerce-backend/pkg/auth"
)

const (
	authorizationHeader = "Authorization"
	bearerScheme        = "Bearer"
	userIDCtxKey        = "userID"
	emailCtxKey         = "email"
	roleCtxKey          = "role"
)

func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeader)
		if authHeader == "" {
			utils.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != bearerScheme {
			utils.Unauthorized(c, "invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := jwtManager.VerifyAccessToken(token)
		if err != nil {
			utils.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		c.Set(userIDCtxKey, claims.UserID)
		c.Set(emailCtxKey, claims.Email)
		c.Set(roleCtxKey, claims.Role)
		c.Next()
	}
}

func GetUserID(c *gin.Context) uint {
	val, exists := c.Get(userIDCtxKey)
	if !exists {
		return 0
	}
	userID, ok := val.(uint)
	if !ok {
		return 0
	}
	return userID
}

func GetUserEmail(c *gin.Context) string {
	val, exists := c.Get(emailCtxKey)
	if !exists {
		return ""
	}
	email, ok := val.(string)
	if !ok {
		return ""
	}
	return email
}

func GetUserRole(c *gin.Context) string {
	val, exists := c.Get(roleCtxKey)
	if !exists {
		return ""
	}
	role, ok := val.(string)
	if !ok {
		return ""
	}
	return role
}
