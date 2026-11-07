package money

import (
	"errors"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		amount, currency string
		want             int64
		wantErr          error
	}{
		{"10.50", "KES", 1050, nil},
		{"10.5", "kes", 1050, nil},
		{"10", "USD", 1000, nil},
		{"0.01", "USD", 1, nil},
		{" 7.25 ", "EUR", 725, nil},
		{"5000", "UGX", 5000, nil},
		{"10.505", "KES", 0, ErrInvalidAmount},
		{"10.5", "UGX", 0, ErrInvalidAmount},
		{"-10", "KES", 0, ErrInvalidAmount},
		{"1e3", "KES", 0, ErrInvalidAmount},
		{"abc", "KES", 0, ErrInvalidAmount},
		{"", "KES", 0, ErrInvalidAmount},
		{".5", "KES", 0, ErrInvalidAmount},
		{"10.", "KES", 0, ErrInvalidAmount},
		{"99999999999999999999", "KES", 0, ErrInvalidAmount},
		{"10", "XYZ", 0, ErrUnsupportedCurrency},
	}
	for _, tt := range tests {
		got, err := Parse(tt.amount, tt.currency)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("Parse(%q, %q) error = %v, want %v", tt.amount, tt.currency, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("Parse(%q, %q) = %d, want %d", tt.amount, tt.currency, got, tt.want)
		}
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		minor    int64
		currency string
		want     string
	}{
		{1050, "KES", "10.50"},
		{1, "USD", "0.01"},
		{0, "USD", "0.00"},
		{100, "EUR", "1.00"},
		{5000, "UGX", "5000"},
	}
	for _, tt := range tests {
		got, err := Format(tt.minor, tt.currency)
		if err != nil || got != tt.want {
			t.Errorf("Format(%d, %q) = %q, %v; want %q", tt.minor, tt.currency, got, err, tt.want)
		}
	}
}
