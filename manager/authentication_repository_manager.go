package manager

import "go-gin-jwt/repository"

type AuthenticationRepositoryManagerEntity interface {
	GetAuthenticationRepository() repository.AuthenticationRepositoryEntity
}

type AuthenticationRepositoryManager struct {
	authenticationRepository repository.AuthenticationRepositoryEntity
}

func NewAuthenticationRepositoryManager() AuthenticationRepositoryManagerEntity {
	authenticationRepository := repository.AuthenticationRepository{}
	return &AuthenticationRepositoryManager{authenticationRepository: &authenticationRepository}
}

func (a *AuthenticationRepositoryManager) GetAuthenticationRepository() repository.AuthenticationRepositoryEntity {
	return a.authenticationRepository
}
