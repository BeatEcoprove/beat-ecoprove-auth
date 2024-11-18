package middlewares

import (
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
		authRepository repositories.IAuthRepository
		tokenService   services.ITokenService
	}
)

func NewAuthorizationMiddleware(
	authRepository repositories.IAuthRepository,
	tokenService services.ITokenService,
) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		authRepository: authRepository,
		tokenService:   tokenService,
	}
}

func (am *AuthorizationMiddleware) tokenHandler(ctx *fiber.Ctx, validateToken func(claims *services.AuthClaims, token *jwt.Token) error) error {
	return jwtware.New(jwtware.Config{
		Claims: &services.AuthClaims{},
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwt.SigningMethodRS256.Name,
			Key:    services.PubKey,
		},
		SuccessHandler: func(ctx *fiber.Ctx) error {
			token, claims, err := GetClaims(ctx)

			if err != nil {
				return shared.WriteProblemDetails(ctx, *fails.INVALID_ACCESS_TOKEN)
			}

			if err := validateToken(claims, token); err != nil {
				return err
			}

			if ok := am.authRepository.ExistsUserWithId(claims.Subject); !ok {
				return shared.WriteProblemDetails(ctx, *fails.DONT_HAVE_ACCESS_TO_RESOURCE)
			}

			return ctx.Next()
		},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return shared.WriteProblemDetails(ctx, *fails.DONT_HAVE_ACCESS_TO_RESOURCE)
		},
	})(ctx)
}

func (am *AuthorizationMiddleware) AccessTokenHandler(ctx *fiber.Ctx) error {
	return am.tokenHandler(ctx, func(claims *services.AuthClaims, token *jwt.Token) error {
		err := am.tokenService.ValidateToken(claims.Subject, token.Raw, services.AccessTokenKey)

		if err != nil {
			return fails.INVALID_ACCESS_TOKEN
		}

		return nil
	})
}

func (am *AuthorizationMiddleware) RefreshTokenHandler(ctx *fiber.Ctx) error {
	return am.tokenHandler(ctx, func(claims *services.AuthClaims, token *jwt.Token) error {
		err := am.tokenService.ValidateToken(claims.Subject, token.Raw, services.RefreshTokenKey)

		if err != nil {
			return fails.INVALID_REFRESH_TOKEN
		}

		return nil
	})
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
