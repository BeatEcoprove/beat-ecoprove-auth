package repositories

import (
	"github.com/BeatEcoprove/identityService/internal/domain"
	interfaces "github.com/BeatEcoprove/identityService/pkg/adapters"
)

type (
	ProfileRepository struct {
		interfaces.RepositoryBase
	}

	IProfileRepository interface {
		interfaces.Repository
		GetMainProfileByAuthId(authId string) (*domain.Profile, error)
		GetAttachProfiles(authId string) ([]domain.Profile, error)
	}
)

func NewProfileRepository(database interfaces.Database) *ProfileRepository {
	return &ProfileRepository{
		RepositoryBase: *interfaces.NewRepositoryBase(database),
	}
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
