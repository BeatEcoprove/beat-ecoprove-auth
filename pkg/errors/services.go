package fails

import (
	"github.com/BeatEcoprove/identityService/pkg/shared"
)

var (
	PASSWORD_PROVIDE = shared.NewConflitError(
		"password-pv",
		"Password Validation",
		"Provide password field",
	)

	PASSWORD_BTW_6_16 = shared.NewConflitError(
		"password-btw-6-16",
		"Password Validation",
		"Password must be between 6 and 16 characters",
	)

	PASSWORD_MUST_CONTAIN_ONE_NUMBER = shared.NewConflitError(
		"password-mst-number",
		"Password Validation",
		"Password must contain at least one number",
	)

	PASSWORD_MUST_CONTAIN_ONE_CAPITAL = shared.NewConflitError(
		"password-mst-capital",
		"Password Validation",
		"Password must contain at least one capital letter",
	)

	PASSWORD_MUST_CONTAIN_NON_CAPITAL = shared.NewConflitError(
		"password-mst-n-capital",
		"Password Validation",
		"Password must contain at least one non-capital letter",
	)
)
