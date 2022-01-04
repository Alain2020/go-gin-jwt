package usecase

import (
	"go-gin-jwt/model"
	"go-gin-jwt/repository"
	"go-gin-jwt/service"
)

type AuthenticationUseCaseEntity interface {
	Login(credential model.Credential) (*model.TokenDetails, error)
	Logout(accessUuid string) error
}

type AuthenticationUseCase struct {
	authenticationRepository repository.AuthenticationRepositoryEntity
	tokenService             service.TokenServiceEntity
}

func NewAuthenticationUseCase(
	authenticationRepository repository.AuthenticationRepositoryEntity,
	tokenService service.TokenServiceEntity,
) AuthenticationUseCaseEntity {
	return &AuthenticationUseCase{authenticationRepository: authenticationRepository, tokenService: tokenService}
}

func (a *AuthenticationUseCase) Login(credential model.Credential) (*model.TokenDetails, error) {
	_, err := a.authenticationRepository.AuthenticateUser(credential)
	if err != nil {
		return nil, err
	}

	tokenDetails, err := a.tokenService.CreateAccessToken(&credential)
	if err != nil {
		return nil, err
	}
	err = a.tokenService.StoreAccessToken(credential.Username, tokenDetails)
	if err != nil {
		return nil, err
	}
	return tokenDetails, err
}

func (a *AuthenticationUseCase) Logout(accessUuid string) error {
	return a.tokenService.DeleteAccessToken(accessUuid)
}
