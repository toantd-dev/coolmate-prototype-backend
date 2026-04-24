// BDD-style acceptance tests for the Auth service.
//
// Each TestUS_* function maps 1:1 to a user story from USER_STORIES.md and
// uses sub-tests named "AC<n>/<scenario>" to mark individual acceptance
// criteria. Leading comment carries an `AC:` tag so scripts can compute
// requirement coverage (see scripts/ac-coverage.sh).
//
// This file is ADDITIVE: the existing *_test.go files remain the authority
// for line coverage. This file adds AC-level traceability + Given/When/Then
// documentation of business rules.
package services

import (
	"errors"
	"testing"

	"github.com/coolmate/ecommerce-backend/internal/models"
	"github.com/coolmate/ecommerce-backend/internal/utils"
	"github.com/coolmate/ecommerce-backend/pkg/auth"
	"github.com/coolmate/ecommerce-backend/pkg/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// newAuthSvcBDD builds an AuthService with a fresh MockUserRepository (defined
// in vendor_service_test.go), a real JWTManager with a fixed test secret, and
// a nil-backed cache (Auth service does not call cache directly).
func newAuthSvcBDD(t *testing.T) (*AuthService, *MockUserRepository, *auth.JWTManager) {
	t.Helper()
	repo := new(MockUserRepository)
	jwtm := auth.NewJWTManager("bdd-test-secret", 15, 7)
	cm := cache.NewCacheManager(nil)
	return NewAuthService(repo, jwtm, cm), repo, jwtm
}

// =============================================================================
// US-AUTH-001 · Register a new account
//
// AC: US-AUTH-001 AC1, AC2
// Source: USER_STORIES.md § 1. Authentication
// =============================================================================

func TestUS_AUTH_001_Register(t *testing.T) {
	// AC1
	// Given  no user exists with email "new@example.com"
	// When   Register is called with a valid payload
	// Then   a User row is created with a hashed password
	//        and access + refresh tokens are returned
	t.Run("AC1/happy_path_creates_user_and_issues_tokens", func(t *testing.T) {
		// --- Given ---
		svc, repo, _ := newAuthSvcBDD(t)
		repo.On("GetByEmail", "new@example.com").Return(nil, errors.New("not found"))
		repo.On("Create", mock.AnythingOfType("*models.User")).
			Run(func(args mock.Arguments) {
				u := args.Get(0).(*models.User)
				u.ID = 42 // simulate DB-assigned ID
			}).
			Return(nil)
		repo.On("SaveRefreshToken", mock.AnythingOfType("*models.RefreshToken")).Return(nil)

		req := &RegisterRequest{
			Email:     "new@example.com",
			Password:  "password123",
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "+1-555-0100",
		}

		// --- When ---
		resp, err := svc.Register(req)

		// --- Then ---
		require.NoError(t, err, "registration should succeed for a fresh email")
		require.NotNil(t, resp)
		assert.NotEmpty(t, resp.AccessToken, "access token must be issued")
		assert.NotEmpty(t, resp.RefreshToken, "refresh token must be issued")
		assert.NotEqual(t, resp.AccessToken, resp.RefreshToken, "tokens must differ")
		assert.Equal(t, "new@example.com", resp.User.Email)
		assert.Equal(t, "customer", resp.User.Role, "default role is customer")

		// Validate the password was hashed, not stored plaintext
		createdUser := repo.Calls[1].Arguments.Get(0).(*models.User)
		assert.NotEqual(t, "password123", createdUser.PasswordHash, "password must be hashed")
		assert.True(t, utils.VerifyPassword(createdUser.PasswordHash, "password123"),
			"hash must verify against the original password")

		repo.AssertExpectations(t)
	})

	// AC2
	// Given  a user with email "taken@example.com" already exists
	// When   Register is called with that email
	// Then   a conflict error is returned and no user row is written
	t.Run("AC2/conflict_when_email_already_registered", func(t *testing.T) {
		// --- Given ---
		svc, repo, _ := newAuthSvcBDD(t)
		existing := &models.User{ID: 1, Email: "taken@example.com"}
		repo.On("GetByEmail", "taken@example.com").Return(existing, nil)

		req := &RegisterRequest{
			Email:    "taken@example.com",
			Password: "password123",
		}

		// --- When ---
		resp, err := svc.Register(req)

		// --- Then ---
		assert.Nil(t, resp, "no auth response on conflict")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
		repo.AssertNotCalled(t, "Create")
		repo.AssertNotCalled(t, "SaveRefreshToken")
	})
}

// =============================================================================
// US-AUTH-002 · Login issues a token pair
//
// AC: US-AUTH-002 AC1, AC2
// Source: USER_STORIES.md § 1. Authentication
// =============================================================================

func TestUS_AUTH_002_Login(t *testing.T) {
	const goodPwd = "correctPassword1"

	// AC1
	// Given  a user whose password hash matches "correctPassword1"
	// When   Login is called with matching credentials
	// Then   access + refresh tokens are returned
	//        and a refresh-token row is saved for revocation
	t.Run("AC1/valid_credentials_issue_token_pair", func(t *testing.T) {
		// --- Given ---
		svc, repo, _ := newAuthSvcBDD(t)
		hash, err := utils.HashPassword(goodPwd)
		require.NoError(t, err)
		user := &models.User{
			ID:           7,
			Email:        "alice@example.com",
			PasswordHash: hash,
			Status:       models.UserActive,
			Role:         models.RoleCustomer,
		}
		repo.On("GetByEmail", "alice@example.com").Return(user, nil)
		repo.On("SaveRefreshToken", mock.AnythingOfType("*models.RefreshToken")).Return(nil)

		// --- When ---
		resp, err := svc.Login(&LoginRequest{Email: "alice@example.com", Password: goodPwd})

		// --- Then ---
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.Equal(t, uint(7), resp.User.ID)
		repo.AssertCalled(t, "SaveRefreshToken", mock.AnythingOfType("*models.RefreshToken"))
	})

	// AC2
	// Given  either (a) wrong password or (b) email not in the system
	// When   Login is called
	// Then   an unauthorized error is returned with the same shape
	//        (no user enumeration leak) and nothing is persisted
	t.Run("AC2a/wrong_password_returns_unauthorized", func(t *testing.T) {
		// --- Given ---
		svc, repo, _ := newAuthSvcBDD(t)
		hash, _ := utils.HashPassword(goodPwd)
		repo.On("GetByEmail", "alice@example.com").Return(&models.User{
			ID: 7, Email: "alice@example.com", PasswordHash: hash, Status: models.UserActive,
		}, nil)

		// --- When ---
		resp, err := svc.Login(&LoginRequest{Email: "alice@example.com", Password: "wrong"})

		// --- Then ---
		assert.Nil(t, resp)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid email or password",
			"must use the generic message to prevent enumeration")
		repo.AssertNotCalled(t, "SaveRefreshToken")
	})

	t.Run("AC2b/unknown_email_returns_unauthorized", func(t *testing.T) {
		// --- Given ---
		svc, repo, _ := newAuthSvcBDD(t)
		repo.On("GetByEmail", "ghost@example.com").Return(nil, errors.New("not found"))

		// --- When ---
		resp, err := svc.Login(&LoginRequest{Email: "ghost@example.com", Password: goodPwd})

		// --- Then ---
		assert.Nil(t, resp)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid email or password",
			"identical message to wrong-password case prevents enumeration")
		repo.AssertNotCalled(t, "SaveRefreshToken")
	})
}

