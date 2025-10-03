package repositories

import (
	"github.com/joefazee/learning-go-shop/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	return nil, nil
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	return nil, nil
}
func (r *UserRepository) GetByEmailAndActive(email string, isActive bool) (*models.User, error) {
	return nil, nil
}
func (r *UserRepository) Create(user *models.User) error {
	return nil
}
func (r *UserRepository) Update(user *models.User) error {
	return nil
}
func (r *UserRepository) Delete(id uint) error {
	return nil
}

func (r *UserRepository) CreateRefreshToken(token *models.RefreshToken) error {
	return nil
}
func (r *UserRepository) GetValidRefreshToken(token string) (*models.RefreshToken, error) {
	return nil, nil
}
func (r *UserRepository) DeleteRefreshToken(token string) error {
	return nil
}
func (r *UserRepository) DeleteRefreshTokenByID(id uint) error {
	return nil
}
