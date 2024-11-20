package internal

import (
	"strconv"

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
	checkFieldUseCase     *usecases.CheckFieldUseCase

	authMiddleware *middlewares.AuthorizationMiddleware
}

func NewAuthController(
	signUpUseCase *usecases.SignUpUseCase,
	loginUseCase *usecases.LoginUseCase,
	attachProfileUseCase *usecases.AttachProfileUseCase,
	refreshTokensUseCase *usecases.RefreshTokensUseCase,
	forgotPasswordUseCase *usecases.ForgotPasswordUseCase,
	resetPasswdUseCase *usecases.ResetPasswdUseCase,
	checkFieldUseCase *usecases.CheckFieldUseCase,
	authMiddleware *middlewares.AuthorizationMiddleware,
) *AuthController {
	return &AuthController{
		signUpUseCase:         signUpUseCase,
		loginUseCase:          loginUseCase,
		attachProfileUseCase:  attachProfileUseCase,
		refreshTokensUseCase:  refreshTokensUseCase,
		forgotPasswordUseCase: forgotPasswordUseCase,
		resetPasswdUseCase:    resetPasswdUseCase,
		checkFieldUseCase:     checkFieldUseCase,
		authMiddleware:        authMiddleware,
	}
}

func (c *AuthController) Route(router fiber.Router) {
	authRoutes := router.Group(AUTH_CONTROLLER_NAME)

	authRoutes.Post("profile", c.authMiddleware.AccessTokenHandler, c.AttachProfile)
	authRoutes.Get("token", c.authMiddleware.AccessTokenHandler, c.Token)
	authRoutes.Get("refresh-token", c.authMiddleware.RefreshTokenHandler, c.RefreshTokens)

	authRoutes.Get("check-field", c.CheckField)

	authRoutes.Post("reset-password", c.ResetPassword)
	authRoutes.Post("forgot-password", c.ForgotPassword)

	authRoutes.Post("login", c.Login)
	authRoutes.Post("sign-up", c.SignUp)
}

// // ShowAccount godoc
//
//	@Summary	Attach a profile to the `created account`.
//	@Tags		Profiles
//	@Accept		application/json
//	@Produce	json
//
//	@Param		data			body		contracts.AttachProfileRequest	true	"AttachProfile Payload"
//	@Success	200				{object}	contracts.ProfileResponse "Attached Profile Id"
//	@security	Bearer
//
// @Failure  400       {object}  shared.ProblemDetailsExtendend   "Invalid parameters"
// @Failure  404       {object}  shared.ProblemDetails   "User not found"
// @Failure  404       {object}  shared.ProblemDetails   "Grant Type not found"
// @Failure  500       {object}  shared.ProblemDetails   "Server failed to provide an valid response"
//
//	@Router		/profile [post]
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

// ShowAccount godoc
//
//	@Summary	Validate the provided `token` and returns success or failed weather the token was signed by the `public key`.
//	@Tags		Access Credentials
//	@Accept		application/json
//	@Produce	json
//
//	@Success	200				{object}	contracts.AccountResponse "Token Payload"
//	@security	Bearer
//
// @Failure  401       {object}  shared.ProblemDetails   "Authentication Failed"
// @Failure  403       {object}  shared.ProblemDetails   "Don't have access to this resource"
// @Failure  500       {object}  shared.ProblemDetails   "Server failed to provide an valid response"
//
//	@Router		/token [get]
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

