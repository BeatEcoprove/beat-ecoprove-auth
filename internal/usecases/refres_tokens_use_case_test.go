package usecases

import (
	"errors"
	"testing"

	"github.com/BeatEcoprove/identityService/internal/domain"
	fails "github.com/BeatEcoprove/identityService/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func getRefreshTokenInput(includeProfile bool) *RefreshTokensInput {
	var profileId string = ""

	if includeProfile {
		profileId = uuid.New().String()
	}

	return &RefreshTokensInput{
		AuthId:    uuid.New().String(),
		ProfileId: profileId,
	}
}

func getTestProfiles(authId string, length int) []domain.Profile {
	var attachedProfiles []domain.Profile = make([]domain.Profile, length)

	attachedProfiles = append(attachedProfiles, *domain.NewProfile(authId, domain.Main))
	for i := 0; i < length; i++ {
		attachedProfiles = append(attachedProfiles, *domain.NewProfile(authId, domain.Sub))
	}

	return attachedProfiles
}

func Test_Refresh_Token_UseCase(t *testing.T) {
	InitTest()
	SetupRedis()

	var sut *RefreshTokensUseCase = NewRefreshTokensUseCase(
		AuthRepository,
		ProfileRepository,
		TokenService,
	)

	t.Run("Should generate tokens for a single main profile", func(t *testing.T) {
		t.Run("Should not refresh tokens if the user does not exists", func(t *testing.T) {
			var input *RefreshTokensInput = getRefreshTokenInput(true)

			// Act
			AuthRepository.On("Get", input.AuthId).Return(&domain.IdentityUser{}, errors.ErrUnsupported)

			_, err := sut.Handle(*input)

			// Assert
			evaluateError(t, fails.USER_NOT_FOUND, err)
		})

		t.Run("Should fail when does not have access to attached profile", func(t *testing.T) {
			var input *RefreshTokensInput = getRefreshTokenInput(true)

			var data LoginInputFaker = LoginInputFaker{}
			generateFakeData(&data)
			data.Password = DefaultPassword

			identityUser, err := getIdentityUser(
				data.Email,
				data.Password,
				int(domain.Client),
			)

			assert.Equal(t, err, nil)

			// Act
			AuthRepository.On("Get", input.AuthId).Return(identityUser, nil)
			ProfileRepository.On("IsProfileFromUserId", input.AuthId, input.ProfileId).Return(false)

			_, err = sut.Handle(*input)

			// Assert
			evaluateError(t, fails.PROFILE_DOES_NOT_BELONG_TO_USER, err)
		})

		t.Run("Should fail if it can't find the profile attached to the account", func(t *testing.T) {
			var input *RefreshTokensInput = getRefreshTokenInput(true)

			var data LoginInputFaker = LoginInputFaker{}
			generateFakeData(&data)
			data.Password = DefaultPassword

			identityUser, err := getIdentityUser(
				data.Email,
				data.Password,
				int(domain.Client),
			)

			assert.Equal(t, err, nil)

			// Act
			AuthRepository.On("Get", input.AuthId).Return(identityUser, nil)
			ProfileRepository.On("IsProfileFromUserId", input.AuthId, input.ProfileId).Return(true)
			ProfileRepository.On("Get", input.ProfileId).Return(&domain.Profile{}, errors.ErrUnsupported)

			_, err = sut.Handle(*input)

			// Assert
			evaluateError(t, fails.PROFILE_NOT_FOUND, err)
		})

		t.Run("Should generate a token pair for authorization", func(t *testing.T) {
			var input *RefreshTokensInput = getRefreshTokenInput(true)

			var data LoginInputFaker = LoginInputFaker{}
			generateFakeData(&data)
			data.Password = DefaultPassword

			identityUser, err := getIdentityUser(
				data.Email,
				data.Password,
				int(domain.Client),
			)

			assert.Equal(t, err, nil)

			testOneProfile := domain.NewProfile(identityUser.GetId(), domain.Main)

			// Act
			AuthRepository.On("Get", input.AuthId).Return(identityUser, nil)
			ProfileRepository.On("IsProfileFromUserId", input.AuthId, input.ProfileId).Return(true)
			ProfileRepository.On("Get", input.ProfileId).Return(testOneProfile, nil)

			response, err := sut.Handle(*input)

			// Assert
			role, _ := domain.GetRole(identityUser.Role)

			assert.Nil(t, err)
			assert.NotEmpty(t, response.AccessToken)
			assert.Greater(t, response.AccessTokenExpiration, 0)
			assert.NotEmpty(t, response.RefreshToken)
			assert.Greater(t, response.RefreshTokenExpiration, 0)
			assert.Equal(t, response.Details.Email, identityUser.Email)
			assert.Equal(t, response.Details.Role, role)
		})

		t.Run("Should generate a token pair for authorization, by multiple profiles", func(t *testing.T) {
			var input *RefreshTokensInput = getRefreshTokenInput(false)

			var data LoginInputFaker = LoginInputFaker{}
			generateFakeData(&data)
			data.Password = DefaultPassword

			identityUser, err := getIdentityUser(
				data.Email,
				data.Password,
				int(domain.Client),
			)

			assert.Equal(t, err, nil)

			testProfiles := getTestProfiles(identityUser.GetId(), 10)

			// Act
			AuthRepository.On("Get", input.AuthId).Return(identityUser, nil)
			ProfileRepository.On("IsProfileFromUserId", input.AuthId, input.ProfileId).Return(true)
			ProfileRepository.On("GetAttachProfiles", input.AuthId).Return(testProfiles, nil)

			response, err := sut.Handle(*input)

			// Assert
			role, _ := domain.GetRole(identityUser.Role)

			assert.Nil(t, err)
			assert.NotEmpty(t, response.AccessToken)
			assert.Greater(t, response.AccessTokenExpiration, 0)
			assert.NotEmpty(t, response.RefreshToken)
			assert.Greater(t, response.RefreshTokenExpiration, 0)
			assert.Equal(t, response.Details.Email, identityUser.Email)
			assert.Equal(t, response.Details.Role, role)
		})
	})
}
