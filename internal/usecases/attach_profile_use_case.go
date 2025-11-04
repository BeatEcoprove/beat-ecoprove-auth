package usecases

import (
	"github.com/BeatEcoprove/identityService/internal/domain"
	"github.com/BeatEcoprove/identityService/internal/repositories"
	"github.com/BeatEcoprove/identityService/internal/usecases/helpers"
	"github.com/BeatEcoprove/identityService/pkg/contracts"
	fails "github.com/BeatEcoprove/identityService/pkg/errors"
	"github.com/BeatEcoprove/identityService/pkg/mappers"
	"github.com/BeatEcoprove/identityService/pkg/services"
)

type (
	// input
	AttachProfileInput struct {
		AuthId           string
		ProfileGrantType int
	}

	AttachProfileUseCase struct {
		authRepo             repositories.IAuthRepository
		profileRepo          repositories.IProfileRepository
		tokenService         services.ITokenService
		createProfileService helpers.IProfileCreateService
	}
)

func NewAttachProfileUseCase(
	authRepo repositories.IAuthRepository,
	profileRepo repositories.IProfileRepository,
	tokenService services.ITokenService,
	createProfileService helpers.IProfileCreateService,
) *AttachProfileUseCase {
	return &AttachProfileUseCase{
		authRepo:             authRepo,
		profileRepo:          profileRepo,
		tokenService:         tokenService,
		createProfileService: createProfileService,
	}
}

func (apu *AttachProfileUseCase) Handle(request AttachProfileInput) (*contracts.AuthResponse, error) {
	grantType := domain.GrantType(request.ProfileGrantType)
	_, err := domain.GetGrantType(grantType)

	if err != nil {
		return nil, fails.GRANT_TYPE_NOT_FOUND
	}

	identityUser, err := apu.authRepo.Get(request.AuthId)

	if err != nil {
		return nil, fails.USER_NOT_FOUND
	}

	createProfileTransaction, err := apu.profileRepo.BeginTransaction()

	if err != nil {
		return nil, fails.InternalServerError()
	}

	profile, err := apu.createProfileService.CreateProfile(createProfileTransaction, helpers.CreateProfileInput{
		AuthID:    identityUser.ID,
		Email:     identityUser.Email,
		Role:      identityUser.GetRole(),
		GrantType: grantType,
	})

	if err != nil {
		return nil, fails.InternalServerError()
	}

	identityUser.IsActive = false
	accessToken, refreshToken, err := apu.tokenService.CreateAuthenticationTokens(services.TokenPayload{
		UserID:     identityUser.ID,
		Email:      identityUser.Email,
		ProfileID:  profile.ID,
		Scope:      domain.GetPermissions(*identityUser),
		ProfileIds: make([]string, 0),
		Role:       string(identityUser.GetRole()),
	})

	if err != nil {
		return nil, fails.InternalServerError()
	}

	if err := createProfileTransaction.Commit(); err != nil {
		return nil, fails.InternalServerError()
	}

	return mappers.ToAuthResponse(
		identityUser,
		accessToken,
		refreshToken,
	), nil
}
