package handlers

import (
	"fmt"

	"github.com/BeatEcoprove/identityService/internal/domain"
	"github.com/BeatEcoprove/identityService/internal/domain/events"
	"github.com/BeatEcoprove/identityService/internal/repositories"
	"github.com/BeatEcoprove/identityService/internal/usecases/helpers"
)

type ProfileCreatedHandler struct {
	authRepo             repositories.IAuthRepository
	profileRepo          repositories.IProfileRepository
	profileCreateService helpers.IProfileCreateService
}

func NewProfileCreatedHandler(
	authRepo repositories.IAuthRepository,
	profileRepo repositories.IProfileRepository,
	profileCreateService helpers.IProfileCreateService,
) *ProfileCreatedHandler {
	return &ProfileCreatedHandler{
		authRepo:             authRepo,
		profileRepo:          profileRepo,
		profileCreateService: profileCreateService,
	}
}

func (p *ProfileCreatedHandler) canCreateProfile(user domain.IdentityUser) bool {
	profiles, err := p.profileRepo.GetAttachProfiles(user.ID)

	if err != nil {
		return false
	}

	return len(profiles) > 1 && user.IsActive || !user.IsActive && len(profiles) == 1
}

func (p *ProfileCreatedHandler) Call(payload any) error {
	event, ok := payload.(*events.ProfileCreatedEvent)

	if !ok {
		return fmt.Errorf("failed to cast, payload %+v", payload)
	}

	foundAuth, err := p.authRepo.Get(event.AuthId)

	if err != nil {
		return fmt.Errorf("failed to find profile with id %s", event.ProfileId)
	}

	if !p.canCreateProfile(*foundAuth) {
		return fmt.Errorf("faulted account, profile can't be created")
	}

	if !p.profileRepo.IsProfileFromUserId(event.AuthId, event.ProfileId) {
		return fmt.Errorf("access revoked to profile")
	}

	if err := p.profileCreateService.MarkProfile(event.ProfileId); err != nil {
		return fmt.Errorf("failed to activate profile")
	}

	foundAuth.IsActive = true

	if err := p.authRepo.Update(foundAuth); err != nil {
		return fmt.Errorf("failed to activate account")
	}

	return nil
}
