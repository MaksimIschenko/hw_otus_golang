package hw02unpackstring

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

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
		{input: "a!b,c$", expected: "a!b,c$"},
		{input: `a!b,5c$`, expected: "a!b,,,,,c$"},
		{input: `a!b,5c$\4`, expected: "a!b,,,,,c$4"},
		{input: `\5\5\5`, expected: `555`},
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
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestDeleteLastRune(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "abc", expected: "ab"},
		{input: "a", expected: ""},
		{input: "", expected: ""},
		{input: "a\n", expected: "a"},
		{input: "a\n\n", expected: "a\n"},
		{input: "a5", expected: "a"},
		{input: "555", expected: "55"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			var sb strings.Builder
			sb.WriteString(tc.input)
			DeleteLastRune(&sb)
			require.Equal(t, tc.expected, sb.String())
		})
	}
}
