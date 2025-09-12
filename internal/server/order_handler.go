package server

import (
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/joefazee/learning-go-shop/internal/dto"
	"github.com/joefazee/learning-go-shop/internal/utils"
)

// @Summary Create an order
// @Description Create an order from the current user's cart
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Success 201 {object} utils.Response{data=dto.OrderResponse} "Order created successfully"
// @Failure 400 {object} utils.Response "Cart is empty or insufficient stock"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /orders [post]
func (s *Server) createOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	order, err := s.orderService.CreateOrder(userID)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to create order", err)
		return
	}

	utils.CreatedResponse(c, "Order created successfully", order)
}

// @Summary Get user's orders
// @Description Retrieve paginated list of user's orders
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse{data=[]dto.OrderResponse} "Orders retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /orders [get]
func (s *Server) getOrders(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, meta, err := s.orderService.GetOrders(userID, page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch orders", err)
		return
	}

	utils.PaginatedSuccessResponse(c, "Orders retrieved successfully", orders, *meta)
}

// @Summary Get order by ID
// @Description Retrieve detailed information about a specific order
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} utils.Response{data=dto.OrderResponse} "Order retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid order ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Order not found"
// @Router /orders/{id} [get]
func (s *Server) getOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid order ID", err)
		return
	}

	order, err := s.orderService.GetOrder(userID, uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Order not found")
		return
	}

	utils.SuccessResponse(c, "Order retrieved successfully", order)
}
