package authenticator

import (
	"context"
	"errors"
	"fmt"
	"go-gin-jwt/model"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Token interface {
	CreateAccessToken(credential *model.Credential) (*TokenDetails, error)
	VerifyAccessToken(tokenString string) (*AccessDetails, error)
	StoreAccessToken(userName string, tokenDetails *TokenDetails) error
	FetchAccessToken(accessDetails *AccessDetails) (string, error)
}

type TokenConfig struct {
	ApplicationName     string
	JwtSignatureKey     string
	JwtSigningMethod    *jwt.SigningMethodHMAC
	AccessTokenLifeTime time.Duration
	Client              *redis.Client
}

type TokenDetails struct {
	AccessToken string
	AccessUuid  string
	AtExpires   int64
}

type AccessDetails struct {
	AccessUuid string
	UserName   string
}
type token struct {
	Config TokenConfig
}

func NewTokenService(config TokenConfig) Token {
	return &token{
		Config: config,
	}
}

func (t *token) CreateAccessToken(credential *model.Credential) (*TokenDetails, error) {
	td := &TokenDetails{}
	now := time.Now().UTC()
	end := now.Add(t.Config.AccessTokenLifeTime)

	td.AtExpires = end.Unix()
	td.AccessUuid = uuid.New().String()

	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:   t.Config.ApplicationName,
			IssuedAt: time.Now().Unix(),
		},
		Username:   credential.Username,
		Email:      credential.Email,
		AccessUUID: td.AccessUuid,
	}
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = end.Unix()

	token := jwt.NewWithClaims(t.Config.JwtSigningMethod, claims)
	newToken, err := token.SignedString([]byte(t.Config.JwtSignatureKey))
	td.AccessToken = newToken
	return td, err
}

func (t *token) VerifyAccessToken(tokenString string) (*AccessDetails, error) {
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
	accessUUID := claims["AccessUUID"].(string)
	userName := claims["Username"].(string)
	return &AccessDetails{
		AccessUuid: accessUUID,
		UserName:   userName,
	}, nil

}

func (t *token) StoreAccessToken(userName string, tokenDetails *TokenDetails) error {
	at := time.Unix(tokenDetails.AtExpires, 0)
	now := time.Now()
	err := t.Config.Client.Set(context.Background(), tokenDetails.AccessUuid, userName, at.Sub(now)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (t *token) FetchAccessToken(accessDetails *AccessDetails) (string, error) {
	if accessDetails != nil {
		userName, err := t.Config.Client.Get(context.Background(), accessDetails.AccessUuid).Result()
		if err != nil {
			return "", err
		}
		return userName, nil
	} else {
		return "", errors.New("invalid Access")
	}
}
