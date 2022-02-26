package services

import (
	"github.com/loupeznik/better-wapi/src/models"
)

type AuthService interface {
	ValidateCredentials(credentials models.Login) bool
}

type authService struct {
	config *models.Config
}

func NewAuthService(config *models.Config) *authService {
	return &authService{config: config}
}

func (s *authService) ValidateCredentials(credentials models.Login) bool {
	return credentials.Login == s.config.UserLogin && credentials.Secret == s.config.UserSecret
}
