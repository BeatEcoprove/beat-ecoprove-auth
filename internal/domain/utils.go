package domain

func GetRole(role Role) (string, error) {
	switch role {
	case Client:
		return "client", nil
	case Admin:
		return "admin", nil
	case Organization:
		return "organization", nil
	}

	return "", ErrUndefinedRole
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
