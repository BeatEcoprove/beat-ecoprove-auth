package usecases

import (
	"github.com/BeatEcoprove/identityService/internal/domain"
	"github.com/BeatEcoprove/identityService/internal/repositories"
	"github.com/BeatEcoprove/identityService/pkg/contracts"
	fails "github.com/BeatEcoprove/identityService/pkg/errors"
	"github.com/BeatEcoprove/identityService/pkg/mappers"
	"github.com/BeatEcoprove/identityService/pkg/services"
)

type (
	// input
	LoginInput struct {
		Email    string
		Password string
	}

	LoginUseCase struct {
		authRepo     repositories.IAuthRepository
		profileRepo  repositories.IProfileRepository
		tokenService services.ITokenService
	}
)

func NewLoginUseCase(
	authRepo repositories.IAuthRepository,
	profileRepo repositories.IProfileRepository,
	tokenService services.ITokenService,
) *LoginUseCase {
	return &LoginUseCase{
		authRepo:     authRepo,
		profileRepo:  profileRepo,
		tokenService: tokenService,
	}
}

func (as *LoginUseCase) Handle(input LoginInput) (*contracts.AuthResponse, error) {
	// Setps to create an user
	/// 1. Exists already this user registered?
	if ok := as.authRepo.ExistsUserWithEmail(input.Email); !ok {
		return nil, fails.USER_AUTH_FAILED
	}

	if err := services.ValidatePassword(input.Password); err != nil {
		return nil, fails.USER_AUTH_FAILED
	}

	identityUser, err := as.authRepo.GetUserByEmail(input.Email)

	if err != nil {
		return nil, fails.USER_AUTH_FAILED
	}

	if !services.CheckPasswordHash(input.Password, identityUser.Salt, identityUser.Password) {
		return nil, fails.USER_AUTH_FAILED
	}

	/// Get the correct profile
	attachedProfiles, err := as.profileRepo.GetAttachProfiles(identityUser.ID)

	if err != nil {
		return nil, fails.USER_AUTH_FAILED
	}

	role, err := domain.GetRole(identityUser.Role)

	if err != nil {
		return nil, fails.USER_AUTH_FAILED
	}

	// filter profiles
	mainProfile, subProfiles := domain.FilterProfiles(attachedProfiles)

	/// 4. Generate the authorization tokens -> Generate Tokens
	accessToken, refreshToken, err := as.tokenService.CreateAuthenticationTokens(services.TokenPayload{
		UserId:     identityUser.ID,
		Email:      identityUser.Email,
		ProfileId:  mainProfile.ID,
		ProfileIds: mappers.MapProfileIdsToString(subProfiles),
		Role:       role,
	})

	if err != nil {
		return nil, fails.InternalServerError()
	}

	return mappers.ToAuthResponse(
		identityUser,
		mainProfile,
		subProfiles,
		role,
		accessToken,
		refreshToken,
	), nil
}
