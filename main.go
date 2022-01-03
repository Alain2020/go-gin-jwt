package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authHeader struct {
	AuthorizationHeader string `header:"Authorization"`
}

func main() {
	r := gin.Default()
	r.Use(AuthtokenMidddleware())
	r.POST("/login", func(c *gin.Context) {
		var user Credential
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "can't bind struct",
			})
			return
		}

		if user.Username == "enigma" && user.Password == "123" {
			c.JSON(http.StatusOK, gin.H{
				"token": "123",
			})
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	})

	r.GET("/customer", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "customer"})
	})
	err := r.Run("localhost:3000")
	if err != nil {
		panic(err)
	}

}

func AuthtokenMidddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/login" {
			c.Next()
		} else {
			h := authHeader{}
			if err := c.ShouldBindHeader(&h); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
				c.Abort()
			}

			if h.AuthorizationHeader == "1234" {
				c.Next()
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
				c.Abort()

			}
		}
	}
}
