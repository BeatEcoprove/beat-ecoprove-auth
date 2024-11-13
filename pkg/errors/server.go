package fails

import (
	"github.com/BeatEcoprove/identityService/pkg/shared"
)

func InternalServerError() *shared.Error {
	return shared.NewInternalError(
		"internal-server",
		"Something went wrong, please try again later",
		"The server could not process this request, please try again later",
	)
}
