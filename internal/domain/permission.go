package domain

type Permission string

const (
	// profile permissions
	ProfileCreate Permission = "profile:create"

	// group permissions
	GroupCreate = "group:create"
	GroupView   = "group:view"
)

var (
	roles map[AuthRole][]Permission
)

func GetPermissions(role AuthRole) []string {
	permissions := roles[role]
	result := make([]string, len(permissions))

	for i, p := range permissions {
		result[i] = string(p)
	}

	return result
}

func InitPermissions() {
	roles = make(map[AuthRole][]Permission)

	roles[AuthAnonymous] = []Permission{
		ProfileCreate,
	}

	roles[AuthClient] = []Permission{
		ProfileCreate,
		GroupCreate,
		GroupView,
	}
}
