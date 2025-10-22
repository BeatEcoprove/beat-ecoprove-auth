package usecases

import "github.com/BeatEcoprove/identityService/internal/repositories"

type (
	CheckFieldInput struct {
		Email string
	}

	CheckFieldUseCase struct {
		authRepo repositories.IAuthRepository
	}
)

func NewCheckFieldUseCase(
	authRepo repositories.IAuthRepository,
) *CheckFieldUseCase {
	return &CheckFieldUseCase{
		authRepo: authRepo,
	}
}

func (cfu *CheckFieldUseCase) Handle(request CheckFieldInput) bool {
	return cfu.authRepo.ExistsUserWithEmail(request.Email)
}