// ShowAccount godoc
//
//	@Summary	It receives the `refresh token` and validates it. Then, creates a new refresh token and access token revoking the other ones.
//	@Tags		Access Credentials
//	@Accept		application/json
//	@Produce	json
//
//	@Param		profile_id			query		string	false	"switch access to (profile_id) or if not provided default to main profile"
//	@Success	200				{object}	contracts.AuthResponse "Access Credentials"
//	@security	Bearer
//
// @Failure  401       {object}  shared.ProblemDetails   "Authentication Failed"
// @Failure  403       {object}  shared.ProblemDetails   "Don't have access to this resource"
// @Failure  500       {object}  shared.ProblemDetails   "Server failed to provide an valid response"
//
//	@Router		/refresh-token [get]
func (c *AuthController) RefreshTokens(ctx *fiber.Ctx) error {
	profileId := ctx.Query("profile_id", "")

	authId, err := middlewares.GetUserId(ctx)

	if err != nil {
		return err
	}

	if err := shared.Validate(&contracts.RefreshTokensRequest{ProfileId: profileId}); err != nil && profileId != "" {
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

// // ShowAccount godoc
//
//	@Summary	Checks if the email is not already registered on the platform.
//	@Tags		Validation
//	@Accept		application/json
//	@Produce	json
//
//	@Param		email			query		string	true	"user email"
//	@Success	200				{object}	contracts.GenericResponse "Value"
//
// @Failure  409       {object}  shared.ProblemDetails   "Invalid Email"
// @Failure  500       {object}  shared.ProblemDetails   "Server failed to provide an valid response"
//
//	@Router		/check-field [get]
func (c *AuthController) CheckField(ctx *fiber.Ctx) error {
	email := ctx.Query("email", "")

	if err := shared.Validate(&contracts.CheckEmailFieldRequest{Email: email}); err != nil {
		return fails.BAD_EMAIL
	}

	return ctx.Status(fiber.StatusOK).JSON(&contracts.GenericResponse{
		Message: strconv.FormatBool(c.checkFieldUseCase.Handle(usecases.CheckFieldInput{
			Email: email,
		})),
	})
}

// // ShowAccount godoc
//
//	@Summary	After reviving an email with the code, this endpoint will allow the user to reset his password.
//	@Tags		Authentication
//	@Accept		application/json
//	@Produce	json
//
//	@Param		data			body		contracts.ResetPasswordRequest	true	"ResetPassword Payload"
//	@Success	200				{object}	contracts.GenericResponse "Response"
//
// @Failure  401       {object}  shared.ProblemDetails   "Authentication Failed"
// @Failure  403       {object}  shared.ProblemDetails   "Don't have access to this resource"
// @Failure  500       {object}  shared.ProblemDetails   "Server failed to provide an valid response"
//
//	@Router		/reset-password [post]
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

// // ShowAccount godoc
//
//	@Summary	It will be sent an email containing an code that will be used to provide a new password.
//	@Tags		Authentication
//	@Accept		application/json
//	@Produce	json
//
//	@Param		data			body		contracts.ForgotPasswordRequest	true	"ForgotPassword Payload"
//	@Success	200				{object}	contracts.GenericResponse "Response"
//
// @Failure  401       {object}  shared.ProblemDetails   "Authentication Failed"
// @Failure  403       {object}  shared.ProblemDetails   "Don't have access to this resource"
// @Failure  500       {object}  shared.ProblemDetails   "Server failed to provide an valid response"
//
//	@Router		/forgot-password [post]
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

// // ShowAccount godoc
//
//	@Summary	Gives the `access token` and `refresh token` that is needed to interact with the platform, along with all user's profiles ids.
//	@Tags		Authentication
//	@Accept		application/json
//	@Produce	json
//
//	@Param		data			body		contracts.LoginRequest	true	"LogIn Payload"
//	@Success	200				{object}	contracts.AuthResponse "Get Access Credentials"
//
// @Failure  400       {object}  shared.ProblemDetailsExtendend   "Invalid parameters"
// @Failure  401       {object}  shared.ProblemDetails   "Authentication Failed"
// @Failure  500       {object}  shared.ProblemDetails   "Server failed to provide an valid response"
//
//	@Router		/login [post]
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var loginRequest contracts.LoginRequest

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

// // ShowAccount godoc
//
//	@Summary	Register an access key to obtain a `profileId`, which allows you to create a profile on the platform.
//	@Tags		Authentication
//	@Accept		application/json
//	@Produce	json
//
//	@Param		data			body		contracts.SignUpRequest	true	"Sign Up Payload"
//	@Success	201				{object}	contracts.AuthResponse "Account Created"
//
// @Failure  400       {object}  shared.ProblemDetailsExtendend   "Invalid parameters"
// @Failure  404       {object}  shared.ProblemDetails   "Role not found"
// @Failure  409       {object}  shared.ProblemDetails   "Invalid Password or Email already used"
// @Failure  500       {object}  shared.ProblemDetails   "Server failed to provide an valid response"
//
//	@Router		/sign-up [post]
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
