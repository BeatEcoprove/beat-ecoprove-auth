package fails

import (
	"github.com/BeatEcoprove/identityService/pkg/shared"
)

var (
	DONT_HAVE_ACCESS_TO_RESOURCE = shared.NewForbiddenError(
		"dont-have-access-to-resource",
		"Authentication Failed",
		"You don't have access to this resource, please try again later",
	)

	INVALID_ACCESS_TOKEN = shared.NewUnauthorizedError(
		"invalid-access-token",
		"Authorization Failed",
		"Please provide an valid access token, try again later",
	)

	INVALID_REFRESH_TOKEN = shared.NewUnauthorizedError(
		"invalid-refresh-token",
		"Authorization Failed",
		"Please provide an valid refresh token, try again later",
	)
)

func InternalServerError() *shared.Error {
	return shared.NewInternalError(
		"internal-server",
		"Something went wrong",
		"The server could not process this request, please try again later",
	)
}
