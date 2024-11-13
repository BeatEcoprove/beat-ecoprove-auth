package usecases

import (
	"encoding/pem"
	"os"
	"testing"

	"github.com/BeatEcoprove/identityService/config"
	"github.com/BeatEcoprove/identityService/internal/adapters"
	"github.com/BeatEcoprove/identityService/internal/domain"
	"github.com/BeatEcoprove/identityService/internal/repositories"
	interfaces "github.com/BeatEcoprove/identityService/pkg/adapters"
	fails "github.com/BeatEcoprove/identityService/pkg/errors"
	"github.com/BeatEcoprove/identityService/pkg/services"
	"github.com/BeatEcoprove/identityService/pkg/shared"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
)

type (
	SignUpInputFaker struct {
		Email    string `faker:"email"`
		Password string
		Role     int `faker:"oneof: 0, 1"`
	}
)

var (
	db interfaces.Database

	authRepository    repositories.IAuthRepository
	profileRepository repositories.IProfileRepository

	sut *SignUpUseCase
)

func TestMain(m *testing.M) {
	// before
	dotEnvPath, err := config.GetDotEnvPath()

	if err != nil {
		panic(err)
	}

	config.LoadEnv(dotEnvPath)

	// load pki
	keys, err := services.CreatePKI()

	if err != nil {
		panic(err)
	}

	privateKey := pem.EncodeToMemory(keys.PrivateKey)
	publicKey := pem.EncodeToMemory(keys.PublicKey)

	err = services.LoadKeys(publicKey, privateKey)

	if err != nil {
		panic(err)
	}

	if _, err = services.CreateJWKS(); err != nil {
		panic(err)
	}

	postgres, connectionString, err := config.SpawnPostgresqlContainer()

	if err != nil {
		panic(err)
	}

	if err := config.Migrate(connectionString, false); err != nil {
		panic(err)
	}

	db = adapters.GetDatabaseWithConnectionString(connectionString)

	exitCode := m.Run()

	// After
	config.KillContainer(postgres)

	os.Exit(exitCode)
}

func Test_SignUp_UseCase(t *testing.T) {
	// Arrange
	authRepository = repositories.NewAuthRepository(db)
	profileRepository = repositories.NewProfileRepository(db)
	sut = NewSignUpUseCase(authRepository, profileRepository)

	t.Run("Should not create an account if the email is already in use", func(t *testing.T) {
		// Arrange
		input := SignUpInputFaker{}

		if err := faker.FakeData(&input); err != nil {
			panic(err)
		}
		input.Password = "Password1"

		// Act
		sut.Handle(SignUpInput(input))
		_, err := sut.Handle(SignUpInput(input))

		// Assert
		var fail *shared.Error
		assert.ErrorAs(t, err, &fail)

		if fail != nil {
			assert.Equal(t, fails.USER_ALREADY_EXISTS.Id, fail.Id)
		}
	})

	t.Run("Should create an account and return the correct response", func(t *testing.T) {
		// Arrange
		input := SignUpInputFaker{}

		if err := faker.FakeData(&input); err != nil {
			panic(err)
		}
		input.Password = "Password1"

		// Act
		response, err := sut.Handle(SignUpInput(input))

		// Assert
		role, _ := domain.GetRole(domain.Role(input.Role))

		assert.Nil(t, err)
		assert.NotEmpty(t, response.AccessToken)
		assert.Greater(t, response.AccessTokenExpiration, 0)
		assert.NotEmpty(t, response.RefreshToken)
		assert.Greater(t, response.RefreshTokenExpiration, 0)
		assert.Equal(t, response.Details.Email, input.Email)
		assert.NotEmpty(t, response.Details.ProfileId)
		assert.Equal(t, response.Details.Role, role)
	})
}
