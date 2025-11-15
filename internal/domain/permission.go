package domain

type Permission string

const (
	// profile permissions
	ProfileCreate Permission = "profile:create"
	ProfileView   Permission = "profile:view"
	ProfileDelete Permission = "profile:delete"
	ProfileUpdate Permission = "profile:update"

	// workers
	WorkerCreate Permission = "worker:create"
	WorkerView   Permission = "worker:view"
	WorkerDelete Permission = "worker:delete"
	WorkerSwitch Permission = "worker:switch"

	// stores
	StoreCreate Permission = "store:create"
	StoreView   Permission = "store:view"
	StoreDelete Permission = "store:delete"

	// stores
	ServiceView   Permission = "service:view"
	ServiceUpdate Permission = "service:update"

	// Ratings
	RatingCreate Permission = "rating:create"
	RatingView   Permission = "rating:view"

	// Adverts
	AdvertCreate Permission = "advert:create"
	AdvertView   Permission = "advert:view"
	AdvertDelete Permission = "advert:delete"

	// Providers
	ProviderCreate Permission = "provider:create"
	ProviderView   Permission = "provider:view"

	// Bucket
	BucketCreate Permission = "bucket:create"
	BucketView   Permission = "bucket:view"
	BucketDelete Permission = "bucket:delete"

	// Cloth
	ClothCreate  Permission = "cloth:create"
	ClothView    Permission = "cloth:view"
	ClothDelete  Permission = "cloth:delete"
	ClothHistory Permission = "cloth:history"

	// Outfit
	OutfitView Permission = "outfit:view"

	// Feedback
	FeedbackCreate Permission = "feedback:create"

	// Maintenance
	MaintenanceCreate Permission = "maintenance:create"

	// Currency
	CurrencyConvert Permission = "currency:convert"

	// Color
	ColorView Permission = "color:view"

	// Brand
	BrandView   Permission = "brand:view"
	BrandCreate Permission = "brand:create"

	// group permissions
	GroupCreate Permission = "group:create"
	GroupView   Permission = "group:view"
	GroupDelete Permission = "group:delete"
	GroupUpdate Permission = "group:update"

	// member permissions
	MemberKick       Permission = "member:kick"
	MemberChangeRole Permission = "member:change_role"

	// invite permissions
	InviteCreate  Permission = "invite:create"
	InviteAccept  Permission = "invite:accept"
	InviteDecline Permission = "invite:decline"

	// messages permissions
	MessageView Permission = "message:view"

	// notifications permissions
	NotificationView Permission = "notification:view"
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
		ProfileUpdate,
		ProfileDelete,
		BucketCreate,
		BucketDelete,
		BucketView,
		ClothCreate,
		ClothView,
		ClothDelete,
		ClothHistory,
		OutfitView,
		FeedbackCreate,
		BrandView,
		BrandCreate,
		ColorView,
		CurrencyConvert,
		MaintenanceCreate,
		AdvertView,
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
		ProviderView,
		ServiceView,
		ServiceUpdate,
	}
}
