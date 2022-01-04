package app

import (
	"go-gin-jwt/delivery"
	"go-gin-jwt/manager"
)

type AppEntity interface {
	Run()
}
type App struct {
	InfrastructureRedisManager manager.InfrastructureRedisManagerEntity
	AppConfigManager           manager.AppConfigManagerManagerEntity
	TokenConfigManager         manager.TokenConfigManagerEntity
	Router                     *delivery.Router
}

func NewApp() AppEntity {
	appConfigManager := manager.NewAppConfigManager()
	infrastructureRedisManager := manager.NewInfrastructureRedisManager(appConfigManager)
	tokenConfigManager := manager.NewTokenConfigManager(infrastructureRedisManager, appConfigManager)
	authenticationRepositoryManager := manager.NewAuthenticationRepositoryManager()
	tokenServiceManager := manager.NewTokenServiceManager(tokenConfigManager.GetTokenConfig())
	authenticationUseCaseManager := manager.NewAuthenticationUseCaseManager(authenticationRepositoryManager.GetAuthenticationRepository(), tokenServiceManager.GetTokenService())

	app := new(App)
	app.AppConfigManager = appConfigManager
	app.InfrastructureRedisManager = infrastructureRedisManager
	app.TokenConfigManager = tokenConfigManager
	app.Router = delivery.NewRouter(authenticationUseCaseManager.GetAuthenticationUseCase(), tokenServiceManager.GetTokenService())
	return app
}

func (a *App) Run() {
	a.Router.RouterEngine.Run()
}
