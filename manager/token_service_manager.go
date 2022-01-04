package manager

import (
	"go-gin-jwt/model"
	"go-gin-jwt/service"
)

type TokenServiceManagerEntity interface {
	GetTokenService() service.TokenServiceEntity
}

type TokenServiceManager struct {
	tokenService service.TokenServiceEntity
}

func NewTokenServiceManager(tokenConfig model.TokenConfig) TokenServiceManagerEntity {
	tokenService := service.NewTokenService(tokenConfig)
	return &TokenServiceManager{tokenService: tokenService}
}

func (t *TokenServiceManager) GetTokenService() service.TokenServiceEntity {
	return t.tokenService
}
