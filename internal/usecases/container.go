package usecases

import "github.com/BeatEcoprove/identityService/internal/usecases/helpers"

type UseCases struct {
	Sign             *SignUpUseCase
	Login            *LoginUseCase
	AttachProfile    *AttachProfileUseCase
	RefreshTokens    *RefreshTokensUseCase
	ForgotPassword   *ForgotPasswordUseCase
	ResetPassword    *ResetPasswdUseCase
	CheckFields      *CheckFieldUseCase
	FetchPermissions *FetchGroupUserPermissionsUseCase

	ProfileCreateService helpers.IProfileCreateService
}
