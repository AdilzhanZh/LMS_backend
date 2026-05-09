package auth

import "github.com/AdilzhanZh/LMS_backend/internal/models"

type TokenManager interface {
	NewAccessToken(user models.User) (string, int64, error)
}
