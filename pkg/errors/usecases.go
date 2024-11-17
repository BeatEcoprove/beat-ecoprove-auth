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
)
