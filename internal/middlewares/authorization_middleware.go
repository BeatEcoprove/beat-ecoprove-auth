package middlewares

import (
	"github.com/BeatEcoprove/identityService/internal/adapters"
	"github.com/BeatEcoprove/identityService/internal/repositories"
	fails "github.com/BeatEcoprove/identityService/pkg/errors"
	"github.com/BeatEcoprove/identityService/pkg/services"
	"github.com/BeatEcoprove/identityService/pkg/shared"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type (
	AuthorizationMiddleware struct {
		tokenService services.ITokenService
	}
)

func NewAuthorizationMiddleware(
	tokenService services.ITokenService,
) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		tokenService: tokenService,
	}
}

func (am *AuthorizationMiddleware) Handle(ctx *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		Claims: &services.AuthClaims{},
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwt.SigningMethodRS256.Name,
			Key:    services.PubKey,
		},
		SuccessHandler: func(ctx *fiber.Ctx) error {
			authRepository := repositories.NewAuthRepository(adapters.GetDatabase())

			token, claims, err := GetClaims(ctx)

			if err != nil {
				return shared.WriteProblemDetails(ctx, *fails.INVALID_ACCESS_TOKEN)
			}

			if err := am.tokenService.ValidateToken(claims.Subject, token.Raw); err != nil {
				return shared.WriteProblemDetails(ctx, *fails.INVALID_ACCESS_TOKEN)
			}

			if ok := authRepository.ExistsUserWithId(claims.Subject); !ok {
				return shared.WriteProblemDetails(ctx, *fails.DONT_HAVE_ACCESS_TO_RESOURCE)
			}

			return ctx.Next()
		},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return shared.WriteProblemDetails(ctx, *fails.DONT_HAVE_ACCESS_TO_RESOURCE)
		},
	})(ctx)
}

func GetClaims(ctx *fiber.Ctx) (*jwt.Token, *services.AuthClaims, error) {
	token, ok := ctx.Locals("user").(*jwt.Token)

	if !ok {
		return nil, nil, fails.INVALID_ACCESS_TOKEN
	}

	claims, ok := token.Claims.(*services.AuthClaims)

	if !ok {
		return nil, nil, fails.INVALID_ACCESS_TOKEN
	}

	return token, claims, nil
}

func GetUserId(ctx *fiber.Ctx) (string, error) {
	_, claims, err := GetClaims(ctx)

	if err != nil {
		return "", err
	}

	return claims.Subject, nil
}
