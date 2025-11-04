package usecases

import (
	"log"

	"github.com/BeatEcoprove/identityService/internal/domain"
	"github.com/BeatEcoprove/identityService/internal/repositories"
	"github.com/BeatEcoprove/identityService/internal/usecases/helpers"
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
		Role     string
	}

	SignUpUseCase struct {
		authRepo     repositories.IAuthRepository
		profileRepo  repositories.IProfileRepository
		tokenService services.ITokenService

		emailService        services.IEmailService
		createProfileHelper helpers.IProfileCreateService
	}
)

func NewSignUpUseCase(
	authRepo repositories.IAuthRepository,
	profileRepo repositories.IProfileRepository,
	tokenService services.ITokenService,
	emailService services.IEmailService,
	createProfileHelper helpers.IProfileCreateService,
) *SignUpUseCase {
	return &SignUpUseCase{
		authRepo:            authRepo,
		profileRepo:         profileRepo,
		tokenService:        tokenService,
		emailService:        emailService,
		createProfileHelper: createProfileHelper,
	}
}

func (as *SignUpUseCase) Handle(input SignUpInput) (*contracts.AuthResponse, error) {
	if ok := as.authRepo.ExistsUserWithEmail(input.Email); ok {
		return nil, fails.USER_ALREADY_EXISTS
	}

	if err := services.ValidatePassword(input.Password); err != nil {
		return nil, err
	}

	role, err := domain.GetRole(domain.AuthRole(input.Role))

	if err != nil {
		return nil, fails.ROLE_NOT_FOUND
	}

	identityUser := domain.NewIdentityUser(
		input.Email,
		input.Password,
		domain.AuthRole(role),
	)

	signUpTransaction, err := as.authRepo.BeginTransaction()

	if err != nil {
		return nil, fails.InternalServerError()
	}

	if err := signUpTransaction.Create(identityUser); err != nil {
		return nil, fails.InternalServerError()
	}

	profile, err := as.createProfileHelper.CreateProfile(signUpTransaction, helpers.CreateProfileInput{
		AuthID: identityUser.ID,
		Email:  identityUser.Email,
		Role:   identityUser.GetRole(),
	})

	if err != nil {
		return nil, fails.InternalServerError()
	}

	accessToken, refreshToken, err := as.tokenService.CreateAuthenticationTokens(services.TokenPayload{
		UserID:     identityUser.ID,
		Email:      identityUser.Email,
		ProfileID:  profile.ID,
		Scope:      domain.GetPermissions(*identityUser),
		ProfileIds: make([]string, 0),
		Role:       string(identityUser.GetRole()),
	})

	if err != nil {
		return nil, fails.InternalServerError()
	}

	if err := signUpTransaction.Commit(); err != nil {
		return nil, fails.InternalServerError()
	}

	if err := as.emailService.Send(services.EmailInput{
		To:       identityUser.Email,
		Template: services.NewConfirmEmailTemplate(),
	}); err != nil {
		log.Println("Failed to send email of account confirmation")
	}

	return mappers.ToAuthResponse(
		identityUser,
		accessToken,
		refreshToken,
	), nil
}
