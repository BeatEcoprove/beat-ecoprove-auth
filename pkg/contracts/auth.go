package contracts

type (
	RefreshTokensRequest struct {
		ProfileId string `validate:"required,uuid"`
	}

	ForgotPasswordRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	ResetPasswordRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Code     string `json:"code"`
		Password string `json:"password" validate:"required,min=8"`
	}

	AttachProfileRequest struct {
		ProfileGrantType int `json:"grant_type"`
	}

	ProfileResponse struct {
		ProfileId string `json:"profile_id"`
	}

	LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	SignUpRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
		Role     int    `json:"role" validate:"numeric"`
	}

	AccountResponse struct {
		UserId     string   `json:"user_id"`
		Email      string   `json:"email"`
		ProfileId  string   `json:"profile_id"`
		ProfileIds []string `json:"profile_ids"`
		Role       string   `json:"role"`
	}

	AuthResponse struct {
		Details                AccountResponse `json:"details"`
		AccessToken            string          `json:"access_token"`
		AccessTokenExpiration  int             `json:"access_token_exp"`
		RefreshToken           string          `json:"refresh_token"`
		RefreshTokenExpiration int             `json:"refresh_token_exp"`
	}

	GenericResponse struct {
		Message string `json:"message"`
	}
)
