package repositories

import (
	"github.com/BeatEcoprove/identityService/internal/domain"
	interfaces "github.com/BeatEcoprove/identityService/pkg/adapters"
)

type (
	AuthRepository struct {
		interfaces.RepositoryBase[*domain.IdentityUser]
	}

	IAuthRepository interface {
		interfaces.Repository[*domain.IdentityUser]
		ExistsUserWithId(id string) bool
		ExistsUserWithEmail(email string) bool
		GetUserByEmail(email string) (*domain.IdentityUser, error)
	}
)

func NewAuthRepository(database interfaces.Database) *AuthRepository {
	return &AuthRepository{
		RepositoryBase: *interfaces.NewRepositoryBase[*domain.IdentityUser](database),
	}
}

func (repo *AuthRepository) ExistsUserWithId(id string) bool {
	return repo.Context.Statement.Where("id = ?", id).First(&domain.IdentityUser{}).Error == nil
}

func (repo *AuthRepository) ExistsUserWithEmail(email string) bool {
	return repo.Context.Statement.Where("email = ?", email).First(&domain.IdentityUser{}).Error == nil
}

func (repo *AuthRepository) GetUserByEmail(email string) (*domain.IdentityUser, error) {
	var identityUser *domain.IdentityUser

	if err := repo.Context.Statement.Where("email = ?", email).First(&identityUser).Error; err != nil {
		return nil, err
	}

	return identityUser, nil
}
