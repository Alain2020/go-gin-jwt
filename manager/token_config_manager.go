package manager

import (
	"go-gin-jwt/model"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenConfigManagerEntity interface {
	GetTokenConfig() model.TokenConfig
}

type TokenConfigManager struct {
	model.TokenConfig
}

func NewTokenConfigManager(
	infrastructureRedisManager InfrastructureRedisManagerEntity,
	appConfigManger AppConfigManagerManagerEntity,
) TokenConfigManagerEntity {

	tokenConfig := model.TokenConfig{
		ApplicationName:     appConfigManger.GetAppConfig().ApplicationName,
		JwtSigningMethod:    jwt.SigningMethodHS256,
		JwtSignatureKey:     appConfigManger.GetAppConfig().JwtSignatureKey,
		AccessTokenLifeTime: 3600 * time.Second,
		Client:              infrastructureRedisManager.GetRedisClient(),
	}
	return &TokenConfigManager{TokenConfig: tokenConfig}
}

func (t *TokenConfigManager) GetTokenConfig() model.TokenConfig {
	return t.TokenConfig
}
