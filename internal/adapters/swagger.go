package adapters

import (
	"fmt"

	docs "github.com/BeatEcoprove/identityService/docs"
	"github.com/gofiber/swagger"
)

func UseSwagger(app *HttpServer, port uint16) {
	docs.SwaggerInfo.Title = "Identity Microservice"
	docs.SwaggerInfo.Description = "API documentation for Identity Microservice"
	docs.SwaggerInfo.Version = fmt.Sprintf("%s.0", app.version)
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", port)
	docs.SwaggerInfo.BasePath = fmt.Sprintf("/api/v%s/auth", app.version)
	docs.SwaggerInfo.Schemes = []string{"http"}

	app.Instance.Get("/swagger/*", swagger.HandlerDefault)
	app.Instance.Get("/swagger/*", swagger.New(swagger.Config{}))
}
