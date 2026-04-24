package repositories

import (
	"github.com/coolmate/ecommerce-backend/internal/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	UpdatePassword(userID uint, passwordHash string) error
	Update(user *models.User) error
	GetRefreshToken(tokenHash string) (*models.RefreshToken, error)
	SaveRefreshToken(token *models.RefreshToken) error
	RevokeRefreshToken(tokenHash string) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Create(user *models.User) error {
	return ur.db.Create(user).Error
}

func (ur *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := ur.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) UpdatePassword(userID uint, passwordHash string) error {
	return ur.db.Model(&models.User{}).Where("id = ?", userID).Update("password_hash", passwordHash).Error
}

func (ur *UserRepository) Update(user *models.User) error {
	return ur.db.Save(user).Error
}

func (ur *UserRepository) GetRefreshToken(tokenHash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	if err := ur.db.Where("token_hash = ? AND revoked_at IS NULL", tokenHash).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (ur *UserRepository) SaveRefreshToken(token *models.RefreshToken) error {
	return ur.db.Create(token).Error
}

func (ur *UserRepository) RevokeRefreshToken(tokenHash string) error {
	return ur.db.Model(&models.RefreshToken{}).Where("token_hash = ?", tokenHash).Update("revoked_at", gorm.Expr("NOW()")).Error
}
