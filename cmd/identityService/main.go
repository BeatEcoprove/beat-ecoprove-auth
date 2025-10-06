package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BeatEcoprove/identityService/config"
	"github.com/BeatEcoprove/identityService/internal"
	"github.com/BeatEcoprove/identityService/pkg/services"
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

func initPKIandJWKS() error {
	if err := generatePki(); err != nil {
		return err
	}

	if err := generateJWKS(); err != nil {
		return err
	}

	return nil
}

func exitGracefully() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("🧹 Shutting down gracefully...")
}

//	@termsOfService				http://swagger.io/terms/
//	@contact.name				API Support
//	@contact.email				fiber@swagger.io
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description Enter the token with the `Bearer: ` prefix, e.g. "Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IjQ5MjRhNmEx..."
//
// @ Schemas http
func main() {
	config.LoadEnv(config.DotEnv)

	err := initPKIandJWKS()

	if err != nil {
		panic(err)
	}

	app, err := internal.NewApp()

	if err != nil {
		panic(err)
	}

	app.ApplyConsumer()
	app.ApplyHttpServer()

	app.Serve()
	defer app.Close()

	exitGracefully()
}
