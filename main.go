package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var ApplicationName = "ENIGMA"
var JwtSigningMethod = jwt.SigningMethodHS256
var JwtSignatureKey = []byte("P@ssw0rd")

type MyClaims struct {
	jwt.StandardClaims
	Username string `json:"Username"`
	Email    string `json:"Email"`
}

type authHeader struct {
	AuthorizationHeader string `header:"Authorization"`
}

func main() {
	r := gin.Default()
	r.Use(AuthtokenMidddleware())
	publicRoute := r.Group("/enigma")
	publicRoute.POST("/auth", func(c *gin.Context) {
		var user Credential
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "can't bind struct",
			})
			return
		}

		if user.Username == "enigma" && user.Password == "123" {
			token, err := GenerateToken(user.Username, "user@gmail.com")
			if err != nil {
				log.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
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
			"user": "user",
		})
	})

	r.GET("/customer", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "customer"})
	})
	err := r.Run("localhost:3000")
	if err != nil {
		panic(err)
	}

}

func parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		} else if method != JwtSigningMethod {
			return nil, fmt.Errorf("signing method invalid")
		}
		return JwtSignatureKey, nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func AuthtokenMidddleware() gin.HandlerFunc {
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
			log.Println(tokenString)
			if tokenString == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
				c.Abort()
			}
			token, err := parseToken(tokenString)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Internal server error",
				})
				c.Abort()
				return
			}

			if token["iss"] == ApplicationName {
				c.Next()
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
				c.Abort()
				return
			}
		}
	}
}

func GenerateToken(userName string, email string) (string, error) {
	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:   ApplicationName,
			IssuedAt: time.Now().Unix(),
		},
		Username: userName,
		Email:    email,
	}

	token := jwt.NewWithClaims(JwtSigningMethod, claims)
	return token.SignedString(JwtSignatureKey)
}
