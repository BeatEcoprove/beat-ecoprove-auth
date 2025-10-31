package domain

type Permission string

const (
	// profile permissions
	ProfileCreate Permission = "profile:create"
	ProfileView   Permission = "profile:view"
	ProfileDelete Permission = "profile:delete"
	ProfileUpdate Permission = "profile:update"

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

	// messages permissions
	MessageView = "message:view"

	// notifications permissions
	NotificationView = "notification:view"
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
		ProfileView,
		ProfileDelete,
		ProfileUpdate,
		GroupCreate,
		GroupView,
		GroupDelete,
		GroupUpdate,
		MemberKick,
		MemberChangeRole,
		InviteAccept,
		InviteCreate,
		InviteDecline,
		MessageView,
		NotificationView,
	}
}
