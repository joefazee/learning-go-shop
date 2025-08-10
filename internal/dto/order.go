package dto

type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

type CartResponse struct {
	ID        uint               `json:"id"`
	UserID    uint               `json:"user_id"`
	CartItems []CartItemResponse `json:"cart_items"`
	Total     float64            `json:"total"`
}

type CartItemResponse struct {
	ID       uint            `json:"id"`
	Product  ProductResponse `json:"product"`
	Quantity int             `json:"quantity"`
	Subtotal float64         `json:"subtotal"`
}

type OrderResponse struct {
	ID          uint                `json:"id"`
	UserID      uint                `json:"user_id"`
	Status      string              `json:"status"`
	TotalAmount float64             `json:"total_amount"`
	OrderItems  []OrderItemResponse `json:"order_items"`
	CreatedAt   string              `json:"created_at"`
}

type OrderItemResponse struct {
	ID       uint            `json:"id"`
	Product  ProductResponse `json:"product"`
	Quantity int             `json:"quantity"`
	Price    float64         `json:"price"`
}
