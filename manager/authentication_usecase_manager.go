package manager

import (
	"go-gin-jwt/repository"
	"go-gin-jwt/service"
	"go-gin-jwt/usecase"
)

type AuthenticationUseCaseManagerEntity interface {
	GetAuthenticationUseCase() usecase.AuthenticationUseCaseEntity
}

type AuthenticationUseCaseManager struct {
	authenticationUseCase usecase.AuthenticationUseCaseEntity
}

func NewAuthenticationUseCaseManager(
	authenticationUseCaseRepository repository.AuthenticationRepositoryEntity,
	tokenService service.TokenServiceEntity,
) AuthenticationUseCaseManagerEntity {
	authenticationUseCase := usecase.NewAuthenticationUseCase(authenticationUseCaseRepository, tokenService)
	return &AuthenticationUseCaseManager{
		authenticationUseCase: authenticationUseCase,
	}
}

func (a *AuthenticationUseCaseManager) GetAuthenticationUseCase() usecase.AuthenticationUseCaseEntity {
	return a.authenticationUseCase
}
