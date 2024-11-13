package contracts

type (
	LoginRequest struct {
		Email    string
		Password string
	}

	SignUpRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
		Role     int    `json:"role" validate:"numeric"`
	}

	AccountResponse struct {
		UserId     string
		Email      string
		ProfileId  string
		ProfileIds []string
		Role       string
	}

	AuthResponse struct {
		Details                AccountResponse
		AccessToken            string
		AccessTokenExpiration  int
		RefreshToken           string
		RefreshTokenExpiration int
	}
)
