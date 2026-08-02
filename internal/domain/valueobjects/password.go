package valueobjects

import (
	"errors"
	"unicode"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	hash string
}

func NewPassword(plain string) (Password, error) {
	if err := validateStrength(plain); err != nil {
		return Password{}, err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, err
	}
	return Password{hash: string(hashed)}, nil
}

func NewPasswordFromHash(hash string) Password {
	return Password{hash: hash}
}

func (p Password) Matches(plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plain)) == nil
}

func (p Password) Hash() string {
	return p.hash
}

func validateStrength(plain string) error {
	var errs []error

	if len(plain) < 8 {
		errs = append(errs, domainerrors.ErrWeakPassword)
	}
	var hasUpper, hasLower, hasDigit, hasSymbol bool
	for _, c := range plain {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsSymbol(c), unicode.IsPunct(c):
			hasSymbol = true
		}
	}
	if !hasUpper {
		errs = append(errs, domainerrors.ErrPasswordNoUpper)
	}
	if !hasLower {
		errs = append(errs, domainerrors.ErrPasswordNoLower)
	}
	if !hasDigit {
		errs = append(errs, domainerrors.ErrPasswordNoNumber)
	}
	if !hasSymbol {
		errs = append(errs, domainerrors.ErrPasswordNoSpecial)
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.Join(errs...)
}
