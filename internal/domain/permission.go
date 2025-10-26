package domain

type Permission string

const (
	// profile permissions
	ProfileCreate Permission = "profile:create"

	// group permissions
	GroupCreate = "group:create"
	GroupView   = "group:view"
	GroupDelete = "group:delete"
	GroupUpdate = "group:update"

	// member permissions
	MemberKick       = "member:kick"
	MemberChangeRole = "member:change_role"

	// invite permissions
	InviteCreate  = "invite:create"
	InviteAccept  = "invite:accept"
	InviteDecline = "invite:decline"
)

var (
	roles map[AuthRole][]Permission
)

func GetPermissions(identityUser IdentityUser) []string {
	var role AuthRole

	if !identityUser.IsActive {
		role = AuthAnonymous
	} else {
		role = identityUser.GetRole()
	}

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
		GroupDelete,
		GroupUpdate,
		MemberKick,
		MemberChangeRole,
		InviteAccept,
		InviteCreate,
		InviteDecline,
	}
}
