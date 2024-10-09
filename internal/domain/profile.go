package domain

import (
	interfaces "github.com/BeatEcoprove/identityService/pkg/domain"
	"gorm.io/gorm"
)

type GrantType int

const (
	Main GrantType = iota
	Sub
)

type Profile struct {
	interfaces.EntityBase
	AuthId string
	Role   GrantType
}

func NewProfile(authId string, role GrantType) *Profile {
	return &Profile{
		AuthId: authId,
		Role:   role,
	}
}

func (b *Profile) TableName() string {
	return "profiles"
}

func (u *Profile) BeforeCreate(tx *gorm.DB) error {
	u.GetId()

	u.DeletedAt = nil
	return nil
}
