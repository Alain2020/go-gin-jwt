package service

import (
	"context"
	"errors"
	"fmt"
	"go-gin-jwt/model"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type TokenServiceEntity interface {
	CreateAccessToken(credential *model.Credential) (*model.TokenDetails, error)
	VerifyAccessToken(tokenString string) (*model.UserCredential, error)
	StoreAccessToken(userName string, tokenDetails *model.TokenDetails) error
	FetchAccessToken(userCredential *model.UserCredential) (string, error)
	DeleteAccessToken(accessUuid string) error
}

type TokenService struct {
	tokenConfig model.TokenConfig
}

func NewTokenService(config model.TokenConfig) TokenServiceEntity {
	return &TokenService{
		tokenConfig: config,
	}
}
func (t *TokenService) CreateAccessToken(credential *model.Credential) (*model.TokenDetails, error) {
	td := &model.TokenDetails{}
	now := time.Now().UTC()
	end := now.Add(t.tokenConfig.AccessTokenLifeTime)

	td.AtExpires = end.Unix()
	td.AccessUuid = uuid.New().String()

	claims := model.JwtClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:   t.tokenConfig.ApplicationName,
			IssuedAt: time.Now().Unix(),
		},
		Username:   credential.Username,
		Email:      credential.Email,
		AccessUUID: td.AccessUuid,
	}
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = end.Unix()

	token := jwt.NewWithClaims(t.tokenConfig.JwtSigningMethod, claims)
	newToken, err := token.SignedString([]byte(t.tokenConfig.JwtSignatureKey))
	td.AccessToken = newToken
	return td, err
}

func (t *TokenService) VerifyAccessToken(tokenString string) (*model.UserCredential, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		} else if method != t.tokenConfig.JwtSigningMethod {
			return nil, fmt.Errorf("signing method invalid")
		}
		return []byte(t.tokenConfig.JwtSignatureKey), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}
	accessUUID := claims["AccessUUID"].(string)
	userName := claims["Username"].(string)
	return &model.UserCredential{
		AccessUuid: accessUUID,
		UserName:   userName,
	}, nil

}

func (t *TokenService) StoreAccessToken(userName string, tokenDetails *model.TokenDetails) error {
	at := time.Unix(tokenDetails.AtExpires, 0)
	now := time.Now()
	err := t.tokenConfig.Client.Set(context.Background(), tokenDetails.AccessUuid, userName, at.Sub(now)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (t *TokenService) FetchAccessToken(userCredential *model.UserCredential) (string, error) {
	if userCredential != nil {
		userName, err := t.tokenConfig.Client.Get(context.Background(), userCredential.AccessUuid).Result()
		if err != nil {
			return "", err
		}
		return userName, nil
	} else {
		return "", errors.New("invalid Access")
	}
}

func (t *TokenService) DeleteAccessToken(accessUuid string) error {
	if accessUuid != "" {
		rd := t.tokenConfig.Client.Del(context.Background(), accessUuid)
		return rd.Err()
	} else {
		return errors.New("token not found")
	}
}
