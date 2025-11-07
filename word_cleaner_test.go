package main

import "testing"

func TestWordCleaner(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"no replacements", "string with no changes", "string with no changes"},
		{"replace kerfuffle", "string with kerfuffle replaced", "string with **** replaced"},
		{"replace sharbert", "string with sharbert replaced", "string with **** replaced"},
		{"replace fornax", "string with fornax replaced", "string with **** replaced"},

		{"test uppercase KERFUFFLE", "string with KERFUFFLE replaced", "string with **** replaced"},
		{"test not replaced due to punctuation", "string with kerfuffle! not replaced", "string with kerfuffle! not replaced"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wordCleaner(tt.input)
			if got != tt.expected {
				t.Errorf("wordCleaner(%q), output: %q, expected: %q", tt.input, got, tt.expected)
			}
		})
	}
}
