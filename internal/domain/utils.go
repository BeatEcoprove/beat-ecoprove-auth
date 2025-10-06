package domain

import "errors"

var (
	ErrUndefinedRole      = errors.New("role not defined")
	ErrUndefinedGrantType = errors.New("grant type not defined")
)

func GetRole(role AuthRole) (string, error) {
	switch role {
	case AuthClient:
		return "client", nil
	case AuthAdmin:
		return "admin", nil
	case AuthOrganization:
		return "organization", nil
	}

	return "", ErrUndefinedRole
}

func GetGrantType(grantType GrantType) (string, error) {
	switch grantType {
	case Main:
		return "main", nil
	case Sub:
		return "sub", nil
	}

	return "", ErrUndefinedGrantType
}

func FilterProfiles(profiles []Profile) (*Profile, []Profile) {
	var mainProfile *Profile
	var subProfiles []Profile = make([]Profile, 0, len(profiles)-1)

	for _, current := range profiles {
		if current.Role == Main {
			mainProfile = &current
			continue
		}

		subProfiles = append(subProfiles, current)
	}

	return mainProfile, subProfiles
}
