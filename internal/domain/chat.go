package domain

import (
	interfaces "github.com/BeatEcoprove/identityService/pkg/domain"
	"gorm.io/gorm"
)

type ChatRole string

const (
	ChatAdmin     ChatRole = "admin"
	ChatModerator ChatRole = "moderator"
	ChatMember    ChatRole = "member"
)

func GetChatRoleByInt(role int) ChatRole {
	switch role {
	case 0:
		return ChatMember
	case 1:
		return ChatModerator
	case 2:
		return ChatAdmin
	default:
		return ChatMember
	}
}

type ChatPermission string

type MemberChatPermission struct {
	interfaces.EntityBase
	GroupId  string
	MemberId string // -> can be the creator // user_id for now
	Role     string
}

func NewMemberChat(groupId, memberId string, role ChatRole) *MemberChatPermission {
	return &MemberChatPermission{
		GroupId:  groupId,
		MemberId: memberId,
		Role:     string(role),
	}
}

func (b *MemberChatPermission) IsMember(memberId string) bool {
	return b.MemberId == memberId
}

func (b *MemberChatPermission) TableName() string {
	return "member_chat"
}

func (u *MemberChatPermission) BeforeCreate(tx *gorm.DB) error {
	u.GetId()

	u.DeletedAt = nil
	return nil
}
