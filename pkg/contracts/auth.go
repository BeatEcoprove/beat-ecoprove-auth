package contracts

type (
	CheckEmailFieldRequest struct {
		Email string `validate:"required,email"`
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
		ProfileID string `json:"profile_id"`
	}

	GroupPermissionsResponse struct {
		MemberID string `json:"member_id"`
		Role     string `json:"role"`
	}

	GroupPermissionsRequest struct {
		GroupID  string `json:"group_id" validate:"uuid"`
		MemberID string `json:"member_id" validate:"uuid"`
	}

	TokenRequest struct {
		GrantType string `json:"grant_type" form:"grant_type"`
	}

	LoginRequest struct {
		TokenRequest
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	RefreshTokenRequest struct {
		TokenRequest
		Token     string `json:"token" form:"token"`
		ProfileID string `json:"profile_id" form:"profile_id" validate:"omitempty,uuid"`
	}

	SignUpRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
		Role     int    `json:"role" validate:"numeric"`
	}

	AccountResponse struct {
		UserID     string   `json:"user_id"`
		Email      string   `json:"email"`
		ProfileID  string   `json:"profile_id"`
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
