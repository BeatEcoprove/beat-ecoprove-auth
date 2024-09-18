package main

import (
	"github.com/BeatEcoprove/identityService/internal"
	"github.com/BeatEcoprove/identityService/internal/adapters"
	"github.com/BeatEcoprove/identityService/pkg/shared"
)

const (
	API_VERSION = "2"
	API_PORT    = "3000"
)

func main() {
	app := adapters.NewHttpServer(API_VERSION)

	authController := internal.NewAuthController()

	app.AddControllers([]shared.Controller{
		authController,
	})

	app.Serve(API_PORT)
}
