package adapters

import (
	"fmt"
	"log"

	"github.com/BeatEcoprove/identityService/pkg/shared"
	"github.com/gofiber/fiber/v2"
)

type HttpServer struct {
	Instance *fiber.App
	version  string
}

func NewHttpServer(version string, middewares ...fiber.Handler) *HttpServer {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	return &HttpServer{
		Instance: app,
		version:  version,
	}
}

func (hs *HttpServer) GetVersion() fiber.Router {
	return hs.Instance.Group(fmt.Sprintf("api/v%s", hs.version))
}

func (hs *HttpServer) AddControllers(controllers []shared.Controller) {
	for i := range controllers {
		controllers[i].Route(hs.GetVersion())
	}
}

func (hs *HttpServer) AddStaticController(controller shared.Controller) {
	controller.Route(hs.Instance)
}

func (hs *HttpServer) Serve(port uint16) {
	log.Fatal(hs.Instance.Listen(fmt.Sprintf(":%d", port)))
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case *shared.Error:
		return shared.WriteProblemDetails(ctx, *e)
	case *shared.ValidationError:
		return shared.WriteProblemDetailsValidation(ctx, *e)
	default:
		return shared.WriteProblemDetails(
			ctx,
			shared.Error{
				Title:  "Internal Server Error",
				Status: fiber.StatusInternalServerError,
				Detail: err.Error(),
			},
		)
	}
}
