package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword_Success(t *testing.T) {
	password := "securePassword123"

	hash, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestHashPassword_Different(t *testing.T) {
	password := "myPassword"

	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, hash1, hash2)
}

func TestHashPassword_Empty(t *testing.T) {
	hash, err := HashPassword("")

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestHashPassword_Special(t *testing.T) {
	password := "Pass@123!Spec"

	hash, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestVerifyPassword_Correct(t *testing.T) {
	password := "correctPassword123"

	hash, _ := HashPassword(password)
	result := VerifyPassword(hash, password)

	assert.True(t, result)
}

func TestVerifyPassword_Incorrect(t *testing.T) {
	password := "correctPassword"
	wrongPassword := "wrongPassword"

	hash, _ := HashPassword(password)
	result := VerifyPassword(hash, wrongPassword)

	assert.False(t, result)
}

func TestVerifyPassword_Empty(t *testing.T) {
	password := "somePassword"

	hash, _ := HashPassword(password)
	result := VerifyPassword(hash, "")

	assert.False(t, result)
}

func TestVerifyPassword_InvalidHash(t *testing.T) {
	result := VerifyPassword("invalid_hash", "password")

	assert.False(t, result)
}

func TestVerifyPassword_EmptyHash(t *testing.T) {
	result := VerifyPassword("", "password")

	assert.False(t, result)
}

func TestVerifyPassword_MatchingPassword(t *testing.T) {
	passwords := []string{
		"simple",
		"Complex@Password123",
		"with spaces",
		"unicode_日本語",
	}

	for _, pwd := range passwords {
		hash, _ := HashPassword(pwd)
		assert.True(t, VerifyPassword(hash, pwd))
		assert.False(t, VerifyPassword(hash, pwd+"extra"))
	}
}
