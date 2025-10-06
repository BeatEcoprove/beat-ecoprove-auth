package usecases

import (
	"testing"

	"github.com/BeatEcoprove/identityService/internal/usecases/utils"
	"github.com/BeatEcoprove/identityService/pkg/services"
	"github.com/BeatEcoprove/identityService/pkg/shared"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type (
	InputFaker struct {
		Email    string `faker:"email"`
		Password string
		Role     int `faker:"oneof: 0, 1"`
	}
)

func InitTest() {
	utils.TestSetup()

	Redis = new(utils.MockRedis)
	RabbitMq = new(utils.MockRabbitMq)

	AuthRepository = new(utils.MockAuthRepository)
	ProfileRepository = new(utils.MockProfileRepository)

	TokenService = services.NewTokenService(Redis)
	EmailService = services.NewEmailService(RabbitMq)
	PGService = services.NewPGService(Redis)
}

func SetupRabbitmq() {
	RabbitMq.On("PublishMessage", mock.Anything).Return(nil)
	RabbitMq.On("Close").Return(nil)
}

func SetupRedis() {
	Redis.On("GetAndDelValue", mock.Anything).Return("", nil)
	Redis.On("SetValue", mock.Anything, mock.Anything, mock.Anything).Return(nil)
}

const DefaultPassword = "Password1"

var (
	Redis    *utils.MockRedis
	RabbitMq *utils.MockRabbitMq

	AuthRepository    *utils.MockAuthRepository
	ProfileRepository *utils.MockProfileRepository

	TokenService services.ITokenService
	EmailService services.IEmailService
	PGService    services.IPGService
)

func generateFakeData(input any) {
	if err := faker.FakeData(input); err != nil {
		panic(err)
	}
}

func evaluateError(t *testing.T, shouldBe *shared.Error, err error) {
	var fail *shared.Error
	assert.ErrorAs(t, err, &fail)

	if fail != nil {
		assert.Equal(t, shouldBe.Id, fail.Id)
	}
}
