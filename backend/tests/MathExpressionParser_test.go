package main

import (
	"216/internal/orchestrator/Services"
	"testing"
)

func TestIsValidExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"( 2 + 3 ) * 5", true},
		{"(2+3)*5)", false},
		{"2 + 3", true},
		{"(2 + 3", false},
		{"(2 + 3) * (5 + 7)", true},
		{"(2 + 3)) * 5", false},
		{"* (2 + 3)*5", false},
		{"(2 + 3)*", false},
		{"+ (2 + 3)*5", false},
	}

	for _, test := range tests {
		result := Services.IsValidExpression(test.input)
		if result != test.expected {
			t.Errorf("For expression '%s', got %t, want %t", test.input, result, test.expected)
		}
	}
}

func TestToPostfix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"3 + 4 * 2 / ( 1 - 5 )", "3 4 2 * 1 5 - / +"},
		{"-4 + 5 / -2", "-4 5 -2 / +"},
		{"( 6 / 2 ) + 7 * 8", "6 2 / 7 8 * +"},
	}

	for _, test := range tests {
		result := Services.ToPostfix(test.input)
		if result != test.expected {
			t.Errorf("For input '%s', got '%s', want '%s'", test.input, result, test.expected)
		}
	}
}
