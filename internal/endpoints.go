package internal

import (
	"errors"

	"github.com/BeatEcoprove/identityService/internal/services"
	"github.com/BeatEcoprove/identityService/pkg/contracts"
	"github.com/gofiber/fiber/v2"
)

const (
	CONTROLLER_NAME = "auth"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Route(router fiber.Router) {
	authRoutes := router.Group(CONTROLLER_NAME)

	authRoutes.Post("sign-up", c.SignUp)
}

func (c *AuthController) SignUp(ctx *fiber.Ctx) error {
	var signUpRequest contracts.SignUpRequest

	if err := ctx.BodyParser(&signUpRequest); err != nil {
		return err
	}

	response, err := c.authService.SignUp(services.SignUpInput{
		Email:    signUpRequest.Email,
		Password: signUpRequest.Password,
		Role:     signUpRequest.Role,
	})

	if err != nil {
		return errors.New("an error has occorred in creating an user")
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
