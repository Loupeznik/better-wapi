package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/loupeznik/better-wapi/src/helpers"
	"github.com/loupeznik/better-wapi/src/models"
)

type AuthService struct {
	config *models.Config
}

type authCustomClaims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}

func NewAuthService(config *models.Config) *AuthService {
	return &AuthService{config: config}
}

func (s *AuthService) ValidateCredentials(credentials models.Login) bool {
	return credentials.Login == s.config.UserLogin && credentials.Secret == s.config.UserSecret
}

func (s *AuthService) IssueToken(credentials models.Login) (string, error) {
	if !s.ValidateCredentials(credentials) {
		return "", fmt.Errorf("invalid credentials")
	}

	claims := &authCustomClaims{
		credentials.Login,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		},
	}

	fmt.Println(s.config.JsonWebKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	result, err := token.SignedString([]byte(s.config.JsonWebKey))

	if err != nil {
		helpers.Log("error", fmt.Sprintf("Error while signing token: %s", err))
		return "", err
	}

	return result, err
}

func (s *AuthService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &authCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JsonWebKey), nil
	})
}
