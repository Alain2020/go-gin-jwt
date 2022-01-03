package authenticator

import (
	"fmt"
	"go-gin-jwt/model"
	"time"

	"github.com/golang-jwt/jwt"
)

type Token interface {
	CreateAccessToken(credential *model.Credential) (string, error)
	VerifyAccessToken(tokenString string) (jwt.MapClaims, error)
}

type TokenConfig struct {
	ApplicationName     string
	JwtSignatureKey     string
	JwtSigningMethod    *jwt.SigningMethodHMAC
	AccessTokenLifeTime time.Duration
}
type token struct {
	Config TokenConfig
}

func NewTokenService(config TokenConfig) Token {
	return &token{
		Config: config,
	}
}

func (t *token) CreateAccessToken(credential *model.Credential) (string, error) {
	now := time.Now().UTC()
	end := now.Add(t.Config.AccessTokenLifeTime)
	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:   t.Config.ApplicationName,
			IssuedAt: time.Now().Unix(),
		},
		Username: credential.Username,
		Email:    credential.Email,
	}
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = end.Unix()

	token := jwt.NewWithClaims(t.Config.JwtSigningMethod, claims)
	return token.SignedString([]byte(t.Config.JwtSignatureKey))
}

func (t *token) VerifyAccessToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		} else if method != t.Config.JwtSigningMethod {
			return nil, fmt.Errorf("signing method invalid")
		}
		return []byte(t.Config.JwtSignatureKey), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}
	return claims, nil
}
