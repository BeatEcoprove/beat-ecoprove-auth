package utils

import (
	interfaces "github.com/BeatEcoprove/identityService/pkg/domain"
	"gorm.io/gorm"

	"github.com/BeatEcoprove/identityService/internal/domain"
	"github.com/BeatEcoprove/identityService/pkg/adapters"
	"github.com/stretchr/testify/mock"
)

type (
	MockRepositoryBase[T interfaces.Entity] struct {
		mock.Mock
	}

	MockTransaction struct {
		mock.Mock
		MockRepositoryBase[interfaces.Entity]
	}

	MockAuthRepository struct {
		MockRepositoryBase[*domain.IdentityUser]
	}

	MockProfileRepository struct {
		MockRepositoryBase[*domain.Profile]
	}
)

func (tran *MockTransaction) Rollback() error {
	args := tran.Called()
	return args.Error(0)
}

func (tran *MockTransaction) Commit() error {
	args := tran.Called()
	return args.Error(0)
}

func (repo *MockRepositoryBase[T]) BeginTransaction() (adapters.Transaction[interfaces.Entity], error) {
	args := repo.Called()
	return args.Get(0).(*MockTransaction), args.Error(1)
}

func (repo *MockRepositoryBase[T]) GetOrm() adapters.Orm {
	return &gorm.DB{}
}

func (repo *MockRepositoryBase[T]) Create(entity T) error {
	args := repo.Called()
	return args.Error(0)
}

func (repo *MockRepositoryBase[T]) Delete(entity T) error {
	args := repo.Called()
	return args.Error(0)
}

func (repo *MockRepositoryBase[T]) Update(entity T) error {
	args := repo.Called(entity)
	return args.Error(0)
}

func (repo *MockRepositoryBase[T]) Get(id string) (T, error) {
	args := repo.Called(id)
	return args.Get(0).(T), args.Error(1)
}

func (repo *MockAuthRepository) ExistsUserWithId(id string) bool {
	args := repo.Called(id)
	return args.Bool(0)
}

func (repo *MockAuthRepository) ExistsUserWithEmail(email string) bool {
	args := repo.Called(email)
	return args.Bool(0)
}

func (repo *MockAuthRepository) GetUserByEmail(email string) (*domain.IdentityUser, error) {
	args := repo.Called(email)
	return args.Get(0).(*domain.IdentityUser), args.Error(1)
}

func (repo *MockProfileRepository) IsProfileFromUserId(authId, profileId string) bool {
	args := repo.Called(authId, profileId)
	return args.Bool(0)
}

func (repo *MockProfileRepository) GetMainProfileByAuthId(authId string) (*domain.Profile, error) {
	args := repo.Called(authId)
	return args.Get(0).(*domain.Profile), args.Error(1)
}

func (repo *MockProfileRepository) GetAttachProfiles(authId string) ([]domain.Profile, error) {
	args := repo.Called(authId)
	return args.Get(0).([]domain.Profile), args.Error(1)
}
