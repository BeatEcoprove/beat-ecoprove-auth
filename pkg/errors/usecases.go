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

	ROLE_NOT_FOUND = shared.NewNotFoundError(
		"role-not-found",
		"Role not found",
		"PLease provide an valid role",
	)
)
