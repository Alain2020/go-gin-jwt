package delivery

import (
	"go-gin-jwt/delivery/middleware"
	"go-gin-jwt/service"
	"go-gin-jwt/usecase"

	"github.com/gin-gonic/gin"
)

type RouteDeliveryHttpEntity interface {
	InitRoute(publicRouter *gin.RouterGroup)
}

type Router struct {
	routes       []RouteDeliveryHttpEntity
	RouterEngine *gin.Engine
	publicRoute  *gin.RouterGroup
}

func NewRouter(authenticationUseCase usecase.AuthenticationUseCaseEntity, tokenService service.TokenServiceEntity) *Router {
	router := new(Router)

	router.RouterEngine = gin.Default()
	router.routes = []RouteDeliveryHttpEntity{
		NewAuthenticationRoute(authenticationUseCase),
		NewUserRoute(),
	}

	tokenValidator := middleware.NewTokenValidator(tokenService)
	router.publicRoute = router.RouterEngine.Group("/api")
	router.publicRoute.Use(tokenValidator.RequireToken())
	router.initRoutes()
	return router
}

func (r *Router) initRoutes() {
	for _, route := range r.routes {
		route.InitRoute(r.publicRoute)
	}
}
