package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	pb "github.com/cheggaaa/pb/v3"
)

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB

	BarWidth = 100
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileNotFound          = errors.New("file not found")
	ErrUnknownFileSize       = errors.New("unknown file size")
	ErrReadFile              = errors.New("error reading file")
	ErrWriteFile             = errors.New("error writing file")
	ErrSeekFile              = errors.New("error seeking file")
	ErrCopyFile              = errors.New("error copying file")
)

// Check arguments.
func CheckArgs(from, to string, offset, limit int64) error {
	if from == "" {
		return errors.New("source file path is empty")
	}
	if to == "" {
		return errors.New("destination file path is empty")
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
		return nil, ErrFileNotFound
	}

	if fileInfo.IsDir() {
		return nil, ErrUnsupportedFile
	}

	if fileInfo.Size() == 0 {
		return nil, ErrUnknownFileSize
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
		return 0, ErrFileNotFound
	}
	return fileInfo.Size(), nil
}

// Get unit and divisor.
func getUnit(size int64) (unit string, divisor int64) {
	switch {
	case size < KB:
		unit, divisor = "B", 1
	case size < MB:
		unit, divisor = "KB", KB
	case size < GB:
		unit, divisor = "MB", MB
	default:
		unit, divisor = "GB", GB
	}
	return
}

// Copy file.
func Copy(fromPath, toPath string, offset, limit int64) error {
	if err := CheckArgs(fromPath, toPath, offset, limit); err != nil {
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
		return ErrReadFile
	}
	defer fromFile.Close()
	logger.Info("Opened source file", "path", fromPath)

	// Create destination file
	toFile, err := os.Create(toPath)
	if err != nil {
		logger.Error("Error creating destination file", "error", err)
		return ErrWriteFile
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
			return ErrSeekFile
		}
		logger.Info("Set file offset", "offset", offset)
	}

	// Определяем размер копируемых данных
	fileSize := fileInfo.Size()
	unit, divisor := getUnit(fileSize)

	// Создание и запуск прогресс-бара
	tmpl := fmt.Sprintf(`{{ bar . "[" "=" ">" " " "]"}} {{counters .}} (%s)`, unit)
	bar := pb.ProgressBarTemplate(tmpl).Start64(fileSize / divisor)
	bar.SetMaxWidth(BarWidth)

	// Обертка для обновления прогресс-бара
	barReader := bar.NewProxyReader(fromFile)

	// Копирование данных
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
			return ErrCopyFile
		}
	}

	// Завершаем прогресс-бар
	bar.Finish()
	logger.Info("File copied successfully", "from", fromPath, "to", toPath)

	return nil
}
