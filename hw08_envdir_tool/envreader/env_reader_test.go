package envreader

import (
	"errors"
	"os"
	"testing"
)

func TestReadDir(t *testing.T) {
	tests := []struct {
		name        string
		dir         string
		expected    Environment
		expectedErr error
	}{
		{
			name: "Correct envdir",
			dir:  "../testdata/env",
			expected: Environment{
				"BAR":   EnvValue{"bar", false},
				"EMPTY": {"", true},
				"FOO":   {"   foo\nwith new line", false},
				"HELLO": {"hello", false},
				"UNSET": {"", true},
			},
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ReadDir(tt.dir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Fatalf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func TestEmptyDir(t *testing.T) {
	tests := []struct {
		name        string
		dir         string
		expected    Environment
		expectedErr error
	}{
		{
			name:        "Empty envdir",
			dir:         "../testdata/tmp",
			expected:    Environment{},
			expectedErr: nil,
		},
	}

	os.Mkdir("../testdata/tmp", 0o755)
	defer os.RemoveAll("../testdata/tmp")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ReadDir(tt.dir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Fatalf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func TestReadDirError(t *testing.T) {
	tests := []struct {
		name        string
		dir         string
		expected    Environment
		expectedErr error
	}{
		{
			name:        "Error reading envdir",
			dir:         "../testdata/err",
			expected:    Environment{},
			expectedErr: errors.New("open ../testdata/err: no such file or directory"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ReadDir(tt.dir)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
	}
}

func TestProcessBytes(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected []byte
	}{
		{
			name:     "Empty data",
			data:     []byte{},
			expected: []byte{},
		},
		{
			name:     "Data with spaces",
			data:     []byte("   "),
			expected: []byte{},
		},
		{
			name:     "Data with spaces and new line",
			data:     []byte("   \n"),
			expected: []byte("   \n"),
		},
		{
			name:     "Data with spaces, null and new line",
			data:     []byte{0x20, 0x20, 0x20, 0x66, 0x6f, 0x6f, 0x0, 0x77, 0x69, 0x74, 0x68},
			expected: []byte("   foo\nwith"),
		},
		{
			name:     "Data without spaces and without new line",
			data:     []byte("foo"),
			expected: []byte("foo"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessBytes(tt.data)
			if string(result) != string(tt.expected) {
				t.Fatalf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}
