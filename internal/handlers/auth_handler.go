package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/coolmate/ecommerce-backend/internal/services"
	"github.com/coolmate/ecommerce-backend/internal/utils"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary      Register a new customer
// @Description  Creates a customer account and returns access + refresh tokens.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      services.RegisterRequest  true  "Registration payload"
// @Success      201   {object}  utils.APIResponse
// @Failure      400   {object}  utils.APIResponse
// @Router       /auth/register [post]
func (ah *AuthHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	authResp, err := ah.authService.Register(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", authResp)
}

// Login godoc
// @Summary      Login
// @Description  Authenticates a user and returns access + refresh tokens.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      services.LoginRequest  true  "Login payload"
// @Success      200   {object}  utils.APIResponse
// @Failure      400   {object}  utils.APIResponse
// @Router       /auth/login [post]
func (ah *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	authResp, err := ah.authService.Login(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", authResp)
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Exchange a valid refresh token for a new access + refresh token pair (rotation).
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      map[string]string  true  "{ refreshToken }"
// @Success      200   {object}  utils.APIResponse
// @Failure      401   {object}  utils.APIResponse
// @Router       /auth/refresh [post]
func (ah *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	authResp, err := ah.authService.Refresh(req.RefreshToken)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Token refreshed successfully", authResp)
}

// Logout godoc
// @Summary      Logout
// @Description  Revokes the provided refresh token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      map[string]string  true  "{ refreshToken }"
// @Success      200   {object}  utils.APIResponse
// @Router       /auth/logout [post]
func (ah *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ah.authService.Logout(req.RefreshToken); err != nil {
		utils.InternalServerError(c, "Failed to logout")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Logged out successfully", nil)
}

// RegisterVendor godoc
// @Summary      Register as vendor (STUB)
// @Description  Currently a stub — returns 200 OK with no side-effects. Phase 1 TODO.
// @Tags         Vendor
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  utils.APIResponse
// @Router       /vendor/register [post]
func (ah *AuthHandler) RegisterVendor(c *gin.Context) {
	// To be implemented - vendor registration
	utils.SuccessResponse(c, http.StatusOK, "Vendor registration endpoint", nil)
}
