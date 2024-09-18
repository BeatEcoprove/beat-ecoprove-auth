package shared

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	MIMEApplicationProblemJSON string = "application/problem+json"
	RFC_URL                    string = "https://datatracker.ietf.org/doc/html/rfc2616#section-10"
)

type (
	ProblemDetails struct {
		Type       string                 `json:"type"`
		Title      string                 `json:"title"`
		Status     int                    `json:"status"`
		Detail     string                 `json:"detail,omitempty"`
		Instance   string                 `json:"instance,omitempty"`
		Extensions map[string]interface{} `json:"extensions,omitempty"`
	}
)

func parseStatusIndex(status int) (string, string) {
	statusCode := strconv.Itoa(status)

	firstDigit := string(statusCode[0])
	lastDigit := string(statusCode[len(statusCode)-1])

	lastDigitInt, _ := strconv.Atoi(lastDigit)

	nextDigit := (lastDigitInt + 1) % 10

	return firstDigit, strconv.Itoa(nextDigit)
}

func WriteProblemDetails(ctx *fiber.Ctx, err Error, extensions ...map[string]interface{}) error {
	section, index := parseStatusIndex(err.Status)

	problem := ProblemDetails{
		Type: fmt.Sprintf("%s.%s.%s", RFC_URL, section, index),
		// Title:    locales.LocalizeHeader(ctx, err.Id, err.Title),
		Title:  err.Title,
		Status: err.Status,
		// Detail:   locales.Localize(ctx, err.Id, err.Detail),
		Detail:   err.Detail,
		Instance: ctx.Path(),
	}

	if len(extensions) > 0 {
		problem.AddExtensions(extensions[0])
	}

	ctx.Status(err.Status)
	return ctx.JSON(problem, MIMEApplicationProblemJSON)
}

func (pd *ProblemDetails) AddExtensions(extensions map[string]interface{}) {
	pd.Extensions = extensions
}
