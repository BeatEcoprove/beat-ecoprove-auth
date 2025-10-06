package internal

import (
	"github.com/BeatEcoprove/identityService/config"
	"github.com/BeatEcoprove/identityService/internal/adapters"
	"github.com/BeatEcoprove/identityService/internal/domain/events"
	"github.com/BeatEcoprove/identityService/internal/domain/handlers"
	"github.com/BeatEcoprove/identityService/internal/middlewares"
	"github.com/BeatEcoprove/identityService/internal/repositories"
	"github.com/BeatEcoprove/identityService/internal/usecases"
	interfaces "github.com/BeatEcoprove/identityService/pkg/adapters"
	"github.com/BeatEcoprove/identityService/pkg/services"
	"github.com/BeatEcoprove/identityService/pkg/shared"
)

const (
	API_VERSION = "2"
)

type App struct {
	Db            interfaces.Database
	Redis         interfaces.Redis
	Publisher     interfaces.Broker
	Consumer      interfaces.Consumer
	HttpServer    *adapters.HttpServer
	Controllers   *Controllers
	Middlewares   *middlewares.Middlewares
	Repositories  *repositories.Repositories
	UseCases      *usecases.UseCases
	Services      *services.Services
	EventHandlers *handlers.EventHandlers
}

type Controllers struct {
	Static *StaticController
	Auth   *AuthController
}

func NewApp() (*App, error) {
	db := adapters.GetDatabase()
	redis := adapters.GetRedis()

	kafkaPub, kafkaSub, err := initKafka()

	if err != nil {
		return nil, err
	}

	repos := &repositories.Repositories{
		Auth:       repositories.NewAuthRepository(db),
		Profile:    repositories.NewProfileRepository(db),
		MemberChat: repositories.NewMemberChatRepository(db),
	}

	services := &services.Services{
		Token: services.NewTokenService(redis),
		PG:    services.NewPGService(redis),
		Email: services.NewEmailService(kafkaPub),
	}

	usecases := &usecases.UseCases{
		Sign:           usecases.NewSignUpUseCase(repos.Auth, repos.Profile, services.Token, services.Email, kafkaPub),
		Login:          usecases.NewLoginUseCase(repos.Auth, repos.Profile, services.Token),
		AttachProfile:  usecases.NewAttachProfileUseCase(repos.Auth, repos.Profile),
		RefreshTokens:  usecases.NewRefreshTokensUseCase(repos.Auth, repos.Profile, services.Token),
		ForgotPassword: usecases.NewForgotPasswordUseCase(repos.Auth, services.PG, services.Email),
		ResetPassword:  usecases.NewResetPasswdUseCase(repos.Auth, services.PG, services.Email),
		CheckFields:    usecases.NewCheckFieldUseCase(repos.Auth),
	}

	middlewares := &middlewares.Middlewares{
		Authorization: middlewares.NewAuthorizationMiddleware(repos.Auth, services.Token),
	}

	controllers := &Controllers{
		Static: NewStaticController(),
		Auth: NewAuthController(
			usecases.Sign,
			usecases.Login,
			usecases.AttachProfile,
			usecases.RefreshTokens,
			usecases.ForgotPassword,
			usecases.ResetPassword,
			usecases.CheckFields,
			middlewares.Authorization,
		),
	}

	eventHandlers := &handlers.EventHandlers{
		GroupCreated:   handlers.NewGroupCreatedHandler(repos.MemberChat, repos.Auth),
		InviteAccepted: handlers.NewInviteAcceptedHandler(repos.MemberChat, repos.Auth),
	}

	httpServer := adapters.NewHttpServer(API_VERSION)

	return &App{
		Db:            db,
		Redis:         redis,
		Publisher:     kafkaPub,
		Consumer:      kafkaSub,
		HttpServer:    httpServer,
		Controllers:   controllers,
		Middlewares:   middlewares,
		Repositories:  repos,
		UseCases:      usecases,
		Services:      services,
		EventHandlers: eventHandlers,
	}, nil
}

func (app *App) ApplyHttpServer() {
	env := config.GetConfig()
	adapters.UseSwagger(app.HttpServer, env.BEAT_IDENTITY_SERVER)

	app.HttpServer.AddStaticController(app.Controllers.Static)
	app.HttpServer.AddControllers([]shared.Controller{
		app.Controllers.Auth,
	})
}

func (app *App) ApplyConsumer() {
	app.Consumer.Register(app.EventHandlers.GroupCreated, &events.GroupCreatedEvent{})
	app.Consumer.Register(app.EventHandlers.InviteAccepted, &events.InviteAcceptedEvent{})
}

func (app *App) Serve() {
	env := config.GetConfig()

	go app.HttpServer.Serve(env.BEAT_IDENTITY_SERVER)
	go app.Consumer.Consume()
}

func initKafka() (*adapters.KafkaPublisher, *adapters.KafkaConsumer, error) {
	kafkaPublisher, err := adapters.GetKafkaPublisher()

	if err != nil {
		panic(err)
	}

	kafkaConsumer, err := adapters.GetKafkaConsumer()

	if err != nil {
		panic(err)
	}

	return kafkaPublisher, kafkaConsumer, nil
}

func (app *App) Close() error {
	app.Db.Close()
	app.Redis.Close()
	app.Publisher.Close()
	app.Consumer.Close()

	return nil
}
