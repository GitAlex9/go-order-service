package valueobjects

import (
	"regexp"
	"strings"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type Email struct {
	value string
}

func NewEmail(raw string) (Email, error) {
	value := strings.TrimSpace(strings.ToLower(raw))
	if value == "" {
		return Email{}, domainerrors.ErrInvalidEmail //Todo Criar error 'email can't be empty'
	}
	if !emailRegex.MatchString(value) {
		return Email{}, domainerrors.ErrInvalidEmail
	}
	return Email{value: value}, nil
}

func (e Email) String() string {
	return e.value
}

func (e Email) Equals(other Email) bool {
	return e.value == other.value
}
