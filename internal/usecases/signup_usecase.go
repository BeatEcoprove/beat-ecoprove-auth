package usecases

import (
	"github.com/BeatEcoprove/identityService/internal/domain"
	"github.com/BeatEcoprove/identityService/internal/repositories"
	"github.com/BeatEcoprove/identityService/pkg/contracts"
	fails "github.com/BeatEcoprove/identityService/pkg/errors"
	"github.com/BeatEcoprove/identityService/pkg/services"
)

type (
	// input
	SignUpInput struct {
		Email    string
		Password string
		Role     int
	}

	SignUpUseCase struct {
		authRepo     repositories.IAuthRepository
		profileRepo  repositories.IProfileRepository
		tokenService services.ITokenService
	}
)

func NewSignUpUseCase(
	authRepo repositories.IAuthRepository,
	profileRepo repositories.IProfileRepository,
	tokenService services.ITokenService,
) *SignUpUseCase {
	return &SignUpUseCase{
		authRepo:     authRepo,
		profileRepo:  profileRepo,
		tokenService: tokenService,
	}
}

func (as *SignUpUseCase) Handle(input SignUpInput) (*contracts.AuthResponse, error) {
	// Setps to create an user
	/// 1. Exists already this user registered?
	if ok := as.authRepo.ExistsUserWithEmail(input.Email); ok {
		return nil, fails.USER_ALREADY_EXISTS
	}

	if err := services.ValidatePassword(input.Password); err != nil {
		return nil, err
	}

	role, err := domain.GetRole(domain.Role(input.Role))

	if err != nil {
		return nil, fails.ROLE_NOT_FOUND
	}

	/// 2. If not then register the user
	identityUser := domain.NewIdentityUser(
		input.Email,
		input.Password,
		domain.Role(input.Role),
	)

	signUpTransaction, err := as.authRepo.BeginTransaction()

	if err != nil {
		return nil, fails.InternalServerError()
	}

	if err := signUpTransaction.Create(identityUser); err != nil {
		return nil, fails.InternalServerError()
	}

	/// 3. Generate an ProfileID -> Which will be the Main Profile
	//// Generate the profile info and the type
	profile := domain.NewProfile(identityUser.ID, domain.Main)

	if err := signUpTransaction.Create(profile); err != nil {
		return nil, fails.InternalServerError()
	}

	/// 4. Generate the authorization tokens -> Generate Tokens
	accessToken, refreshToken, err := as.tokenService.CreateAuthenticationTokens(services.TokenPayload{
		UserId:     identityUser.ID,
		Email:      identityUser.Email,
		ProfileId:  profile.ID,
		ProfileIds: make([]string, 0),
		Role:       role,
	})

	if err != nil {
		return nil, fails.InternalServerError()
	}

	/// 5. Redis to keep up with the tokens that are black listed

	if err := signUpTransaction.Commit().Error; err != nil {
		return nil, fails.InternalServerError()
	}

	return &contracts.AuthResponse{
		Details: contracts.AccountResponse{
			UserId:     identityUser.ID,
			Email:      identityUser.Email,
			ProfileId:  profile.ID,
			ProfileIds: make([]string, 0),
			Role:       role,
		},
		AccessToken:            accessToken.Token,
		AccessTokenExpiration:  accessToken.ExpireAt,
		RefreshToken:           refreshToken.Token,
		RefreshTokenExpiration: refreshToken.ExpireAt,
	}, nil
}
