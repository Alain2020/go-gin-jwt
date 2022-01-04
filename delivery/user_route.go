package delivery

import (
	"go-gin-jwt/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
}

func NewUserRoute() RouteDeliveryHttpEntity {
	return &UserRoute{}
}

func (u *UserRoute) InitRoute(publicRouter *gin.RouterGroup) {
	{
		userRoute := publicRouter.Group("/user")
		{
			userRoute.GET("", u.handleGetUser)
		}
	}
}

func (u *UserRoute) handleGetUser(c *gin.Context) {
	username := c.GetString("username")
	c.JSON(http.StatusOK, utils.CreateHTTPRespond(http.StatusOK, "success", username))
}
