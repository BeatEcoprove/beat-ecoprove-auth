package adapters

import (
	"github.com/BeatEcoprove/identityService/pkg/domain"
)

type Repository interface {
	Create(entity domain.Entity) error
	Delete(entity domain.Entity) error
	Get(id string) (domain.Entity, error)
}

type RepositoryBase struct {
	Context Orm
}

func NewRepositoryBase(database Database) *RepositoryBase {
	return &RepositoryBase{
		Context: database.GetOrm(),
	}
}

func (repo *RepositoryBase) Create(entity domain.Entity) error {
	if err := repo.Context.Statement.Create(entity).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RepositoryBase) Delete(entity domain.Entity) error {
	if err := repo.Context.Statement.Delete(entity).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RepositoryBase) Get(id string) (domain.Entity, error) {
	var entity domain.Entity

	if err := repo.Context.Statement.Where("id = ?", id).First(&entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}
