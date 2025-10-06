package usecases

import (
	"testing"

	"github.com/BeatEcoprove/identityService/internal/domain"
	"github.com/BeatEcoprove/identityService/internal/usecases/utils"
	fails "github.com/BeatEcoprove/identityService/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type (
	SignUpInputFaker struct {
		Email    string `faker:"email"`
		Password string
		Role     int `faker:"oneof: 0, 1"`
	}
)

func Test_SignUp_UseCase(t *testing.T) {
	InitTest()
	SetupRedis()

	var sut *SignUpUseCase = NewSignUpUseCase(
		AuthRepository,
		ProfileRepository,
		TokenService,
		EmailService,
	)

	t.Run("Should not create an account if the email is already in use", func(t *testing.T) {
		// Arrange
		var input SignUpInputFaker = SignUpInputFaker{}
		generateFakeData(&input)
		input.Password = DefaultPassword

		// Act
		AuthRepository.On("ExistsUserWithEmail", input.Email).Return(true)

		sut.Handle(SignUpInput(input))
		_, err := sut.Handle(SignUpInput(input))

		// Assert
		evaluateError(t, fails.USER_ALREADY_EXISTS, err)
	})

	t.Run("Should create an account and return the correct response", func(t *testing.T) {
		// Arrange
		transRepo := new(utils.MockTransaction)

		var input SignUpInputFaker = SignUpInputFaker{}
		generateFakeData(&input)
		input.Password = DefaultPassword

		// Act
		AuthRepository.On("ExistsUserWithEmail", input.Email).Return(false)
		AuthRepository.On("BeginTransaction").Return(transRepo, nil)
		transRepo.MockRepositoryBase.On("Create", mock.Anything).Return(nil)
		transRepo.On("Commit", mock.Anything).Return(nil)

		response, err := sut.Handle(SignUpInput(input))

		// Assert
		role, _ := domain.GetRole(domain.AuthRole(input.Role))

		assert.Nil(t, err)
		assert.NotEmpty(t, response.AccessToken)
		assert.Greater(t, response.AccessTokenExpiration, 0)
		assert.NotEmpty(t, response.RefreshToken)
		assert.Greater(t, response.RefreshTokenExpiration, 0)
		assert.Equal(t, response.Details.Email, input.Email)
		assert.Equal(t, response.Details.Role, role)
	})
}
