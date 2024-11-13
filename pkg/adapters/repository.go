package adapters

import (
	"github.com/BeatEcoprove/identityService/pkg/domain"
)

type Repository interface {
	Create(entity domain.Entity) error
	Delete(entity domain.Entity) error
	Get(id string) (domain.Entity, error)

	BeginTransaction() (*Transaction, error)
	GetOrm() Orm
}

type (
	Transaction struct {
		Repository
	}

	RepositoryBase struct {
		Context Orm
	}
)

func NewTransaction(repository Repository) *Transaction {
	return &Transaction{
		Repository: repository,
	}
}

func (tran *Transaction) Rollback() Orm {
	return tran.GetOrm().Statement.Rollback()
}

func (tran *Transaction) Commit() Orm {
	return tran.GetOrm().Statement.Commit()
}

func NewRepositoryBase(database Database) *RepositoryBase {
	return &RepositoryBase{
		Context: database.GetOrm(),
	}
}

func (repo *RepositoryBase) BeginTransaction() (*Transaction, error) {
	cloneRepo := *repo
	cloneRepo.Context = repo.Context.Statement.Begin()

	return NewTransaction(&cloneRepo), nil
}

func (repo *RepositoryBase) GetOrm() Orm {
	return repo.Context
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
