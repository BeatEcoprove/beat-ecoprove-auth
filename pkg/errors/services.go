package fails

import (
	"github.com/BeatEcoprove/identityService/pkg/shared"
)

var (
	PASSWORD_PROVIDE = shared.NewConflitError(
		"password-pv",
		"Auth.Validation.ProvidePassword.Title",
		"Auth.Validation.ProvidePassword.Description",
	)

	PASSWORD_BTW_6_16 = shared.NewConflitError(
		"password-btw-6-16",
		"Auth.Validation.PasswordBetween6And16.Title",
		"Auth.Validation.PasswordBetween6And16.Description",
	)

	PASSWORD_MUST_CONTAIN_ONE_NUMBER = shared.NewConflitError(
		"password-mst-number",
		"Auth.Validation.PasswordMustContainOneNumber.Title",
		"Auth.Validation.PasswordMustContainOneNumber.Description",
	)

	PASSWORD_MUST_CONTAIN_ONE_CAPITAL = shared.NewConflitError(
		"password-mst-capital",
		"Auth.Validation.PasswordMustContainOneCapital.Title",
		"Auth.Validation.PasswordMustContainOneCapital.Description",
	)

	PASSWORD_MUST_CONTAIN_NON_CAPITAL = shared.NewConflitError(
		"password-mst-n-capital",
		"Auth.Validation.PasswordMustContainNonCapital.Title",
		"Auth.Validation.PasswordMustContainNonCapital.Description",
	)

	BAD_UUID = shared.NewConflitError(
		"bad-uuid",
		"Auth.Validation.InvalidUUID.Title",
		"Auth.Validation.InvalidUUID.Description",
	)

	BAD_EMAIL = shared.NewConflitError(
		"bad-email",
		"Auth.Validation.InvalidEmail.Title",
		"Auth.Validation.InvalidEmail.Description",
	)
)
