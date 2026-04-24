package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/coolmate/ecommerce-backend/internal/models"
	"github.com/coolmate/ecommerce-backend/internal/utils"
)

func RequireRole(allowedRoles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := GetUserRole(c)
		if role == "" {
			utils.Unauthorized(c, "user not authenticated")
			c.Abort()
			return
		}

		authorized := false
		for _, allowedRole := range allowedRoles {
			if string(allowedRole) == role {
				authorized = true
				break
			}
		}

		if !authorized {
			utils.Forbidden(c, "insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return RequireRole(models.RoleSuperAdmin, models.RoleAdmin)
}

func RequireVendor() gin.HandlerFunc {
	return RequireRole(models.RoleVendor)
}

func RequireSuperAdmin() gin.HandlerFunc {
	return RequireRole(models.RoleSuperAdmin)
}
