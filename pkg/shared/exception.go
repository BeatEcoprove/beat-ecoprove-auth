package shared

import "github.com/gofiber/fiber/v2"

type (
	Error struct {
		Id     string
		Status int
		Title  string
		Detail string
	}

	ValidationError struct {
		Id      string
		Status  int
		Title   string
		Details map[string]string
	}
)

func NewBadRequestError(id, title string, details map[string]string) *ValidationError {
	return &ValidationError{
		Id:      id,
		Status:  fiber.StatusBadRequest,
		Title:   title,
		Details: details,
	}
}

func NewUnauthorizedError(id, title, detail string) *Error {
	return &Error{
		Id:     id,
		Status: fiber.StatusUnauthorized,
		Title:  title,
		Detail: detail,
	}
}

func NewForbiddenError(id, title, detail string) *Error {
	return &Error{
		Id:     id,
		Status: fiber.StatusForbidden,
		Title:  title,
		Detail: detail,
	}
}

func NewNotFoundError(id, title, detail string) *Error {
	return &Error{
		Id:     id,
		Status: fiber.StatusNotFound,
		Title:  title,
		Detail: detail,
	}
}

func NewConflitError(id, title, detail string) *Error {
	return &Error{
		Id:     id,
		Status: fiber.StatusConflict,
		Title:  title,
		Detail: detail,
	}
}

func NewUnsupportedMediaError(id, title, detail string) *Error {
	return &Error{
		Id:     id,
		Status: fiber.StatusUnsupportedMediaType,
		Title:  title,
		Detail: detail,
	}
}

func NewInternalError(id, title, detail string) *Error {
	return &Error{
		Id:     id,
		Status: fiber.StatusInternalServerError,
		Title:  title,
		Detail: detail,
	}
}

func NewError(id string, status int, title string, detail string) *Error {
	return &Error{
		Id:     id,
		Status: status,
		Title:  title,
		Detail: detail,
	}
}

func (e *Error) Error() string {
	return e.Detail
}

func (e *ValidationError) Error() string {
	return e.Title
}