// =============================================================================
// US-AUTH-003 · Refresh rotates, Logout revokes
//
// AC: US-AUTH-003 AC1, AC2, AC3
// Source: USER_STORIES.md § 1. Authentication
// =============================================================================

func TestUS_AUTH_003_RefreshAndLogout(t *testing.T) {
	// AC1
	// Given  a valid refresh token for user 42
	// When   Refresh is called
	// Then   a fresh access+refresh pair is returned
	//        and the old refresh token hash is revoked (rotation)
	t.Run("AC1/refresh_rotates_token_pair", func(t *testing.T) {
		// --- Given ---
		svc, repo, jwtm := newAuthSvcBDD(t)
		user := &models.User{ID: 42, Email: "rotate@example.com", Role: models.RoleCustomer}
		oldRefresh, err := jwtm.GenerateRefreshToken(user)
		require.NoError(t, err)
		repo.On("GetByID", uint(42)).Return(user, nil)
		repo.On("SaveRefreshToken", mock.AnythingOfType("*models.RefreshToken")).Return(nil)
		repo.On("RevokeRefreshToken", mock.AnythingOfType("string")).Return(nil)

		// --- When ---
		resp, err := svc.Refresh(oldRefresh)

		// --- Then ---
		//
		// Rotation is a behavior, not string inequality: the OLD refresh hash
		// is revoked and a NEW hash is persisted. Note that JWT `iat` has
		// 1-second resolution, so two refreshes in the same second can produce
		// the same JWT string. The contract that matters is the repo call
		// pattern — that's what the server uses to reject a leaked token.
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.NotEmpty(t, resp.AccessToken, "new access token must be issued")
		assert.NotEmpty(t, resp.RefreshToken, "new refresh token must be issued")
		repo.AssertCalled(t, "SaveRefreshToken", mock.AnythingOfType("*models.RefreshToken"))
		repo.AssertCalled(t, "RevokeRefreshToken", mock.AnythingOfType("string"))
	})

	// AC2
	// Given  a refresh token that was never issued (or has been revoked)
	// When   Refresh is called
	// Then   an error is returned and no tokens are issued
	t.Run("AC2/refresh_rejects_invalid_token", func(t *testing.T) {
		// --- Given ---
		svc, repo, _ := newAuthSvcBDD(t)
		garbage := "not.a.valid.jwt"

		// --- When ---
		resp, err := svc.Refresh(garbage)

		// --- Then ---
		assert.Nil(t, resp)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid refresh token")
		repo.AssertNotCalled(t, "GetByID")
		repo.AssertNotCalled(t, "SaveRefreshToken")
	})

	// AC3
	// Given  a refresh token belonging to user 42
	// When   Logout is called with that token
	// Then   the token hash is revoked at the repository
	t.Run("AC3/logout_revokes_refresh_token", func(t *testing.T) {
		// --- Given ---
		svc, repo, jwtm := newAuthSvcBDD(t)
		user := &models.User{ID: 42, Email: "logout@example.com", Role: models.RoleCustomer}
		token, err := jwtm.GenerateRefreshToken(user)
		require.NoError(t, err)
		repo.On("RevokeRefreshToken", mock.AnythingOfType("string")).Return(nil)

		// --- When ---
		err = svc.Logout(token)

		// --- Then ---
		require.NoError(t, err)
		repo.AssertCalled(t, "RevokeRefreshToken", mock.AnythingOfType("string"))
	})
}
