package repository

import (
	"errors"
	"go-gin-jwt/model"
)

type AuthenticationRepositoryEntity interface {
	AuthenticateUser(credential model.Credential) (*model.Credential, error)
}

type AuthenticationRepository struct{}

func NewAuthenticationRepository() AuthenticationRepositoryEntity {
	return &AuthenticationRepository{}
}

func (a *AuthenticationRepository) AuthenticateUser(credential model.Credential) (*model.Credential, error) {
	if credential.Username == "user" && credential.Password == "password" {
		return &credential, nil
	}

	return nil, errors.New("wrong credential")
}
