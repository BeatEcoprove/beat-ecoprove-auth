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
	RefreshTokensInput struct {
		AuthId    string
		ProfileId string
	}

	RefreshTokensUseCase struct {
		authRepo     repositories.IAuthRepository
		profileRepo  repositories.IProfileRepository
		tokenService services.ITokenService
	}
)

func NewRefreshTokensUseCase(
	authRepo repositories.IAuthRepository,
	profileRepo repositories.IProfileRepository,
	tokenService services.ITokenService,
) *RefreshTokensUseCase {
	return &RefreshTokensUseCase{
		authRepo:     authRepo,
		profileRepo:  profileRepo,
		tokenService: tokenService,
	}
}

func (rtu *RefreshTokensUseCase) Handle(request RefreshTokensInput) (*contracts.AuthResponse, error) {
	identityUser, err := rtu.authRepo.Get(request.AuthId)

	if err != nil {
		return nil, fails.USER_NOT_FOUND
	}

	role, err := domain.GetRole(identityUser.Role)

	if err != nil {
		return nil, fails.ROLE_NOT_FOUND
	}

	mainProfile, subProfiles, err := rtu.getProfiles(request.AuthId, request.ProfileId)

	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := rtu.tokenService.CreateAuthenticationTokens(services.TokenPayload{
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

func (rtu *RefreshTokensUseCase) getProfiles(authId, profileId string) (*domain.Profile, []domain.Profile, error) {
	if profileId != "" {
		if ok := rtu.profileRepo.IsProfileFromUserId(authId, profileId); !ok {
			return nil, nil, fails.PROFILE_DOES_NOT_BELONG_TO_USER
		}

		mainProfile, err := rtu.profileRepo.Get(profileId)

		if err != nil {
			return nil, nil, fails.PROFILE_NOT_FOUND
		}

		return mainProfile, make([]domain.Profile, 0), nil
	}

	attachedProfiles, err := rtu.profileRepo.GetAttachProfiles(authId)

	if err != nil {
		return nil, nil, fails.PROFILES_NOT_FOUND
	}

	mainProfile, subProfiles := domain.FilterProfiles(attachedProfiles)
	return mainProfile, subProfiles, nil
}
