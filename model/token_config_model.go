package model

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

type TokenConfig struct {
	ApplicationName     string
	JwtSignatureKey     string
	JwtSigningMethod    *jwt.SigningMethodHMAC
	AccessTokenLifeTime time.Duration
	Client              *redis.Client
}
