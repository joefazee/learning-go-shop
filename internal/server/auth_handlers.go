package server

import (
	"github.com/gin-gonic/gin"
	"github.com/joefazee/learning-go-shop/internal/dto"
	"github.com/joefazee/learning-go-shop/internal/services"
	"github.com/joefazee/learning-go-shop/internal/utils"
)

func (s *Server) register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	authService := services.NewAuthService(s.db, s.config)
	response, err := authService.Register(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Registration failed", err)
		return
	}

	utils.CreatedResponse(c, "User registered successfully", response)
}

func (s *Server) login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	authService := services.NewAuthService(s.db, s.config)
	response, err := authService.Login(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "Login failed")
		return
	}

	utils.SuccessResponse(c, "Login successful", response)
}

func (s *Server) refreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	authService := services.NewAuthService(s.db, s.config)
	response, err := authService.RefreshToken(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "Token refresh failed")
		return
	}

	utils.SuccessResponse(c, "Token refreshed successfully", response)
}

func (s *Server) logout(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	authService := services.NewAuthService(s.db, s.config)
	if err := authService.Logout(req.RefreshToken); err != nil {
		utils.InternalServerErrorResponse(c, "Logout failed", err)
		return
	}

	utils.SuccessResponse(c, "Logout successful", nil)
}

func (s *Server) getProfile(c *gin.Context) {

	userID := c.GetUint("user_id")
	userService := services.NewUserService(s.db)
	profile, err := userService.GetProfile(userID)
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	utils.SuccessResponse(c, "Profile retrieved successfully", profile)
}

func (s *Server) updateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	userService := services.NewUserService(s.db)
	profile, err := userService.UpdateProfile(userID, &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update profile", err)
		return
	}
	utils.SuccessResponse(c, "Profile updated successfully", profile)
}
