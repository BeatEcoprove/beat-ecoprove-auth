package usecases

import (
	"errors"
	"testing"

	"github.com/BeatEcoprove/identityService/internal/domain"
	fails "github.com/BeatEcoprove/identityService/pkg/errors"
	"github.com/BeatEcoprove/identityService/pkg/shared"
	"github.com/stretchr/testify/assert"
)

type (
	ForgotPasswordInputFaker struct {
		Email string `faker:"email"`
	}
)

func Test_ForgotPassword_UseCase(t *testing.T) {
	InitTest()
	SetupRedis()
	SetupRabbitmq()

	var sut *ForgotPasswordUseCase = NewForgotPasswordUseCase(
		AuthRepository,
		PGService,
		EmailService,
	)

	t.Run("Shoul fail when fetching the user", func(t *testing.T) {
		var input ForgotPasswordInput = ForgotPasswordInput{}
		generateFakeData(&input)

		// Act
		AuthRepository.On("GetUserByEmail", input.Email).Return(&domain.IdentityUser{}, errors.ErrUnsupported)
		_, err := sut.Handle(input)

		// Assert
		var fail *shared.Error
		assert.ErrorAs(t, err, &fail)

		if fail != nil {
			assert.Equal(t, fails.USER_NOT_FOUND.Id, fail.Id)
		}
	})

	t.Run("Should send a email with the forgot code", func(t *testing.T) {
		var input ForgotPasswordInput = ForgotPasswordInput{}
		generateFakeData(&input)

		// Act
		identityIdenty, err := getIdentityUser(
			input.Email,
			DefaultPassword,
			int(domain.Main),
		)

		assert.Equal(t, err, nil)

		// Act
		AuthRepository.On("GetUserByEmail", identityIdenty.Email).Return(identityIdenty, nil)
		response, err := sut.Handle(input)

		// Assert
		assert.Nil(t, err)
		assert.NotEmpty(t, response.Message)
	})

	t.Run("Should generate and encrypt a valid code", func(t *testing.T) {
		var input ForgotPasswordInput = ForgotPasswordInput{}
		generateFakeData(&input)

		// Act
		identityIdenty, err := getIdentityUser(
			input.Email,
			DefaultPassword,
			int(domain.Main),
		)

		assert.Equal(t, err, nil)

		// Act
		AuthRepository.On("GetUserByEmail", identityIdenty.Email).Return(identityIdenty, nil)
		sut.Handle(input)

		// Assert
		lastEmail, err := EmailService.Last()
		assert.Nil(t, err)
		assert.NotNil(t, lastEmail)
	})
}
