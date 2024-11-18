package main

import (
	"github.com/BeatEcoprove/identityService/config"
	"github.com/BeatEcoprove/identityService/internal"
	"github.com/BeatEcoprove/identityService/internal/adapters"
	"github.com/BeatEcoprove/identityService/internal/middlewares"
	"github.com/BeatEcoprove/identityService/internal/repositories"
	"github.com/BeatEcoprove/identityService/internal/usecases"
	"github.com/BeatEcoprove/identityService/pkg/services"
	"github.com/BeatEcoprove/identityService/pkg/shared"
)

const (
	API_VERSION = "2"
)

func generatePki() error {
	// Generate PKIS
	if !services.MustCreatePKI() {
		if err := services.GenerateServerPKI(); err != nil {
			return err
		}
	}

	// load pki
	publicKey, privateKey, err := services.ReadKeys()

	if err != nil {
		return err
	}

	err = services.LoadKeys(privateKey, publicKey)

	if err != nil {
		return err
	}

	return nil
}

func generateJWKS() error {
	if _, err := services.NewJWKS(); err != nil {
		return err
	}

	return nil
}

func main() {
	config.LoadEnv(config.DotEnv)
	env := config.GetCofig()

	if err := generatePki(); err != nil {
		panic(err)
	}

	if err := generateJWKS(); err != nil {
		panic(err)
	}

	// adapters
	db := adapters.GetDatabase()
	defer db.Close()

	redis := adapters.GetRedis()
	defer redis.Close()

	rabbitMQ, err := adapters.GetRabbitMqConnection()

	if err != nil {
		panic(err)
	}

	defer rabbitMQ.Close()

	app := adapters.NewHttpServer(API_VERSION)

	// repositories
	authRepository := repositories.NewAuthRepository(db)
	profileRepository := repositories.NewProfileRepository(db)

	// services
	tokenService := services.NewTokenService(redis)
	pgService := services.NewPGService(redis)
	emailService := services.NewEmailService(rabbitMQ)

	// midlewares
	authMiddleware := middlewares.NewAuthorizationMiddleware(authRepository, tokenService)

	// use cases
	signUpUseCase := usecases.NewSignUpUseCase(authRepository, profileRepository, tokenService, emailService)
	loginUseCase := usecases.NewLoginUseCase(authRepository, profileRepository, tokenService)
	attachProfileUseCase := usecases.NewAttachProfileUseCase(authRepository, profileRepository)
	refreshTokensUseCase := usecases.NewRefreshTokensUseCase(authRepository, profileRepository, tokenService)
	forgotPasswordUseCase := usecases.NewForgotPasswordUseCase(authRepository, pgService, emailService)
	resetPasswdUseCase := usecases.NewResetPasswdUseCase(authRepository, pgService, emailService)
	checkFieldUseCase := usecases.NewCheckFieldUseCase(authRepository)

	// controllers
	staticController := internal.NewStaticController()
	authController := internal.NewAuthController(
		signUpUseCase,
		loginUseCase,
		attachProfileUseCase,
		refreshTokensUseCase,
		forgotPasswordUseCase,
		resetPasswdUseCase,
		checkFieldUseCase,
		authMiddleware,
	)

	app.AddStaticController(staticController)
	app.AddControllers([]shared.Controller{
		authController,
	})

	app.Serve(env.BEAT_IDENTITY_SERVER)
}
