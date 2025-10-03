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
	var cart models.Cart
	if err := r.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}
func (r *CartRepository) Create(cart *models.Cart) error {
	return r.db.Create(cart).Error
}
func (r *CartRepository) Update(cart *models.Cart) error {
	return r.db.Save(cart).Error
}
func (r *CartRepository) Delete(id uint) error {
	return r.db.Delete(&models.Cart{}, id).Error
}
