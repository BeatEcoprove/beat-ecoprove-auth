package domain

import (
	interfaces "github.com/BeatEcoprove/identityService/pkg/domain"
	"github.com/BeatEcoprove/identityService/pkg/services"
	"gorm.io/gorm"
)

type AuthRole string

const (
	AuthAnonymous    AuthRole = "anonymous"
	AuthClient       AuthRole = "client"
	AuthOrganization AuthRole = "organization"
	AuthAdmin        AuthRole = "admin"
)

type IdentityUser struct {
	interfaces.EntityBase
	Email    string
	Password string
	Salt     string `gorm:"column:salt"`
	IsActive bool
	role     AuthRole
}

func NewIdentityUser(email, password string, role AuthRole) *IdentityUser {
	return &IdentityUser{
		Email:    email,
		Password: password,
		role:     role,
		IsActive: false,
	}
}

func (b *IdentityUser) GetRole() AuthRole {
	if b.IsActive {
		return b.role
	}

	return AuthAnonymous
}

func (b *IdentityUser) TableName() string {
	return "auths"
}

func (u *IdentityUser) SetPassword(value string) error {
	salt, err := services.GenerateSalt(services.SALT_COST)

	if err != nil {
		return err
	}

	password, err := services.HashPassword(value, salt)

	if err != nil {
		return err
	}

	u.Salt = salt
	u.Password = password
	return nil
}

func (u *IdentityUser) BeforeCreate(tx *gorm.DB) error {
	u.GetId()

	u.SetPassword(u.Password)
	u.DeletedAt = nil
	return nil
}
