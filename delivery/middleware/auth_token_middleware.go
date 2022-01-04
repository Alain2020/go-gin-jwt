package middleware

import (
	"go-gin-jwt/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type authHeader struct {
	AuthorizationHeader string `header:"Authorization"`
}

type AuthTokenMiddleware struct {
	tokenService service.TokenServiceEntity
}

func NewTokenValidator(tokenService service.TokenServiceEntity) *AuthTokenMiddleware {
	return &AuthTokenMiddleware{
		tokenService: tokenService,
	}
}

func (a *AuthTokenMiddleware) RequireToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/auth/login" {
			c.Next()
		} else {
			h := authHeader{}
			if err := c.ShouldBindHeader(&h); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
				c.Abort()
			}

			tokenString := strings.Replace(h.AuthorizationHeader, "Bearer ", "", -1)
			if tokenString == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
				c.Abort()
			}

			token, err := a.tokenService.VerifyAccessToken(tokenString)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "unauthorized",
				})
				c.Abort()
				return
			}

			userName, err := a.tokenService.FetchAccessToken(token)
			if userName == "" || err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "unauthorized",
				})
				c.Abort()
				return
			}

			if token != nil {
				c.Set("username", userName)
				c.Set("uuid", token.AccessUuid)
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
				return
			}
		}
	}

}
