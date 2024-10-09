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
	for i := 0; i < len(controllers); i++ {
		controllers[i].Route(hs.GetVersion())
	}
}

func (hs *HttpServer) Serve(port string) {
	log.Fatal(hs.Instance.Listen(fmt.Sprintf(":%s", port)))
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case *shared.Error:
		return shared.WriteProblemDetails(ctx, *e)
	case *fiber.Error:
		return shared.WriteProblemDetails(
			ctx,
			shared.Error{
				Title:  "Bad Input",
				Status: e.Code,
				Detail: e.Message,
			},
		)
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
