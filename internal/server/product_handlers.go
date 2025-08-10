package server

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joefazee/learning-go-shop/internal/dto"
	"github.com/joefazee/learning-go-shop/internal/services"
	"github.com/joefazee/learning-go-shop/internal/utils"
)

func (s *Server) createCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	productService := services.NewProductService(s.db)
	category, err := productService.CreateCategory(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create category", err)
		return
	}

	utils.CreatedResponse(c, "Category created successfully", category)

}

func (s *Server) getCategories(c *gin.Context) {
	productService := services.NewProductService(s.db)
	categories, err := productService.GetCategories()
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

	productService := services.NewProductService(s.db)
	category, err := productService.UpdateCategory(uint(id), &req)
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

	productService := services.NewProductService(s.db)
	if err := productService.DeleteCategory(uint(id)); err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete category", err)
		return
	}

	utils.SuccessResponse(c, "Category deleted successfully", nil)
}
