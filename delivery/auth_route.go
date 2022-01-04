package delivery

import (
	"go-gin-jwt/model"
	"go-gin-jwt/usecase"
	"go-gin-jwt/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthenticationRoute struct {
	authenticationUseCase usecase.AuthenticationUseCaseEntity
}

func NewAuthenticationRoute(authenticationUseCase usecase.AuthenticationUseCaseEntity) RouteDeliveryHttpEntity {
	return &AuthenticationRoute{authenticationUseCase: authenticationUseCase}
}

func (a *AuthenticationRoute) InitRoute(publicRouter *gin.RouterGroup) {
	authenticationRoute := publicRouter.Group("/auth")
	{
		authenticationRoute.POST("login", a.handleLogin)
		authenticationRoute.POST("logout", a.handleLogout)
	}
}

func (m *AuthenticationRoute) handleLogin(c *gin.Context) {
	var credential model.Credential
	err := c.ShouldBindJSON(&credential)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateHTTPRespond(http.StatusBadRequest, "bad request", err.Error()))
		return
	}

	tokenDetails, err := m.authenticationUseCase.Login(credential)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.CreateHTTPRespond(http.StatusUnauthorized, "unauthorized", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.CreateHTTPRespond(http.StatusOK, "success", tokenDetails))
}

func (m *AuthenticationRoute) handleLogout(c *gin.Context) {
	accessUuid := c.GetString("uuid")
	err := m.authenticationUseCase.Logout(accessUuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateHTTPRespond(http.StatusBadRequest, "bad request", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.CreateHTTPRespond(http.StatusOK, "success", nil))
}
