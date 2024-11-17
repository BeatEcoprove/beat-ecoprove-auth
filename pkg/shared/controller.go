package shared

import (
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	Route(router fiber.Router)
}
