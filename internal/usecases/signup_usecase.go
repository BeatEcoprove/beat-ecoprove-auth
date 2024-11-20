package usecases

import (
	"github.com/BeatEcoprove/identityService/internal/domain"
	"github.com/BeatEcoprove/identityService/internal/repositories"
	"github.com/BeatEcoprove/identityService/pkg/contracts"
	fails "github.com/BeatEcoprove/identityService/pkg/errors"
	"github.com/BeatEcoprove/identityService/pkg/mappers"
	"github.com/BeatEcoprove/identityService/pkg/services"
)

type (
	// input
	SignUpInput struct {
		Email    string
		Password string
		Role     int
	}

	SignUpUseCase struct {
		authRepo     repositories.IAuthRepository
		profileRepo  repositories.IProfileRepository
		tokenService services.ITokenService

		emailService services.IEmailService
	}
)

func NewSignUpUseCase(
	authRepo repositories.IAuthRepository,
	profileRepo repositories.IProfileRepository,
	tokenService services.ITokenService,
	emailService services.IEmailService,
) *SignUpUseCase {
	return &SignUpUseCase{
		authRepo:     authRepo,
		profileRepo:  profileRepo,
		tokenService: tokenService,
		emailService: emailService,
	}
}

func (as *SignUpUseCase) Handle(input SignUpInput) (*contracts.AuthResponse, error) {
	if ok := as.authRepo.ExistsUserWithEmail(input.Email); ok {
		return nil, fails.USER_ALREADY_EXISTS
	}

	if err := services.ValidatePassword(input.Password); err != nil {
		return nil, err
	}

	role, err := domain.GetRole(domain.Role(input.Role))

	if err != nil {
		return nil, fails.ROLE_NOT_FOUND
	}

	identityUser := domain.NewIdentityUser(
		input.Email,
		input.Password,
		domain.Role(input.Role),
	)

	signUpTransaction, err := as.authRepo.BeginTransaction()

	if err != nil {
		return nil, fails.InternalServerError()
	}

	if err := signUpTransaction.Create(identityUser); err != nil {
		return nil, fails.InternalServerError()
	}

	profile := domain.NewProfile(identityUser.ID, domain.Main)

	if err := signUpTransaction.Create(profile); err != nil {
		return nil, fails.InternalServerError()
	}

	accessToken, refreshToken, err := as.tokenService.CreateAuthenticationTokens(services.TokenPayload{
		UserId:     identityUser.ID,
		Email:      identityUser.Email,
		ProfileId:  profile.ID,
		ProfileIds: make([]string, 0),
		Role:       role,
	})

	if err != nil {
		return nil, fails.InternalServerError()
	}

	if err := signUpTransaction.Commit().Error; err != nil {
		return nil, fails.InternalServerError()
	}

	return mappers.ToAuthResponse(
		identityUser,
		profile,
		make([]domain.Profile, 0),
		role,
		accessToken,
		refreshToken,
	), nil
}
