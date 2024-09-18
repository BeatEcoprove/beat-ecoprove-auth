package shared

import "github.com/gofiber/fiber/v2"

type Controller interface {
	Route(router fiber.Router)
}

//
// func GetUserId(ctx *fiber.Ctx) string {
// 	user := ctx.Locals("user").(*jwt.Token)
// 	claims, _ := user.Claims.(*jwtService.AuthClaims)
//
// 	return claims.Subject
// }
//
// func GetRole(ctx *fiber.Ctx) string {
// 	user := ctx.Locals("user").(*jwt.Token)
// 	claims, _ := user.Claims.(*jwtService.AuthClaims)
//
// 	return claims.Role
// }
