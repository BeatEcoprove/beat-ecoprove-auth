package usecases

import (
	"github.com/BeatEcoprove/identityService/internal/repositories"
	"github.com/BeatEcoprove/identityService/pkg/contracts"
	fails "github.com/BeatEcoprove/identityService/pkg/errors"
	"github.com/BeatEcoprove/identityService/pkg/services"
)

type (
	// input
	ResetPasswdInput struct {
		Email    string
		Code     string
		Password string
	}

	ResetPasswdUseCase struct {
		authRepo     repositories.IAuthRepository
		pgService    services.IPGService
		emailService services.IEmailService
	}
)

func NewResetPasswdUseCase(
	authRepo repositories.IAuthRepository,
	pgService services.IPGService,
	emailService services.IEmailService,
) *ResetPasswdUseCase {
	return &ResetPasswdUseCase{
		authRepo:     authRepo,
		pgService:    pgService,
		emailService: emailService,
	}
}

func (rpu *ResetPasswdUseCase) Handle(request ResetPasswdInput) (*contracts.GenericResponse, error) {
	identityUser, err := rpu.authRepo.GetUserByEmail(request.Email)

	if err != nil {
		return nil, fails.USER_NOT_FOUND
	}

	if err := rpu.pgService.ValidateCode(identityUser.ID, request.Code); err != nil {
		return nil, fails.CODE_NOT_VALID
	}

	if err := services.ValidatePassword(request.Password); err != nil {
		return nil, err
	}

	if err := identityUser.SetPassword(request.Password); err != nil {
		return nil, fails.InternalServerError()
	}

	if err := rpu.authRepo.Update(identityUser); err != nil {
		return nil, fails.InternalServerError()
	}

	return &contracts.GenericResponse{
		Message: "The password was changed with success.",
	}, nil
}
