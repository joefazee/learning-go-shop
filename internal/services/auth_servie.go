package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/joefazee/learning-go-shop/internal/config"
	"github.com/joefazee/learning-go-shop/internal/dto"
	"github.com/joefazee/learning-go-shop/internal/events"
	"github.com/joefazee/learning-go-shop/internal/models"
	"github.com/joefazee/learning-go-shop/internal/notifications"
	"github.com/joefazee/learning-go-shop/internal/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	db             *gorm.DB
	config         *config.Config
	eventPublisher events.Publisher
}

func NewAuthService(db *gorm.DB, config *config.Config, eventPublisher events.Publisher) *AuthService {
	return &AuthService{
		db:             db,
		config:         config,
		eventPublisher: eventPublisher,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if user exists
	var existingUser models.User
	if err := s.db.Where("email  = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("you cannot register with this email")
	}
	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      models.UserRoleCustomer,
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}
	// create a cart
	cart := models.Cart{UserID: user.ID}
	if err := s.db.Create(&cart).Error; err != nil {
		fmt.Println("Unable to create cart")
	}

	// generate token
	return s.generateAuthResponse(&user)

}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	var user models.User
	if err := s.db.Where("email = ? AND is_active = ?", req.Email, true).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return s.generateAuthResponse(&user)
}

func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claims, err := utils.ValidateToken(req.RefreshToken, s.config.JWT.Secret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	var refreshToken models.RefreshToken
	if err := s.db.Where("token = ? AND expires_at > ?", req.RefreshToken, time.Now()).First(&refreshToken).Error; err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	var user models.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	s.db.Delete(&refreshToken)

	return s.generateAuthResponse(&user)
}

func (s *AuthService) Logout(refreshToken string) error {
	return s.db.Where("token = ?", refreshToken).Delete(&models.RefreshToken{}).Error
}

func (s *AuthService) generateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	accessToken, refreshToken, err := utils.GenerateTokenPair(
		&s.config.JWT,
		user.ID,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.config.JWT.RefreshTokenExpires),
	}

	s.db.Create(&refreshTokenModel)

	err = s.eventPublisher.Publish(notifications.UserLoggedIn, user, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("unable to publish user login event: %w", err)
	}

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Role:      string(user.Role),
			IsActive:  user.IsActive,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
