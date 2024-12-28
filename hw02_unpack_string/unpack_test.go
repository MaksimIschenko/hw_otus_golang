package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsLetterOrNewString(t *testing.T) {
	tests := []struct {
		input    rune
		expected bool
	}{
		{input: rune('d'), expected: true},
		{input: rune('\n'), expected: true},
		{input: rune('5'), expected: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(string(tc.input), func(t *testing.T) {
			result := IsLetterOrNewString(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestIsContainAllowedSymbols(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{input: "abcd", expected: true},
		{input: "abcd12345", expected: true},
		{input: "12345", expected: true},
		{input: "a1b2c3d4e5", expected: true},
		{input: "", expected: true},
		{input: "abcd\n", expected: true},
		{input: "\n\n\n\n\n", expected: true},
		{input: "abcd!", expected: false},
		{input: "!ab\ncd!", expected: false},
		{input: "abcd.", expected: false},
		{input: "!", expected: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result := IsContainAllowedSymbols(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		{input: "d\n5abc\n", expected: "d\n\n\n\n\nabc\n"},
		{input: "d\n5abc\n3", expected: "d\n\n\n\n\nabc\n\n\n"},
		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", "aa!d5"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
