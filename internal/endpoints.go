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
	signUpUseCase *usecases.SignUpUseCase
	loginUseCase  *usecases.LoginUseCase
}

func NewAuthController(
	signUpUseCase *usecases.SignUpUseCase,
	loginUseCase *usecases.LoginUseCase,
) *AuthController {
	return &AuthController{
		signUpUseCase: signUpUseCase,
		loginUseCase:  loginUseCase,
	}
}

func (c *AuthController) Route(router fiber.Router) {
	authRoutes := router.Group(AUTH_CONTROLLER_NAME)

	authRoutes.Post("login", c.Login)
	authRoutes.Post("sign-up", c.SignUp)
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var loginRequest contracts.SignUpRequest

	if err := shared.ParseBodyAndValidate(ctx, &loginRequest); err != nil {
		return err
	}

	response, err := c.loginUseCase.Handle(usecases.LoginInput{
		Email:    loginRequest.Email,
		Password: loginRequest.Password,
	})

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *AuthController) SignUp(ctx *fiber.Ctx) error {
	var signUpRequest contracts.SignUpRequest

	if err := shared.ParseBodyAndValidate(ctx, &signUpRequest); err != nil {
		return err
	}

	response, err := c.signUpUseCase.Handle(usecases.SignUpInput{
		Email:    signUpRequest.Email,
		Password: signUpRequest.Password,
		Role:     signUpRequest.Role,
	})

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}
