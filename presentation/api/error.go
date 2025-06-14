package api

type Error struct {
	Errors []string `json:"errors"`
}

func NewError(err error) *Error {
	return &Error{
		Errors: []string{err.Error()},
	}
}
