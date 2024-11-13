package repositories

import (
	"github.com/BeatEcoprove/identityService/internal/domain"
	interfaces "github.com/BeatEcoprove/identityService/pkg/adapters"
)

type (
	AuthRepository struct {
		interfaces.RepositoryBase
	}

	IAuthRepository interface {
		interfaces.Repository
		ExistsUserWithEmail(email string) bool
	}
)

func NewAuthRepository(database interfaces.Database) *AuthRepository {
	return &AuthRepository{
		RepositoryBase: *interfaces.NewRepositoryBase(database),
	}
}

func (repo *AuthRepository) ExistsUserWithEmail(email string) bool {
	return repo.Context.Statement.Where("email = ?", email).First(&domain.IdentityUser{}).Error == nil
}
