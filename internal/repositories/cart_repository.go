package repositories

import (
	"github.com/joefazee/learning-go-shop/internal/models"
	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{
		db: db,
	}
}

func (r *CartRepository) GetByUserID(userID uint) (*models.Cart, error) {
	return nil, nil
}
func (r *CartRepository) Create(cart *models.Cart) error {
	return nil
}
func (r *CartRepository) Update(cart *models.Cart) error {
	return nil
}
func (r *CartRepository) Delete(id uint) error {
	return nil
}
