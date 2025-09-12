package server

import (
	"github.com/gin-gonic/gin"
	"github.com/joefazee/learning-go-shop/internal/dto"
	"github.com/joefazee/learning-go-shop/internal/utils"
)

// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration data"
// @Success 201 {object} utils.Response{data=dto.AuthResponse} "User registered successfully"
// @Failure 400 {object} utils.Response "Invalid request data or user already exists"
// @Router /auth/register [post]
func (s *Server) register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	response, err := s.authService.Register(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Registration failed", err)
		return
	}

	utils.CreatedResponse(c, "User registered successfully", response)
}

// @Summary User login
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User login credentials"
// @Success 200 {object} utils.Response{data=dto.AuthResponse} "Login successful"
// @Failure 401 {object} utils.Response "Invalid credentials"
// @Router /auth/login [post]
func (s *Server) login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	response, err := s.authService.Login(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "Login failed")
		return
	}

	utils.SuccessResponse(c, "Login successful", response)
}

// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} utils.Response{data=dto.AuthResponse} "Token refreshed successfully"
// @Failure 401 {object} utils.Response "Invalid refresh token"
// @Router /auth/refresh [post]
func (s *Server) refreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	response, err := s.authService.RefreshToken(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "Token refresh failed")
		return
	}

	utils.SuccessResponse(c, "Token refreshed successfully", response)
}

// @Summary User logout
// @Description Invalidate refresh token and logout user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token to invalidate"
// @Success 200 {object} utils.Response "Logout successful"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Router /auth/logout [post]
func (s *Server) logout(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	if err := s.authService.Logout(req.RefreshToken); err != nil {
		utils.InternalServerErrorResponse(c, "Logout failed", err)
		return
	}

	utils.SuccessResponse(c, "Logout successful", nil)
}

// @Summary Get user profile
// @Description Get current authenticated user's profile information
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.UserResponse} "Profile retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "User not found"
// @Router /users/profile [get]
func (s *Server) getProfile(c *gin.Context) {

	userID := c.GetUint("user_id")
	profile, err := s.userService.GetProfile(userID)
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	utils.SuccessResponse(c, "Profile retrieved successfully", profile)
}

// @Summary Update user profile
// @Description Update current authenticated user's profile information
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} utils.Response{data=dto.UserResponse} "Profile updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /users/profile [put]
func (s *Server) updateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	profile, err := s.userService.UpdateProfile(userID, &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update profile", err)
		return
	}
	utils.SuccessResponse(c, "Profile updated successfully", profile)
}
