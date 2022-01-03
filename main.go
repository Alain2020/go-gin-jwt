package main

import (
	"go-gin-jwt/authenticator"
	"go-gin-jwt/delivery/middleware"
	"go-gin-jwt/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

func main() {
	r := gin.Default()
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	tokenConfig := authenticator.TokenConfig{
		ApplicationName:     "ENIGMA",
		JwtSigningMethod:    jwt.SigningMethodHS256,
		JwtSignatureKey:     "P@ssw0rd",
		AccessTokenLifeTime: 3600 * time.Second,
		Client:              client,
	}
	tokenService := authenticator.NewTokenService(tokenConfig)
	r.Use(middleware.NewTokenValidator(tokenService).RequireToken())
	publicRoute := r.Group("/enigma")
	publicRoute.POST("/auth", func(c *gin.Context) {
		var user model.Credential
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "can't bind struct",
			})
			return
		}

		if user.Username == "enigma" && user.Password == "123" {
			token, err := tokenService.CreateAccessToken(&user)
			if err != nil {
				log.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			err = tokenService.StoreAccessToken(user.Username, token)
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	})

	publicRoute.GET("/user", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user": c.GetString("username"),
		})
	})
	publicRoute.POST("/logout", func(c *gin.Context) {
		if err := tokenService.DeleteAccessToken(c.GetString("uuid")); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"message": gin.H{"message": "success to logout"}})
			return
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
