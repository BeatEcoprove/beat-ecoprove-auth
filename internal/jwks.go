package internal

import (
	"fmt"

	"github.com/BeatEcoprove/identityService/pkg/services"
	"github.com/gofiber/fiber/v2"
)

type StaticController struct {
}

func NewStaticController() *StaticController {
	return &StaticController{}
}

func (c *StaticController) Route(router fiber.Router) {
	router.Get(".well-known/jwks.json", c.JWKS)
}

func (c *StaticController) JWKS(ctx *fiber.Ctx) error {
	jwks, err := services.NewJWKS()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error generating JWKS: %v", err))
	}

	return ctx.Status(fiber.StatusOK).JSON(jwks)
}
