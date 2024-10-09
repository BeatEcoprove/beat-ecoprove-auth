package repositories

import (
	"github.com/BeatEcoprove/identityService/internal/domain"
	interfaces "github.com/BeatEcoprove/identityService/pkg/adapters"
)

type AuthRepository struct {
	interfaces.RepositoryBase
}

func NewAuthRepository(database interfaces.Database) *AuthRepository {
	return &AuthRepository{
		RepositoryBase: *interfaces.NewRepositoryBase(database),
	}
}

func (repo *AuthRepository) ExistsUserWithEmail(email string) bool {
	if err := repo.Context.Statement.Where("email = ?", email).First(&domain.IdentityUser{}).Error; err != nil {
		return false
	}

	return true
}
