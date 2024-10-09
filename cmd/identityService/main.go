package main

import (
	"github.com/BeatEcoprove/identityService/config"
	"github.com/BeatEcoprove/identityService/internal"
	"github.com/BeatEcoprove/identityService/internal/adapters"
	"github.com/BeatEcoprove/identityService/internal/repositories"
	"github.com/BeatEcoprove/identityService/internal/services"
	"github.com/BeatEcoprove/identityService/pkg/shared"
)

const (
	API_VERSION = "2"
)

func main() {
	config.LoadEnv(".env")
	env := config.GetCofig()

	// adapters
	db := adapters.GetDatabase()
	app := adapters.NewHttpServer(API_VERSION)

	// repositories
	authRepository := repositories.NewAuthRepository(db)
	profileRepository := repositories.NewProfileRepository(db)

	// services
	authService := services.NewAuthService(authRepository, profileRepository)

	// controllers
	authController := internal.NewAuthController(authService)

	app.AddControllers([]shared.Controller{
		authController,
	})

	app.Serve(env.BEAT_IDENTITY_SERVER)
}
