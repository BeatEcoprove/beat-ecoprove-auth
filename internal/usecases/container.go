package usecases

type UseCases struct {
	Sign           *SignUpUseCase
	Login          *LoginUseCase
	AttachProfile  *AttachProfileUseCase
	RefreshTokens  *RefreshTokensUseCase
	ForgotPassword *ForgotPasswordUseCase
	ResetPassword  *ResetPasswdUseCase
	CheckFields    *CheckFieldUseCase
}
