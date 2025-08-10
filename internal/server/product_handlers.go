package server

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joefazee/learning-go-shop/internal/dto"
	"github.com/joefazee/learning-go-shop/internal/utils"
)

func (s *Server) createCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	category, err := s.productService.CreateCategory(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create category", err)
		return
	}

	utils.CreatedResponse(c, "Category created successfully", category)

}

func (s *Server) getCategories(c *gin.Context) {
	categories, err := s.productService.GetCategories()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch categories", err)
		return
	}

	utils.SuccessResponse(c, "Categories retrieved successfully", categories)
}

func (s *Server) updateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid category ID", err)
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	category, err := s.productService.UpdateCategory(uint(id), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update category", err)
		return
	}

	utils.SuccessResponse(c, "Category updated successfully", category)

}

func (s *Server) deleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid category ID", err)
		return
	}

	if err := s.productService.DeleteCategory(uint(id)); err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete category", err)
		return
	}

	utils.SuccessResponse(c, "Category deleted successfully", nil)
}

func (s *Server) createProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	product, err := s.productService.CreateProduct(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create product", err)
		return
	}

	utils.CreatedResponse(c, "Product created successfully", product)
}

func (s *Server) getProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, meta, err := s.productService.GetProducts(page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch products", err)
		return
	}

	utils.PaginatedSuccessResponse(c, "Products retrieved successfully", products, *meta)
}

func (s *Server) getProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid product ID", err)
		return
	}

	product, err := s.productService.GetProduct(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Product not found")
		return
	}

	utils.SuccessResponse(c, "Product retrieved successfully", product)
}

func (s *Server) updateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid product ID", err)
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	product, err := s.productService.UpdateProduct(uint(id), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update product", err)
		return
	}

	utils.SuccessResponse(c, "Product updated successfully", product)
}

func (s *Server) deleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid product ID", err)
		return
	}
	if err := s.productService.DeleteProduct(uint(id)); err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete product", err)
		return
	}

	utils.SuccessResponse(c, "Product deleted successfully", nil)
}
