package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/AdilzhanZh/LMS_backend/internal/auth"
	"github.com/AdilzhanZh/LMS_backend/internal/models"
	"github.com/AdilzhanZh/LMS_backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo         repository.UserRepo
	tokenManager auth.TokenManager
}

func NewAuthService(userRepo repository.UserRepo, tokenManager auth.TokenManager) *AuthService {
	return &AuthService{
		repo:         userRepo,
		tokenManager: tokenManager,
	}
}

func (s *AuthService) Register(ctx context.Context, input models.RegisterUser) (int, error) {
	input.FullName = strings.TrimSpace(input.FullName)
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))

	if input.FullName == "" || input.Email == "" || input.Password == "" {
		return 0, errors.New("full_name, email, and password are required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user := models.CreateUser{
		FullName:     input.FullName,
		Email:        input.Email,
		PasswordHash: string(hash),
	}

	id, err := s.repo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, models.ErrUserAlreadyExists) {
			return 0, err
		}

		return 0, err
	}

	return id, nil
}

func (s *AuthService) Login(ctx context.Context, input models.LoginUser) (models.AuthTokens, error) {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	if input.Email == "" || input.Password == "" {
		return models.AuthTokens{}, errors.New("email and password are required")
	}

	user, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return models.AuthTokens{}, models.ErrInvalidCredentials
		}
		return models.AuthTokens{}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return models.AuthTokens{}, models.ErrInvalidCredentials
	}

	accessToken, accessExp, err := s.tokenManager.NewAccessToken(user)
	if err != nil {
		return models.AuthTokens{}, err
	}

	expiresIn := accessExp - time.Now().Unix()
	if expiresIn < 0 {
		expiresIn = 0
	}

	return models.AuthTokens{
		AccessToken: accessToken,
		//RefreshToken: refreshToken,
		ExpiresIn: expiresIn,
	}, nil
}
