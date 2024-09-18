package internal

import "github.com/gofiber/fiber/v2"

const (
	CONTROLLER_NAME = "auth"
)

type AuthController struct {
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (c *AuthController) Route(router fiber.Router) {
	authRoutes := router.Group(CONTROLLER_NAME)

	authRoutes.Get("hello-world", c.HelloWorld)
}

func (c *AuthController) HelloWorld(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON("Hello World!, from the identity Server!")
}
