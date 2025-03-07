package repositories

import (
	"github.com/BeatEcoprove/identityService/internal/domain"
	interfaces "github.com/BeatEcoprove/identityService/pkg/adapters"
)

type (
	ProfileRepository struct {
		interfaces.RepositoryBase[*domain.Profile]
	}

	IProfileRepository interface {
		interfaces.Repository[*domain.Profile]
		IsProfileFromUserId(authId, profileId string) bool
		GetMainProfileByAuthId(authId string) (*domain.Profile, error)
		GetAttachProfiles(authId string) ([]domain.Profile, error)
	}
)

func NewProfileRepository(database interfaces.Database) *ProfileRepository {
	return &ProfileRepository{
		RepositoryBase: *interfaces.NewRepositoryBase[*domain.Profile](database),
	}
}

func (repo *ProfileRepository) IsProfileFromUserId(authId, profileId string) bool {
	return repo.Context.Statement.Where("auth_id = ?", authId).Where("id = ?", profileId).First(&domain.Profile{}).Error == nil
}

func (repo *ProfileRepository) GetMainProfileByAuthId(authId string) (*domain.Profile, error) {
	var profile *domain.Profile

	if err := repo.Context.Statement.Where("auth_id = ?", authId).Where("role = ?", domain.Main).First(&profile).Error; err != nil {
		return nil, err
	}

	return profile, nil

}

func (repo *ProfileRepository) GetAttachProfiles(authId string) ([]domain.Profile, error) {
	var profile []domain.Profile

	if err := repo.Context.Statement.Where("auth_id = ?", authId).Find(&profile).Error; err != nil {
		return nil, err
	}

	return profile, nil
}
