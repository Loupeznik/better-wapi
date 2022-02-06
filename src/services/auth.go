package services

import "github.com/loupeznik/better-wapi/src/models"

type AuthService interface {
	Login(credentials models.Login)
}
