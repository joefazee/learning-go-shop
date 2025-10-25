package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/joefazee/learning-go-shop/docs"
	"github.com/joefazee/learning-go-shop/internal/config"
	"github.com/joefazee/learning-go-shop/internal/services"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config         *config.Config
	logger         *zerolog.Logger
	authService    services.AuthServiceInterface
	productService services.ProductServiceInterface
	userService    services.UserServiceInterface
	uploadService  services.UploadServiceInterface
	cartService    services.CartServiceInterface
	orderService   services.OrderServiceInterface
}

func New(cfg *config.Config,
	logger *zerolog.Logger,
	authService services.AuthServiceInterface,
	productService services.ProductServiceInterface,
	userService services.UserServiceInterface,
	uploadService services.UploadServiceInterface,
	cartService services.CartServiceInterface,
	orderService services.OrderServiceInterface,
) *Server {
	return &Server{
		config:         cfg,
		logger:         logger,
		authService:    authService,
		productService: productService,
		userService:    userService,
		uploadService:  uploadService,
		cartService:    cartService,
		orderService:   orderService,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.New()

	// Add middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())

	// Add routes
	router.GET("/health", s.healthCheck)

	// Add documentation routes
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.StaticFile("/api-docs", "./docs/rapidoc.html")

	router.Static("/uploads", "./uploads")

	router.GET("/playground", s.playgroundHandler())
	router.GET("/playground/public", s.playgroundPublicHandler())
	router.GET("/playground/protected", s.playgroundProtectedHandler())

	graphqlPublic := router.Group("/graphql/public")
	graphqlPublic.Use(s.graphqlMiddleware())
	graphqlPublic.POST("/", s.graphqlHandler())

	graphqlProtected := router.Group("/graphql")
	graphqlProtected.Use(s.authMiddleware())
	graphqlProtected.Use(s.graphqlMiddleware())
	graphqlProtected.POST("/", s.graphqlHandler())

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{ //nolint:gocritic // I need this for readability
			auth.POST("/register", s.register)
			auth.POST("/login", s.login)
			auth.POST("/refresh", s.refreshToken)
			auth.POST("/logout", s.logout)

		}

		protected := api.Group("/")
		protected.Use(s.authMiddleware())
		{
			// User routes
			users := protected.Group("/users")
			{
				userRoutes := users
				userRoutes.GET("/profile", s.getProfile)
				userRoutes.PUT("/profile", s.updateProfile)
			}

			// category routes
			categories := protected.Group("/categories")
			{
				categoryRoute := categories
				categoryRoute.POST("/", s.adminMiddleware(), s.createCategory)
				categoryRoute.PUT("/:id", s.adminMiddleware(), s.updateCategory)
				categoryRoute.DELETE("/:id", s.adminMiddleware(), s.deleteCategory)
			}

			// product routes
			products := protected.Group("/products")
			{
				productRoutes := products
				productRoutes.POST("/", s.adminMiddleware(), s.createProduct)
				productRoutes.PUT("/:id", s.adminMiddleware(), s.updateProduct)
				productRoutes.DELETE("/:id", s.adminMiddleware(), s.deleteProduct)
				productRoutes.POST("/:id/images", s.adminMiddleware(), s.uploadProductImage)

			}

			// cart routes
			cart := protected.Group("/cart")
			{
				cartRoutes := cart
				cartRoutes.GET("/", s.getCart)
				cartRoutes.POST("/items", s.addToCart)
				cartRoutes.PUT("/items/:id", s.updateCartItem)
				cartRoutes.DELETE("/items/:id", s.removeFromCart)
			}

			// Order routes
			orders := protected.Group("/orders")
			{
				orderRoutes := orders
				orderRoutes.POST("/", s.createOrder)
				orderRoutes.GET("/", s.getOrders)
				orderRoutes.GET("/:id", s.getOrder)
			}
		}

		// public routes
		api.GET("/categories", s.getCategories)
		api.GET("/search", s.searchProducts)
		api.GET("/products", s.getProducts)
		api.GET("/products/:id", s.getProduct)

	}

	return router
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) corsMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}

}
