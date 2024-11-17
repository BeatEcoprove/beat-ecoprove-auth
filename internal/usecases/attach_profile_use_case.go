package usecases

import (
	"github.com/BeatEcoprove/identityService/internal/domain"
	"github.com/BeatEcoprove/identityService/internal/repositories"
	"github.com/BeatEcoprove/identityService/pkg/contracts"
	fails "github.com/BeatEcoprove/identityService/pkg/errors"
)

type (
	// input
	AttachProfileInput struct {
		AuthId           string
		ProfileGrantType int
	}

	AttachProfileUseCase struct {
		authRepo    repositories.IAuthRepository
		profileRepo repositories.IProfileRepository
	}
)

func NewAttachProfileUseCase(
	authRepo repositories.IAuthRepository,
	profileRepo repositories.IProfileRepository,
) *AttachProfileUseCase {
	return &AttachProfileUseCase{
		authRepo:    authRepo,
		profileRepo: profileRepo,
	}
}

func (apu *AttachProfileUseCase) Handle(request AttachProfileInput) (*contracts.ProfileResponse, error) {
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

	if grantType == domain.Main {
		mainProfile, err := apu.profileRepo.GetMainProfileByAuthId(identityUser.ID)

		if err != nil {
			return nil, fails.USER_NOT_FOUND
		}

		mainProfile.Role = domain.Sub

		if err := createProfileTransaction.Update(mainProfile); err != nil {
			return nil, fails.InternalServerError()
		}
	}

	profile := domain.NewProfile(
		identityUser.ID,
		grantType,
	)

	if err := createProfileTransaction.Create(profile); err != nil {
		return nil, fails.InternalServerError()
	}

	if err := createProfileTransaction.Commit().Error; err != nil {
		return nil, fails.InternalServerError()
	}

	return &contracts.ProfileResponse{
		ProfileId: profile.ID,
	}, nil
}
