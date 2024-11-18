package internal

import (
	"github.com/BeatEcoprove/identityService/internal/middlewares"
	"github.com/BeatEcoprove/identityService/internal/usecases"
	"github.com/BeatEcoprove/identityService/pkg/contracts"
	fails "github.com/BeatEcoprove/identityService/pkg/errors"
	"github.com/BeatEcoprove/identityService/pkg/shared"
	"github.com/gofiber/fiber/v2"
)

const (
	AUTH_CONTROLLER_NAME = "auth"
)

type AuthController struct {
	signUpUseCase         *usecases.SignUpUseCase
	loginUseCase          *usecases.LoginUseCase
	attachProfileUseCase  *usecases.AttachProfileUseCase
	refreshTokensUseCase  *usecases.RefreshTokensUseCase
	forgotPasswordUseCase *usecases.ForgotPasswordUseCase
	resetPasswdUseCase    *usecases.ResetPasswdUseCase

	authMiddleware *middlewares.AuthorizationMiddleware
}

func NewAuthController(
	signUpUseCase *usecases.SignUpUseCase,
	loginUseCase *usecases.LoginUseCase,
	attachProfileUseCase *usecases.AttachProfileUseCase,
	refreshTokensUseCase *usecases.RefreshTokensUseCase,
	forgotPasswordUseCase *usecases.ForgotPasswordUseCase,
	resetPasswdUseCase *usecases.ResetPasswdUseCase,
	authMiddleware *middlewares.AuthorizationMiddleware,
) *AuthController {
	return &AuthController{
		signUpUseCase:         signUpUseCase,
		loginUseCase:          loginUseCase,
		attachProfileUseCase:  attachProfileUseCase,
		refreshTokensUseCase:  refreshTokensUseCase,
		forgotPasswordUseCase: forgotPasswordUseCase,
		resetPasswdUseCase:    resetPasswdUseCase,
		authMiddleware:        authMiddleware,
	}
}

func (c *AuthController) Route(router fiber.Router) {
	authRoutes := router.Group(AUTH_CONTROLLER_NAME)

	authRoutes.Post("profile", c.authMiddleware.AccessTokenHandler, c.AttachProfile)
	authRoutes.Get("token", c.authMiddleware.AccessTokenHandler, c.Token)
	authRoutes.Get("refresh-token", c.authMiddleware.RefreshTokenHandler, c.RefreshTokens)

	authRoutes.Post("reset-password", c.ResetPassword)
	authRoutes.Post("forgot-password", c.ForgotPassword)

	authRoutes.Post("login", c.Login)
	authRoutes.Post("sign-up", c.SignUp)
}

func (c *AuthController) AttachProfile(ctx *fiber.Ctx) error {
	var attachProfileRequest contracts.AttachProfileRequest

	authId, err := middlewares.GetUserId(ctx)

	if err != nil {
		return err
	}

	if err := shared.ParseBodyAndValidate(ctx, &attachProfileRequest); err != nil {
		return err
	}

	response, err := c.attachProfileUseCase.Handle(usecases.AttachProfileInput{
		AuthId:           authId,
		ProfileGrantType: attachProfileRequest.ProfileGrantType,
	})

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (c *AuthController) Token(ctx *fiber.Ctx) error {
	_, claims, err := middlewares.GetClaims(ctx)

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(&contracts.AccountResponse{
		UserId:     claims.Subject,
		Email:      claims.Email,
		ProfileId:  claims.ProfileId,
		ProfileIds: claims.ProfileIds,
		Role:       claims.Role,
	})
}

func (c *AuthController) RefreshTokens(ctx *fiber.Ctx) error {
	profileId := ctx.Query("profile_id", "")

	authId, err := middlewares.GetUserId(ctx)

	if err != nil {
		return err
	}

	if err := shared.Validate(&contracts.RefreshTokensRequest{ProfileId: profileId}); err != nil {
		return fails.BAD_UUID
	}

	response, err := c.refreshTokensUseCase.Handle(usecases.RefreshTokensInput{
		AuthId:    authId,
		ProfileId: profileId,
	})

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *AuthController) ResetPassword(ctx *fiber.Ctx) error {
	var resetPasswordRequest contracts.ResetPasswordRequest

	if err := shared.ParseBodyAndValidate(ctx, &resetPasswordRequest); err != nil {
		return err
	}

	response, err := c.resetPasswdUseCase.Handle(usecases.ResetPasswdInput{
		Email:    resetPasswordRequest.Email,
		Code:     resetPasswordRequest.Code,
		Password: resetPasswordRequest.Password,
	})

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *AuthController) ForgotPassword(ctx *fiber.Ctx) error {
	var forgotPasswordRequest contracts.ForgotPasswordRequest

	if err := shared.ParseBodyAndValidate(ctx, &forgotPasswordRequest); err != nil {
		return err
	}

	response, err := c.forgotPasswordUseCase.Handle(usecases.ForgotPasswordInput{
		Email: forgotPasswordRequest.Email,
	})

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
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
