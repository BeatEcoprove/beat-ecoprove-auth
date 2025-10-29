package fails

import "github.com/BeatEcoprove/identityService/pkg/shared"

var (
	USER_ALREADY_EXISTS = shared.NewConflitError(
		"user-already-exists",
		"Auth.User.AlreadyExists.Title",
		"Auth.User.AlreadyExists.Description",
	)

	USER_AUTH_FAILED = shared.NewUnauthorizedError(
		"user-auth-failed",
		"Auth.User.AuthFailed.Title",
		"Auth.User.AuthFailed.Description",
	)

	USER_NOT_FOUND = shared.NewNotFoundError(
		"user-not-found",
		"Auth.User.NotFound.Title",
		"Auth.User.NotFound.Description",
	)

	ROLE_NOT_FOUND = shared.NewNotFoundError(
		"role-not-found",
		"Auth.Role.NotFound.Title",
		"Auth.Role.NotFound.Description",
	)

	GRANT_TYPE_NOT_FOUND = shared.NewNotFoundError(
		"grant-type-not-found",
		"Auth.GrantType.NotFound.Title",
		"Auth.GrantType.NotFound.Description",
	)

	PROFILE_DOES_NOT_BELONG_TO_USER = shared.NewConflitError(
		"profile-does-not-belong-to-user",
		"Auth.Profile.DoesNotBelongToUser.Title",
		"Auth.Profile.DoesNotBelongToUser.Description",
	)

	PROFILE_NOT_FOUND = shared.NewNotFoundError(
		"profile-not-found",
		"Auth.Profile.NotFound.Title",
		"Auth.Profile.NotFound.Description",
	)

	PROFILES_NOT_FOUND = shared.NewNotFoundError(
		"profiles-not-found",
		"Auth.Profiles.NotFound.Title",
		"Auth.Profiles.NotFound.Description",
	)

	CODE_NOT_VALID = shared.NewConflitError(
		"code-not-valid",
		"Auth.Code.NotValid.Title",
		"Auth.Code.NotValid.Description",
	)

	GROUP_NOT_FOUND = shared.NewNotFoundError(
		"group-not-found",
		"Auth.Group.NotFound.Title",
		"Auth.Group.NotFound.Description",
	)

	MEMBER_NOT_FOUND = shared.NewNotFoundError(
		"member-not-found",
		"Auth.Member.NotFound.Title",
		"Auth.Member.NotFound.Description",
	)
)
