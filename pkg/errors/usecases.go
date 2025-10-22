package fails

import "github.com/BeatEcoprove/identityService/pkg/shared"

var (
	USER_ALREADY_EXISTS = shared.NewConflitError(
		"user-already-exists",
		"Email not available",
		"This email is already in use, please choose another one",
	)

	USER_AUTH_FAILED = shared.NewUnauthorizedError(
		"user-auth-failed",
		"Authentication Failed",
		"The password or email must be incorrect, please try again later",
	)

	USER_NOT_FOUND = shared.NewNotFoundError(
		"user-not-found",
		"User not found",
		"The user does not exists",
	)

	ROLE_NOT_FOUND = shared.NewNotFoundError(
		"role-not-found",
		"Role not found",
		"Please provide an valid role",
	)

	GRANT_TYPE_NOT_FOUND = shared.NewNotFoundError(
		"grant-type-not-found",
		"Grant Type not found",
		"Please provide an valid profile grant type",
	)

	PROFILE_DOES_NOT_BELONG_TO_USER = shared.NewConflitError(
		"profile-does-not-belong-to-user",
		"Profile does not belong to user",
		"Please provide an profile id that you have access",
	)

	PROFILE_NOT_FOUND = shared.NewNotFoundError(
		"profile-not-found",
		"Profile not found",
		"The profile does not exists",
	)

	PROFILES_NOT_FOUND = shared.NewNotFoundError(
		"profiles-not-found",
		"Profiles not found",
		"No attached profiles were found",
	)

	CODE_NOT_VALID = shared.NewConflitError(
		"code-not-valid",
		"Code is not valid",
		"Please provide a valid code, try again later",
	)

	GROUP_NOT_FOUND = shared.NewNotFoundError(
		"group-not-found",
		"Group not found",
		"The group that you were looking for doesn't exist",
	)

	MEMBER_NOT_FOUND = shared.NewNotFoundError(
		"member-not-found",
		"Member not found",
		"This member doesn't belong to this group",
	)
)
