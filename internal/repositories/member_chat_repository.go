package repositories

import (
	"github.com/BeatEcoprove/identityService/internal/domain"
	interfaces "github.com/BeatEcoprove/identityService/pkg/adapters"
)

type (
	MemberChatRepository struct {
		interfaces.RepositoryBase[*domain.MemberChatPermission]
	}

	IMemberChatRepository interface {
		interfaces.Repository[*domain.MemberChatPermission]
		ExistsGroupById(id string) bool
	}
)

func NewMemberChatRepository(database interfaces.Database) *MemberChatRepository {
	return &MemberChatRepository{
		RepositoryBase: *interfaces.NewRepositoryBase[*domain.MemberChatPermission](database),
	}
}

func (repo *MemberChatRepository) ExistsGroupById(id string) bool {
	return repo.Context.Statement.Where("id = ?", id).First(&domain.MemberChatPermission{}).Error == nil
}
