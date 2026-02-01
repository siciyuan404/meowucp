package service

import "testing"

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"test@example.com", "t*t@example.com"},
		{"ab@example.com", "ab@example.com"},
		{"a@example.com", "a@example.com"},
		{"", ""},
	}

	for _, tt := range tests {
		result := MaskEmail(tt.input)
		if result != tt.expected {
			t.Errorf("MaskEmail(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestMaskPhone(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"13812345678", "13****5678"},
		{"123456", "12****56"},
		{"123", "***"},
		{"", ""},
	}

	for _, tt := range tests {
		result := MaskPhone(tt.input)
		if result != tt.expected {
			t.Errorf("MaskPhone(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestMaskCreditCard(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"4111111111111111", "************1111"},
		{"4111", "4111"},
		{"4111111111111", "*********1111"},
		{"", ""},
	}

	for _, tt := range tests {
		result := MaskCreditCard(tt.input)
		if result != tt.expected {
			t.Errorf("MaskCreditCard(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}
