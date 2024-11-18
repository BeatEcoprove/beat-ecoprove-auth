package usecases

import (
	"github.com/BeatEcoprove/identityService/internal/repositories"
	"github.com/BeatEcoprove/identityService/pkg/contracts"
	fails "github.com/BeatEcoprove/identityService/pkg/errors"
	"github.com/BeatEcoprove/identityService/pkg/services"
)

type (
	// input
	ForgotPasswordInput struct {
		Email    string
		Password string
	}

	ForgotPasswordUseCase struct {
		authRepo     repositories.IAuthRepository
		pgService    services.IPGService
		emailService services.IEmailService
	}
)

func NewForgotPasswordUseCase(
	authRepo repositories.IAuthRepository,
	pgService services.IPGService,
	emailService services.IEmailService,
) *ForgotPasswordUseCase {
	return &ForgotPasswordUseCase{
		authRepo:     authRepo,
		pgService:    pgService,
		emailService: emailService,
	}
}

func (fpu *ForgotPasswordUseCase) Handle(request ForgotPasswordInput) (*contracts.GenericResponse, error) {
	identityUser, err := fpu.authRepo.GetUserByEmail(request.Email)

	if err != nil {
		return nil, fails.USER_NOT_FOUND
	}

	genCode, err := fpu.pgService.CreateAndStoreCode(identityUser.ID)

	if err != nil {
		return nil, fails.InternalServerError()
	}

	fpu.emailService.Send(services.EmailInput{
		To:         identityUser.Email,
		Subject:    string(*genCode),
		TemplateId: services.Default,
	})

	return &contracts.GenericResponse{
		Message: "It was sent an email with the code to renew your password to your email",
	}, nil
}
