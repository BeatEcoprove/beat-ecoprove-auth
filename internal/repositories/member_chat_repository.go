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
		IsMember(id, memberID string) bool
		GetByGroupId(id string) (*domain.MemberChatPermission, error)
		GetPermissions(id string) ([]domain.MemberChatPermission, error)
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

func (repo *MemberChatRepository) GetPermissions(id string) ([]domain.MemberChatPermission, error) {
	var entries []domain.MemberChatPermission

	if err := repo.Context.Statement.Where("group_id = ?", id).Find(&entries).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

func (repo *MemberChatRepository) GetByGroupId(id string) (*domain.MemberChatPermission, error) {
	var groupEntry *domain.MemberChatPermission

	if err := repo.Context.Statement.Where("group_id = ?", id).First(&groupEntry).Error; err != nil {
		return nil, err
	}

	return groupEntry, nil
}

func (repo *MemberChatRepository) IsMember(id, memberId string) bool {
	return repo.Context.Statement.Where("id = ?", id).Where("member_id = ?", memberId).First(&domain.MemberChatPermission{}).Error == nil
}
