package usecases

import (
	"math/rand"
	"testing"

	"github.com/BeatEcoprove/identityService/internal/domain"
	fails "github.com/BeatEcoprove/identityService/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type (
	LoginInputFaker struct {
		Email    string `faker:"email"`
		Password string
	}
)

func getIdentityUser(
	email string,
	password string,
	role int,
) (*domain.IdentityUser, error) {
	identityUser := &domain.IdentityUser{
		Email: email,
		Role:  domain.Role(role),
	}

	err := identityUser.SetPassword(password)
	return identityUser, err
}

func getProfilesById(identityId string, length int) []domain.Profile {
	attachedProfiles := make([]domain.Profile, length)

	for i := 0; i < length; i++ {
		attachedProfiles = append(attachedProfiles, *domain.NewProfile(
			identityId,
			domain.GrantType(rand.Intn(2)),
		))
	}

	return attachedProfiles
}

func Test_LogIn_UseCase(t *testing.T) {
	InitTest()
	SetupRedis()

	var sut *LoginUseCase = NewLoginUseCase(
		AuthRepository,
		ProfileRepository,
		TokenService,
	)

	t.Run("Should fail to login when user does not exists", func(t *testing.T) {
		var input LoginInputFaker = LoginInputFaker{}
		generateFakeData(&input)

		// Act
		AuthRepository.On("ExistsUserWithEmail", input.Email).Return(false)

		sut.Handle(LoginInput(input))
		_, err := sut.Handle(LoginInput(input))

		// Assert
		evaluateError(t, fails.USER_AUTH_FAILED, err)
	})

	t.Run("Should fail when password is not valid", func(t *testing.T) {
		var input LoginInputFaker = LoginInputFaker{}
		generateFakeData(&input)

		// Act
		AuthRepository.On("ExistsUserWithEmail", input.Email).Return(true)

		sut.Handle(LoginInput(input))
		_, err := sut.Handle(LoginInput(input))

		// Assert
		evaluateError(t, fails.USER_AUTH_FAILED, err)
	})

	t.Run("Should Login User", func(t *testing.T) {
		var input LoginInputFaker = LoginInputFaker{}
		generateFakeData(&input)
		input.Password = DefaultPassword

		identityUser, err := getIdentityUser(
			input.Email,
			input.Password,
			int(domain.Client),
		)

		assert.Equal(t, err, nil)

		profiles := getProfilesById(identityUser.GetId(), 4)

		// Act
		AuthRepository.On("ExistsUserWithEmail", input.Email).Return(true)
		AuthRepository.On("GetUserByEmail", input.Email).Return(identityUser, nil)
		ProfileRepository.On("GetAttachProfiles", identityUser.GetId()).Return(profiles, nil)

		sut.Handle(LoginInput(input))
		response, err := sut.Handle(LoginInput(input))

		// Assert
		role, _ := domain.GetRole(domain.Role(identityUser.Role))

		assert.Nil(t, err)
		assert.NotEmpty(t, response.AccessToken)
		assert.Greater(t, response.AccessTokenExpiration, 0)
		assert.NotEmpty(t, response.RefreshToken)
		assert.Greater(t, response.RefreshTokenExpiration, 0)
		assert.Equal(t, response.Details.Email, input.Email)
		assert.Equal(t, response.Details.Role, role)
	})
}
