package manager

import (
	"go-gin-jwt/model"
	"os"
)

type AppConfigManagerManagerEntity interface {
	GetAppConfig() *model.AppConfig
}

type AppConfigManager struct {
	appConfig model.AppConfig
}

func NewAppConfigManager() AppConfigManagerManagerEntity {
	applicationName := os.Getenv("APP_NAME")
	jwtSignatureKey := os.Getenv("JWT_SIGNATURE_KEY")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	appConfig := model.AppConfig{
		ApplicationName: applicationName,
		JwtSignatureKey: jwtSignatureKey,
		RedisHost:       redisHost,
		RedisPort:       redisPort,
	}
	return &AppConfigManager{appConfig: appConfig}
}

func (a *AppConfigManager) GetAppConfig() *model.AppConfig {
	return &a.appConfig
}
