package helpers

import (
	"log"
	"time"

	"github.com/BeatEcoprove/identityService/internal/domain"
	"github.com/BeatEcoprove/identityService/internal/domain/events"
	"github.com/BeatEcoprove/identityService/internal/repositories"
	"github.com/BeatEcoprove/identityService/pkg/adapters"
	interfaces "github.com/BeatEcoprove/identityService/pkg/domain"
	fails "github.com/BeatEcoprove/identityService/pkg/errors"
)

type (
	CreateProfileInput struct {
		AuthID    string
		Email     string
		Role      domain.AuthRole
		GrantType domain.GrantType
	}

	IProfileCreateService interface {
		CreateProfile(
			trans adapters.Transaction[interfaces.Entity],
			input CreateProfileInput,
		) (*domain.Profile, error)

		MarkProfile(profileID string) error
	}

	ProfileCreateService struct {
		profileRepo repositories.IProfileRepository
		broker      adapters.Broker
		redis       adapters.Redis
	}
)

var profileKeyPrefix = "profile:pending"
var profileExp = time.Duration(15) * time.Minute

func NewProfileKey(profileID string) adapters.RedisKey {
	return adapters.NewRedisKey(profileKeyPrefix, profileID)
}

func NewProfileCreateService(
	profileRepo repositories.IProfileRepository,
	broker adapters.Broker,
	redis adapters.Redis,
) *ProfileCreateService {
	return &ProfileCreateService{
		profileRepo: profileRepo,
		broker:      broker,
		redis:       redis,
	}
}

func (pc ProfileCreateService) checkMainProfile(
	trans adapters.Transaction[interfaces.Entity],
	input CreateProfileInput,
) error {

	if input.GrantType == domain.Main {
		mainProfile, err := pc.profileRepo.GetMainProfileByAuthId(input.AuthID)

		if err != nil {
			return nil
		}

		mainProfile.Role = domain.Sub

		if err := trans.Update(mainProfile); err != nil {
			return fails.InternalServerError()
		}
	}

	return nil
}

func (pc *ProfileCreateService) CreateProfile(
	trans adapters.Transaction[interfaces.Entity],
	input CreateProfileInput,
) (*domain.Profile, error) {
	if err := pc.checkMainProfile(trans, input); err != nil {
		return nil, err
	}

	profile := domain.NewProfile(input.AuthID, input.GrantType)

	if err := trans.Create(profile); err != nil {
		return nil, fails.InternalServerError()
	}

	if err := pc.broker.Publish(&events.UserCreatedEvent{
		AuthID:    input.AuthID,
		ProfileID: profile.ID,
		Email:     input.Email,
		Role:      string(input.Role),
	}, adapters.AuthEventTopic); err != nil {
		log.Printf("failed to send kafka event %s", err.Error())
		return nil, fails.InternalServerError()
	}

	if err := pc.redis.SetValue(NewProfileKey(profile.ID), "1", profileExp); err != nil {
		return nil, err
	}

	return profile, nil
}

func (pc *ProfileCreateService) MarkProfile(profileID string) error {
	_, err := pc.redis.GetAndDelValue(NewProfileKey(profileID))
	return err
}
