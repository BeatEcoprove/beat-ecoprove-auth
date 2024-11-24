package usecases

import (
	"errors"
	"strings"
	"testing"

	"github.com/BeatEcoprove/identityService/internal/domain"
	fails "github.com/BeatEcoprove/identityService/pkg/errors"
	"github.com/BeatEcoprove/identityService/pkg/services"
	"github.com/stretchr/testify/assert"
)

type (
	ResetPasswdInputFaker struct {
		Email    string
		Code     string
		Password string
	}
)

func shuffling(s string) string {
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return strings.ReplaceAll(string(runes), services.Delimiter, "\\\"#$")
}

func generateFakeCode() (string, error) {
	code, err := services.GenerateCode()

	if err != nil {
		return "", err
	}

	return shuffling(code), nil
}

func Test_Reset_Password_UseCase(t *testing.T) {
	InitTest()

	var sut *ResetPasswdUseCase = NewResetPasswdUseCase(
		AuthRepository,
		PGService,
		EmailService,
	)

	t.Run("Should not reset password if the user can't be found", func(t *testing.T) {
		// Assert
		var input ResetPasswdInputFaker = ResetPasswdInputFaker{}
		generateFakeData(&input)

		// Act
		AuthRepository.On("GetUserByEmail", input.Email).Return(&domain.IdentityUser{}, errors.ErrUnsupported)
		_, err := sut.Handle(ResetPasswdInput(input))

		// Assert
		evaluateError(t, fails.USER_NOT_FOUND, err)
	})

	t.Run("Should not reset password if the code its not valid", func(t *testing.T) {
		// Assert
		var input ResetPasswdInputFaker = ResetPasswdInputFaker{}
		generateFakeData(&input)

		identityIdenty, err := getIdentityUser(
			input.Email,
			DefaultPassword,
			int(domain.Main),
		)

		assert.Equal(t, err, nil)

		// Act
		code, err := generateFakeCode()
		assert.Equal(t, err, nil)
		input.Code = code

		AuthRepository.On("GetUserByEmail", input.Email).Return(identityIdenty, errors.ErrUnsupported)
		_, err = sut.Handle(ResetPasswdInput(input))

		// Assert
		evaluateError(t, fails.USER_NOT_FOUND, err)
	})

}
