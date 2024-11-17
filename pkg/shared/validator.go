package shared

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

var (
	messages = map[string]string{
		"email":    "Email is required and must be valid.",
		"password": "Password is required and must be at least 8 characters long.",
		"role":     "Role is required and must be a positive number.",
	}
)

func ValidationFailed(errors map[string]string) *ValidationError {
	return NewBadRequestError(
		"validation-error",
		"Validation Failed",
		errors,
	)
}

func InputUnsupported(contentType string) *Error {
	return NewUnsupportedMediaError(
		"input-not-supported",
		fmt.Sprintf("Content-Type must be %s", contentType),
		"The user must send requests with a server supported format",
	)
}

func Validate(input interface{}) error {
	return validate.Struct(input)
}

func ParseBodyAndValidate(ctx *fiber.Ctx, request interface{}) error {
	errors := make(map[string]string)

	if err := ctx.BodyParser(request); err != nil {
		return InputUnsupported("application/json")
	}

	err := validate.Struct(request)

	if err == nil {
		return nil
	}

	for _, err := range err.(validator.ValidationErrors) {
		key := strings.ToLower(err.Field())

		errors[key] = messages[key]
	}

	return ValidationFailed(errors)
}
