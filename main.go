package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	type authHeader struct {
		AuthorizationHeader string `header:"Authorization"`
	}

	r.GET("/customer", func(c *gin.Context) {
		h := authHeader{}
		if err := c.ShouldBindHeader(&h); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}

		if h.AuthorizationHeader == "1234" {
			c.JSON(http.StatusOK, gin.H{
				"message": "customer",
			})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
	})
	err := r.Run("localhost:3000")
	if err != nil {
		panic(err)
	}

}
