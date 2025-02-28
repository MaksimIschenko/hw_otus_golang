package main

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		name        string
		fromPath    string
		toPath      string
		correctPath string
		offset      int64
		limit       int64
		expectedErr error
	}{
		{
			name:        "successful copy",
			fromPath:    "testdata/input.txt",
			toPath:      "testdata/tmp/out_offset0_limit0.txt",
			correctPath: "testdata/out_offset0_limit0.txt",
			offset:      0,
			limit:       0,
			expectedErr: nil,
		},
		{
			name:        "successful copy with limit 10",
			fromPath:    "testdata/input.txt",
			toPath:      "testdata/tmp/out_offset0_limit10.txt",
			correctPath: "testdata/out_offset0_limit10.txt",
			offset:      0,
			limit:       10,
			expectedErr: nil,
		},
		{
			name:        "successful copy with limit 1000",
			fromPath:    "testdata/input.txt",
			toPath:      "testdata/tmp/out_offset0_limit1000.txt",
			correctPath: "testdata/out_offset0_limit1000.txt",
			offset:      0,
			limit:       1000,
			expectedErr: nil,
		},
		{
			name:        "successful copy with offset 100 and limit 1000",
			fromPath:    "testdata/input.txt",
			toPath:      "testdata/tmp/out_offset100_limit1000.txt",
			correctPath: "testdata/out_offset100_limit1000.txt",
			offset:      100,
			limit:       1000,
			expectedErr: nil,
		},
	}

	err := os.Mkdir("testdata/tmp", 0o755)
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll("testdata/tmp")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Copy(tt.fromPath, tt.toPath, tt.offset, tt.limit)
			require.Equal(t, tt.expectedErr, err)

			if err == nil {
				data, err := os.ReadFile(tt.toPath)
				require.NoError(t, err)

				correctData, err := os.ReadFile(tt.correctPath)
				require.NoError(t, err)

				require.True(t, bytes.Equal(data, correctData), "File contents are not equal:\n%s\n%s", tt.toPath, tt.correctPath)
			}
		})
	}
}

func TestErrArgs(t *testing.T) {
	tests := []struct {
		name        string
		fromPath    string
		toPath      string
		correctPath string
		offset      int64
		limit       int64
		expectedErr error
	}{
		{
			name:        "unsupported file",
			fromPath:    "testdata",
			toPath:      "testdata/tmp/out_offset0_limit0.txt",
			correctPath: "testdata/out_offset0_limit0.txt",
			offset:      0,
			limit:       0,
			expectedErr: ErrUnsupportedFile,
		},
		{
			name:        "offset exceeds file size",
			fromPath:    "testdata/input.txt",
			toPath:      "testdata/tmp/out_offset10000_limit0.txt",
			correctPath: "testdata/out_offset10000_limit0.txt",
			offset:      10000,
			limit:       0,
			expectedErr: ErrOffsetExceedsFileSize,
		},
		{
			name:        "file not found",
			fromPath:    "testdata/notfound.txt",
			toPath:      "testdata/tmp/out_offset0_limit0.txt",
			correctPath: "testdata/out_offset0_limit0.txt",
			offset:      0,
			limit:       0,
			expectedErr: os.ErrNotExist,
		},
		{
			name:        "unknown file size",
			fromPath:    "/dev/urandom",
			toPath:      "testdata/tmp/out_offset0_limit0.txt",
			correctPath: "testdata/out_offset0_limit0.txt",
			offset:      0,
			limit:       0,
			expectedErr: ErrUnsupportedFile,
		},
		// {
		// 	name: "error reading file",
		// 	fromPath: "testdata/input.txt",
		// 	toPath: "testdata/tmp/out_offset0_limit0.txt",
		// 	correctPath: "testdata/out_offset0_limit0.txt",
		// 	offset: 0,
		// 	limit: 0,
		// 	expectedErr: ErrReadFile,
		// },
		// {
		// 	name: "error writing file",
		// 	fromPath: "testdata/input.txt",
		// 	toPath: "testdata/tmp",
		// 	correctPath: "testdata/out_offset0_limit0.txt",
		// 	offset: 0,
		// 	limit: 0,
		// 	expectedErr: ErrWriteFile,
		// },
		// {
		// 	name: "error seeking file",
		// 	fromPath: "testdata/input.txt",
		// 	toPath: "testdata/tmp/out_offset0_limit0.txt",
		// 	correctPath: "testdata/out_offset0_limit0.txt",
		// 	offset: 0,
		// 	limit: 0,
		// 	expectedErr: ErrSeekFile,
		// },
		// {
		// 	name: "error copying file",
		// 	fromPath: "testdata/input.txt",
		// 	toPath: "testdata/tmp/out_offset0_limit0.txt",
		// 	correctPath: "testdata/out_offset0_limit0.txt",
		// 	offset: 0,
		// 	limit: 0,
		// 	expectedErr: ErrCopyFile,
		// },
	}

	err := os.Mkdir("testdata/tmp", 0o755)
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll("testdata/tmp")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Copy(tt.fromPath, tt.toPath, tt.offset, tt.limit)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestCheckArgs(t *testing.T) {
	tests := []struct {
		name        string
		from        string
		to          string
		offset      int64
		limit       int64
		expectedErr error
	}{
		{
			name:        "successful check",
			from:        "testdata/input.txt",
			to:          "testdata/tmp/out_offset0_limit0.txt",
			offset:      0,
			limit:       0,
			expectedErr: nil,
		},
		{
			name:        "empty source file path",
			from:        "",
			to:          "testdata/tmp/out_offset0_limit0.txt",
			offset:      0,
			limit:       0,
			expectedErr: errors.New("source file path is empty"),
		},
		{
			name:        "empty destination file path",
			from:        "testdata/input.txt",
			to:          "",
			offset:      0,
			limit:       0,
			expectedErr: errors.New("destination file path is empty"),
		},
		{
			name:        "negative offset",
			from:        "testdata/input.txt",
			to:          "testdata/tmp/out_offset0_limit0.txt",
			offset:      -1,
			limit:       0,
			expectedErr: errors.New("offset cannot be negative"),
		},
		{
			name:        "negative limit",
			from:        "testdata/input.txt",
			to:          "testdata/tmp/out_offset0_limit0.txt",
			offset:      0,
			limit:       -1,
			expectedErr: errors.New("limit cannot be negative"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckArgs(tt.from, tt.to, tt.offset, tt.limit)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestGetFileSize(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expected    int64
		expectedErr error
	}{
		{
			name:        "successful get file size",
			path:        "testdata/input.txt",
			expected:    6617,
			expectedErr: nil,
		},
		{
			name:        "file not found",
			path:        "testdata/notfound.txt",
			expected:    0,
			expectedErr: os.ErrNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size, err := GetFileSize(tt.path)
			require.Equal(t, tt.expected, size)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
