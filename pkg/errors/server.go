package fails

import (
	"github.com/BeatEcoprove/identityService/pkg/shared"
)

var (
	DONT_HAVE_ACCESS_TO_RESOURCE = shared.NewForbiddenError(
		"dont-have-access-to-resource",
		"Auth.Http.DontHaveAccessToResource.Title",
		"Auth.Http.DontHaveAccessToResource.Description",
	)

	INVALID_ACCESS_TOKEN = shared.NewUnauthorizedError(
		"invalid-access-token",
		"Auth.Http.InvalidAccessToken.Title",
		"Auth.Http.InvalidAccessToken.Description",
	)

	INVALID_REFRESH_TOKEN = shared.NewUnauthorizedError(
		"invalid-refresh-token",
		"Auth.Http.InvalidRefreshToken.Title",
		"Auth.Http.InvalidRefreshToken.Description",
	)
)

func InternalServerError() *shared.Error {
	return shared.NewInternalError(
		"internal-server",
		"Auth.Http.InternalServerError.Title",
		"Auth.Http.InternalServerError.Description",
	)
}
