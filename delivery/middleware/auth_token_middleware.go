package middleware

import (
	"go-gin-jwt/authenticator"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type authHeader struct {
	AuthorizationHeader string `header:"Authorization"`
}

type AuthTokenMiddleware struct {
	accessToken authenticator.Token
}

func NewTokenValidator(accessToken authenticator.Token) *AuthTokenMiddleware {
	return &AuthTokenMiddleware{
		accessToken: accessToken,
	}
}

func (a *AuthTokenMiddleware) RequireToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/enigma/auth" {
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

			token, err := a.accessToken.VerifyAccessToken(tokenString)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "unauthorized",
				})
				c.Abort()
				return
			}

			userName, err := a.accessToken.FetchAccessToken(token)
			if userName == "" || err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "unauthorized",
				})
				c.Abort()
				return
			}

			if token != nil {
				c.Set("username", userName)
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
				return
			}
		}
	}

}
