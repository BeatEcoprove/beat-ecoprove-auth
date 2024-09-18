package shared

type (
	Error struct {
		Id     string
		Status int
		Title  string
		Detail string
	}
)

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
