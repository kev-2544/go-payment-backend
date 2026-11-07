// Package money converts between decimal amount strings (e.g. "10.50")
// and int64 amounts in a currency's smallest unit (e.g. 1050 cents).
// It never uses floats, so there are no rounding errors.
package money

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrUnsupportedCurrency = errors.New("unsupported currency")
	ErrInvalidAmount       = errors.New("invalid amount")
)

// decimals is the number of minor-unit digits for each supported currency.
var decimals = map[string]int{
	"KES": 2,
	"USD": 2,
	"EUR": 2,
	"GBP": 2,
	"UGX": 0,
}

// NormalizeCurrency upper-cases and validates a currency code.
func NormalizeCurrency(currency string) (string, error) {
	code := strings.ToUpper(strings.TrimSpace(currency))
	if _, ok := decimals[code]; !ok {
		return "", fmt.Errorf("%w: %q", ErrUnsupportedCurrency, currency)
	}
	return code, nil
}

// Parse converts a decimal string like "10.50" into minor units (1050).
// It rejects negative numbers, signs, exponents and more decimal places
// than the currency allows.
func Parse(amount, currency string) (int64, error) {
	code, err := NormalizeCurrency(currency)
	if err != nil {
		return 0, err
	}
	places := decimals[code]

	s := strings.TrimSpace(amount)
	whole, frac, hasDot := strings.Cut(s, ".")
	if whole == "" || !isDigits(whole) || (hasDot && (frac == "" || !isDigits(frac))) {
		return 0, fmt.Errorf("%w: %q", ErrInvalidAmount, amount)
	}
	if len(frac) > places {
		return 0, fmt.Errorf("%w: %q has more than %d decimal places for %s", ErrInvalidAmount, amount, places, code)
	}

	// Pad the fraction to the currency's precision: "10.5" -> "10" + "50".
	digits := whole + frac + strings.Repeat("0", places-len(frac))
	minor, err := strconv.ParseInt(digits, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %q is too large", ErrInvalidAmount, amount)
	}
	return minor, nil
}

// Format converts minor units (1050) back to a decimal string ("10.50").
func Format(minor int64, currency string) (string, error) {
	code, err := NormalizeCurrency(currency)
	if err != nil {
		return "", err
	}
	places := decimals[code]

	sign := ""
	if minor < 0 {
		sign = "-"
		minor = -minor
	}
	s := strconv.FormatInt(minor, 10)
	if places == 0 {
		return sign + s, nil
	}
	if len(s) <= places {
		s = strings.Repeat("0", places-len(s)+1) + s
	}
	return sign + s[:len(s)-places] + "." + s[len(s)-places:], nil
}

func isDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
