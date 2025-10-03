package resolver

import (
	"strconv"

	"github.com/joefazee/learning-go-shop/internal/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	authService    services.AuthServiceInterface
	userService    services.UserServiceInterface
	productService services.ProductServiceInterface
	cartService    services.CartServiceInterface
	orderService   services.OrderServiceInterface
}

func NewResolver(authService services.AuthServiceInterface,
	userService services.UserServiceInterface,
	productService services.ProductServiceInterface,
	cartService services.CartServiceInterface,
	orderService services.OrderServiceInterface) *Resolver {

	return &Resolver{
		authService:    authService,
		userService:    userService,
		productService: productService,
		cartService:    cartService,
		orderService:   orderService,
	}

}

func (r *Resolver) parseID(id string) (uint, error) {
	parsed, err := strconv.ParseUint(id, 10, 32)
	return uint(parsed), err
}
