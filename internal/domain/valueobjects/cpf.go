package valueobjects

import (
	"regexp"
	"strconv"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

type CPF struct {
	value string // validar com 11 dígitos, sem máscara
}

var nonDigits = regexp.MustCompile(`\D`)

func NewCPF(raw string) (CPF, error) {
	digits := nonDigits.ReplaceAllString(raw, "")
	if len(digits) != 11 {
		return CPF{}, domainerrors.ErrInsufficientCPFLength
	}
	if allEqual(digits) || !hasValidCheckDigits(digits) {
		return CPF{}, domainerrors.ErrInvalidCPF
	}
	return CPF{value: digits}, nil
}

func (c CPF) String() string { return c.value }

func (c CPF) Formatted() string {
	return c.value[0:3] + "." + c.value[3:6] + "." + c.value[6:9] + "-" + c.value[9:11]
}

func allEqual(digits string) bool {
	for i := 1; i < len(digits); i++ {
		if digits[i] != digits[0] {
			return false
		}
	}
	return true
}

func hasValidCheckDigits(digits string) bool {
	calc := func(length int) int {
		sum, weight := 0, length+1
		for i := 0; i < length; i++ {
			n, _ := strconv.Atoi(string(digits[i]))
			sum += n * weight
			weight--
		}
		if rest := sum % 11; rest < 2 {
			return 0
		} else {
			return 11 - rest
		}
	}
	d1 := calc(9)
	base := digits[0:9] + strconv.Itoa(d1)
	sum, weight := 0, 11
	for i := 0; i < 10; i++ {
		n, _ := strconv.Atoi(string(base[i]))
		sum += n * weight
		weight--
	}
	d2 := 0
	if rest := sum % 11; rest >= 2 {
		d2 = 11 - rest
	}
	return digits[9:11] == strconv.Itoa(d1)+strconv.Itoa(d2)
}
