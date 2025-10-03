package services

import (
	"errors"

	"github.com/joefazee/learning-go-shop/internal/dto"
	"github.com/joefazee/learning-go-shop/internal/models"
	"gorm.io/gorm"
)

var _ CartServiceInterface = (*CartService)(nil)

type CartService struct {
	db *gorm.DB
}

func NewCartService(db *gorm.DB) *CartService {
	return &CartService{db: db}
}

func (s *CartService) GetCart(userID uint) (*dto.CartResponse, error) {
	var cart models.Cart
	err := s.db.Preload("CartItems.Product.Category").
		Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}

	return s.convertToCartResponse(&cart), nil
}

func (s *CartService) AddToCart(userID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error) {

	// Check if product exists
	var product models.Product
	if err := s.db.First(&product, req.ProductID).Error; err != nil {
		return nil, errors.New("product not found")
	}

	if product.Stock < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	// Get or create cart
	var cart models.Cart
	if err := s.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		cart = models.Cart{UserID: userID}
		if err := s.db.Create(&cart).Error; err != nil {
			return nil, err
		}
	}

	// Check if item already exists in cart
	var cartItem models.CartItem
	if err := s.db.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&cartItem).Error; err != nil {
		// Create new cart item
		cartItem = models.CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		s.db.Create(&cartItem)
	} else {
		// Update existing cart item
		cartItem.Quantity += req.Quantity
		if cartItem.Quantity > product.Stock {
			return nil, errors.New("insufficient stock")
		}
		s.db.Save(&cartItem)
	}

	return s.GetCart(userID)
}

func (s *CartService) UpdateCartItem(userID, itemID uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error) {
	var cartItem models.CartItem
	if err := s.db.Joins("JOIN carts ON cart_items.cart_id = carts.id").
		Where("cart_items.id = ? AND carts.user_id = ?", itemID, userID).
		First(&cartItem).Error; err != nil {
		return nil, errors.New("cart item not found")
	}

	var product models.Product
	if err := s.db.First(&product, cartItem.ProductID).Error; err != nil {
		return nil, errors.New("product not found")
	}

	if product.Stock < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	cartItem.Quantity = req.Quantity
	if err := s.db.Save(&cartItem).Error; err != nil {
		return nil, err
	}

	return s.GetCart(userID)
}

func (s *CartService) RemoveFromCart(userID, itemID uint) error {
	return s.db.Where("id = ? AND cart_id IN (?)", itemID,
		s.db.Select("id").Table("carts").
			Where("user_id = ?", userID)).
		Delete(&models.CartItem{}).Error
}

func (s *CartService) convertToCartResponse(cart *models.Cart) *dto.CartResponse {

	cartItems := make([]dto.CartItemResponse, len(cart.CartItems)) // memory allocation
	var total float64

	for i := range cart.CartItems {
		subtotal := float64(cart.CartItems[i].Quantity) * cart.CartItems[i].Product.Price
		total += subtotal

		cartItems[i] = dto.CartItemResponse{
			ID: cart.CartItems[i].ID,
			Product: dto.ProductResponse{
				ID:          cart.CartItems[i].Product.ID,
				CategoryID:  cart.CartItems[i].Product.CategoryID,
				Name:        cart.CartItems[i].Product.Name,
				Description: cart.CartItems[i].Product.Description,
				Price:       cart.CartItems[i].Product.Price,
				Stock:       cart.CartItems[i].Product.Stock,
				SKU:         cart.CartItems[i].Product.SKU,
				IsActive:    cart.CartItems[i].Product.IsActive,
				Category: dto.CategoryResponse{
					ID:          cart.CartItems[i].Product.Category.ID,
					Name:        cart.CartItems[i].Product.Category.Name,
					Description: cart.CartItems[i].Product.Category.Description,
					IsActive:    cart.CartItems[i].Product.Category.IsActive,
				},
			},
			Quantity:  cart.CartItems[i].Quantity,
			Subtotal:  subtotal,
			CreatedAt: cart.CartItems[i].CreatedAt,
			UpdatedAt: cart.CartItems[i].UpdatedAt,
		}
	}

	return &dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CartItems: cartItems,
		Total:     total,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}
}
