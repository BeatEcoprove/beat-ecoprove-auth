package internal

import (
	"github.com/BeatEcoprove/identityService/internal/usecases"
	"github.com/BeatEcoprove/identityService/pkg/contracts"
	"github.com/BeatEcoprove/identityService/pkg/shared"
	"github.com/gofiber/fiber/v2"
)

const (
	AUTH_CONTROLLER_NAME = "auth"
)

type AuthController struct {
	authUseCase *usecases.SignUpUseCase
}

func NewAuthController(
	authService *usecases.SignUpUseCase,
) *AuthController {
	return &AuthController{
		authUseCase: authService,
	}
}

func (c *AuthController) Route(router fiber.Router) {
	authRoutes := router.Group(AUTH_CONTROLLER_NAME)

	authRoutes.Post("sign-up", c.SignUp)
}

func (c *AuthController) SignUp(ctx *fiber.Ctx) error {
	var signUpRequest contracts.SignUpRequest

	if err := shared.ParseBodyAndValidate(ctx, &signUpRequest); err != nil {
		return err
	}

	response, err := c.authUseCase.Handle(usecases.SignUpInput{
		Email:    signUpRequest.Email,
		Password: signUpRequest.Password,
		Role:     signUpRequest.Role,
	})

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}
