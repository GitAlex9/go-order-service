package validation

import "strings"

type FieldError struct {
	Field   string
	Message string
}

type ValidationErrors struct {
	Errors []FieldError
}

func (v *ValidationErrors) Add(field, message string) {
	v.Errors = append(v.Errors, FieldError{Field: field, Message: message})
}

func (v *ValidationErrors) HasErrors() bool {
	return len(v.Errors) > 0
}

func (v *ValidationErrors) Error() string {
	msgs := make([]string, len(v.Errors))
	for i, e := range v.Errors {
		msgs[i] = e.Field + ": " + e.Message
	}
	return strings.Join(msgs, "; ")
}
