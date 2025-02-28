package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	pb "github.com/cheggaaa/pb/v3"
)

const (
	BarWidth = 100
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

// ComparePaths compares two paths.
func ComparePaths(path1, path2 string) (bool, error) {
	absPath1, err := filepath.Abs(path1)
	if err != nil {
		return false, err
	}
	absPath2, err := filepath.Abs(path2)
	if err != nil {
		return false, err
	}
	return absPath1 == absPath2, nil
}

// Check arguments.
func CheckArgs(from, to string, offset, limit int64) error {
	if from == "" {
		return errors.New("source file path is empty")
	}
	if to == "" {
		return errors.New("destination file path is empty")
	}

	samePath, err := ComparePaths(from, to)
	if err != nil {
		return err
	}
	if samePath {
		return errors.New("source and destination file paths are same. Must be different")
	}

	if offset < 0 {
		return errors.New("offset cannot be negative")
	}
	if limit < 0 {
		return errors.New("limit cannot be negative")
	}
	return nil
}

// Check is available file and not a directory.
func CheckFile(fromPath string) (os.FileInfo, error) {
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		return nil, ErrUnsupportedFile
	}

	if fileInfo.Mode()&os.ModeType != 0 {
		return nil, ErrUnsupportedFile
	}

	return fileInfo, nil
}

// Check offset.
func CheckOffset(fileInfo os.FileInfo, offset int64) error {
	if offset >= fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}
	return nil
}

// Get file size.
func GetFileSize(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// Copy file.
func Copy(fromPath, toPath string, offset, limit int64) error {
	if err := CheckArgs(fromPath, toPath, offset, limit); err != nil {
		logger.Error("Error validating arguments", "error", err)
		return err
	}
	logger.Info("Arguments validated", "from", fromPath, "to", toPath, "offset", offset, "limit", limit)

	fileInfo, err := CheckFile(fromPath)
	if err != nil {
		logger.Error("Error checking source file", "error", err)
		return err
	}
	logger.Info("File exists and is valid", "path", fromPath)

	// Open source file
	fromFile, err := os.Open(fromPath)
	if err != nil {
		logger.Error("Error opening source file", "error", err)
		return err
	}
	defer fromFile.Close()
	logger.Info("Opened source file", "path", fromPath)

	// Create destination file
	toFile, err := os.Create(toPath)
	if err != nil {
		logger.Error("Error creating destination file", "error", err)
		return err
	}
	defer toFile.Close()
	logger.Info("Created destination file", "path", toPath)

	// Set offset if needed
	if offset > 0 {
		if err := CheckOffset(fileInfo, offset); err != nil {
			logger.Error("Error checking offset", "error", err)
			return err
		}
		if _, err := fromFile.Seek(offset, io.SeekStart); err != nil {
			logger.Error("Error setting file offset", "error", err)
			return err
		}
		logger.Info("Set file offset", "offset", offset)
	}

	// Get file size and unit
	var sizeToCopy int64
	fileSize := fileInfo.Size()
	fileSizeWithOffset := fileSize - offset
	if limit == 0 || limit > fileSizeWithOffset {
		sizeToCopy = fileSizeWithOffset
	} else {
		sizeToCopy = limit
	}

	// Create progress and start bar
	tmpl := `{{ bar . "[" "=" ">" " " "]"}} {{counters .}}`
	bar := pb.ProgressBarTemplate(tmpl).Start64(sizeToCopy)
	bar.SetMaxWidth(BarWidth)
	barReader := bar.NewProxyReader(fromFile)

	// Copy file
	if limit == 0 {
		_, err = io.Copy(toFile, barReader)
	} else {
		_, err = io.CopyN(toFile, barReader, limit)
	}

	if err != nil {
		if errors.Is(err, io.EOF) {
			logger.Info("EOF reached")
		} else {
			logger.Error("Error during file copy", "error", err)
			errRemove := os.Remove(toPath)
			if errRemove != nil {
				logger.Error("Error removing destination file", "error", errRemove)
			}
			return err
		}
	}

	bar.Finish()
	logger.Info("File copied successfully", "from", fromPath, "to", toPath)

	return nil
}
