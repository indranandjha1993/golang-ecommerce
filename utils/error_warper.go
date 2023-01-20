package utils

type ErrorWarper struct {
	Errors map[string]string `json:"errors"`
}

func (e *ErrorWarper) Add(field, message string) {
	if e.Errors == nil {
		e.Errors = make(map[string]string)
	}
	e.Errors[field] = message
}
