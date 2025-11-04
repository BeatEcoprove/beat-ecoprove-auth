package internal

import (
	"strconv"

	"github.com/BeatEcoprove/identityService/internal/middlewares"
	"github.com/BeatEcoprove/identityService/internal/usecases"
	"github.com/BeatEcoprove/identityService/pkg/contracts"
	fails "github.com/BeatEcoprove/identityService/pkg/errors"
	"github.com/BeatEcoprove/identityService/pkg/services"
	"github.com/BeatEcoprove/identityService/pkg/shared"
	"github.com/gofiber/fiber/v2"
)

const (
	AuthRoutes         = "auth"
	ProfileRoutes      = "profiles"
	AvailabilityRoutes = "availability"
	GroupRoutes        = "groups"

	GrantTypePassword      = "password"
	GrantTypeRefreshTokens = "refresh_token"
)

type AuthController struct {
	signUpUseCase         *usecases.SignUpUseCase
	loginUseCase          *usecases.LoginUseCase
	attachProfileUseCase  *usecases.AttachProfileUseCase
	refreshTokensUseCase  *usecases.RefreshTokensUseCase
	forgotPasswordUseCase *usecases.ForgotPasswordUseCase
	resetPasswdUseCase    *usecases.ResetPasswdUseCase
	checkFieldUseCase     *usecases.CheckFieldUseCase
	fechPermissions       *usecases.FetchGroupUserPermissionsUseCase

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
	fechPermissions *usecases.FetchGroupUserPermissionsUseCase,
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
		fechPermissions:       fechPermissions,
	}
}

func (c *AuthController) Route(router fiber.Router) {
	// heath check - router

	authRoutes := router.Group(AuthRoutes)
	authRoutes.Post("reset-password", c.ResetPassword)
	authRoutes.Post("forgot-password", c.ForgotPassword)
	authRoutes.Post("token", c.Token)
	authRoutes.Post("sign-up", c.SignUp)

	profileRoutes := authRoutes.Group(ProfileRoutes)
	profileRoutes.Post("reserve", c.authMiddleware.AccessTokenHandler, c.AttachProfile)
	profileRoutes.Get("me", c.authMiddleware.AccessTokenHandler, c.Me)

	availabilityRoutes := authRoutes.Group(AvailabilityRoutes)
	availabilityRoutes.Get("check-field", c.CheckField)

	groupRoutes := authRoutes.Group(GroupRoutes)
	groupRoutes.Get("permissions", c.FetchGroupPermissions)
}

// // ShowAccount godoc
//
//	@Summary	OAuth2 endpoints for authorization
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
//	@Router		/token [post]
func (c *AuthController) Token(ctx *fiber.Ctx) error {
	var request contracts.TokenRequest

	if err := shared.ParseBodyAndValidate(ctx, &request); err != nil {
		return err
	}

	switch request.GrantType {
	case GrantTypePassword:
		return c.handleLogin(ctx)
	case GrantTypeRefreshTokens:
		return c.handleRefreshTokens(ctx)
	default:
		return fails.DONT_HAVE_ACCESS_TO_RESOURCE
	}
}

func (c *AuthController) handleLogin(ctx *fiber.Ctx) error {
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

func (c *AuthController) handleRefreshTokens(ctx *fiber.Ctx) error {
	var refreshTokenRequest contracts.RefreshTokenRequest

	if err := shared.ParseBodyAndValidate(ctx, &refreshTokenRequest); err != nil {
		return err
	}

	var claims services.AuthClaims
	if err := services.GetClaims(refreshTokenRequest.Token, &claims, services.Refresh); err != nil {
		return err
	}

	response, err := c.refreshTokensUseCase.Handle(usecases.RefreshTokensInput{
		AuthId:    claims.Subject,
		ProfileId: refreshTokenRequest.ProfileID,
	})

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// // ShowAccount godoc
//
//	@Summary	Attach a profile to the `created account`.
//	@Tags		Account
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
//	@Router		/profiles [post]
func (c *AuthController) AttachProfile(ctx *fiber.Ctx) error {
	authID, err := middlewares.GetUserID(ctx)

	if err != nil {
		return err
	}

	response, err := c.attachProfileUseCase.Handle(usecases.AttachProfileInput{
		AuthId:           authID,
		ProfileGrantType: 1,
	})

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}

// ShowAccount godoc
//
//	@Summary	Validate the provided `token` and returns success or failed weather the token was signed by the `public key`.
//	@Tags		Account
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
//	@Router		/profiles/me [get]
func (c *AuthController) Me(ctx *fiber.Ctx) error {
	_, claims, err := middlewares.GetClaims(ctx)

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(&contracts.AccountResponse{
		UserID:     claims.Subject,
		Email:      claims.Email,
		ProfileID:  claims.ProfileID,
		ProfileIds: claims.ProfileIds,
		Role:       claims.Role,
	})
}

// // ShowAccount godoc
//
//	@Summary	Checks if the email is not already registered on the platform.
//	@Tags		Availability
//	@Accept		application/json
//	@Produce	json
//
//	@Param		email			query		string	true	"user email"
//	@Success	200				{object}	contracts.GenericResponse "Value"
//
// @Failure  409       {object}  shared.ProblemDetails   "Invalid Email"
// @Failure  500       {object}  shared.ProblemDetails   "Server failed to provide an valid response"
//
//	@Router		/availability/check-field [get]
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

// // ShowAccount godoc
//
//	@Summary	Register an access key to obtain a `profileId`, which allows you to create a profile on the platform.
//	@Tags		Internal
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
//	@Router		/groups/permissions [post]
func (c *AuthController) FetchGroupPermissions(ctx *fiber.Ctx) error {
	var fetchPermissionsRequest contracts.GroupPermissionsRequest

	if err := shared.ParseBodyAndValidate(ctx, &fetchPermissionsRequest); err != nil {
		return err
	}

	response, err := c.fechPermissions.Handle(usecases.FetchGroupUserPermissionsInput{
		GroupID:  fetchPermissionsRequest.GroupID,
		MmeberID: fetchPermissionsRequest.MemberID,
	})

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
