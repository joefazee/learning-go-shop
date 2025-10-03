package services

import (
	"errors"
	"fmt"

	"github.com/joefazee/learning-go-shop/internal/dto"
	"github.com/joefazee/learning-go-shop/internal/models"
	"github.com/joefazee/learning-go-shop/internal/utils"
	"gorm.io/gorm"
)

var _ OrderServiceInterface = (*OrderService)(nil)

type OrderService struct {
	db *gorm.DB
}

// NewOrderService creates the order service type
func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

func (s *OrderService) CreateOrder(userID uint) (*dto.OrderResponse, error) {
	var orderResponse *dto.OrderResponse

	err := s.db.Transaction(func(tx *gorm.DB) error {

		var cart models.Cart
		if err := tx.Preload("CartItems.Product").Where("user_id = ?", userID).First(&cart).Error; err != nil {
			return errors.New("cart not found")
		}

		if len(cart.CartItems) == 0 {
			return errors.New("cart is empty")
		}

		// Calculate total and validate stock
		var totalAmount float64
		var orderItems []models.OrderItem

		for i := range cart.CartItems {
			cartItem := &cart.CartItems[i]

			if cartItem.Product.Stock < cartItem.Quantity {
				return fmt.Errorf("insufficient stock for product: %s", cartItem.Product.Name)
			}

			itemTotal := float64(cartItem.Quantity) * cartItem.Product.Price
			totalAmount += itemTotal

			orderItems = append(orderItems, models.OrderItem{
				ProductID: cartItem.ProductID,
				Quantity:  cartItem.Quantity,
				Price:     cartItem.Product.Price,
			})

			// Update product stock
			cartItem.Product.Stock -= cartItem.Quantity
			if err := tx.Save(&cartItem.Product).Error; err != nil {
				return err
			}

			// Create order
			order := models.Order{
				UserID:      userID,
				Status:      models.OrderStatusPending,
				TotalAmount: totalAmount,
				OrderItems:  orderItems,
			}

			if err := tx.Create(&order).Error; err != nil {
				return err
			}

			// Clear cart
			if err := tx.Unscoped().Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
				return err
			}

			response, err := s.getOrderResponse(tx, order.ID)
			if err != nil {
				return err
			}

			orderResponse = response
		}

		return nil // Transaction successful
	})

	if err != nil {
		return nil, err
	}

	return orderResponse, nil

}

func (s *OrderService) GetOrders(userID uint, page, limit int) ([]dto.OrderResponse, *utils.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	var orders []models.Order
	var total int64

	s.db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&total)

	if err := s.db.Preload("OrderItems.Product.Category").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, nil, err
	}

	response := make([]dto.OrderResponse, len(orders))
	for i := range orders {
		order := &orders[i]
		response[i] = s.convertToOrderResponse(order)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	meta := &utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return response, meta, nil
}

func (s *OrderService) GetOrder(userID, orderID uint) (*dto.OrderResponse, error) {
	var order models.Order
	if err := s.db.Preload("OrderItems.Product.Category").
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error; err != nil {
		return nil, err
	}

	response := s.convertToOrderResponse(&order)

	return &response, nil
}

func (s *OrderService) getOrderResponse(tx *gorm.DB, orderID uint) (*dto.OrderResponse, error) {
	var order models.Order
	if err := tx.Preload("OrderItems.Product.Category").First(&order, orderID).Error; err != nil {
		return nil, err
	}

	response := s.convertToOrderResponse(&order)

	return &response, nil
}

func (s *OrderService) convertToOrderResponse(order *models.Order) dto.OrderResponse {
	orderItems := make([]dto.OrderItemResponse, len(order.OrderItems))
	for i := range order.OrderItems {
		item := order.OrderItems[i]

		orderItems[i] = dto.OrderItemResponse{
			ID: item.ID,
			Product: dto.ProductResponse{
				ID:          item.Product.ID,
				CategoryID:  item.Product.CategoryID,
				Name:        item.Product.Name,
				Description: item.Product.Description,
				Price:       item.Product.Price,
				Stock:       item.Product.Stock,
				SKU:         item.Product.SKU,
				IsActive:    item.Product.IsActive,
				Category: dto.CategoryResponse{
					ID:          item.Product.Category.ID,
					Name:        item.Product.Category.Name,
					Description: item.Product.Category.Description,
					IsActive:    item.Product.Category.IsActive,
				},
			},
			Quantity: item.Quantity,
			Price:    item.Price,

			CreatedAt: item.CreatedAt,
		}
	}

	return dto.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Status:      string(order.Status),
		TotalAmount: order.TotalAmount,
		OrderItems:  orderItems,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
}
