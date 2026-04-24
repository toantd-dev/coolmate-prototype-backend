package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// HealthCheck godoc
// @Summary      Liveness + DB ping
// @Description  Returns 200 when the API is up and the database accepts a ping.
// @Tags         Health
// @Produce      json
// @Success      200  {object}  map[string]string  "{ status: healthy }"
// @Failure      503  {object}  map[string]string  "{ status: unhealthy, error: ... }"
// @Router       /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  "database connection failed",
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  "database ping failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}
