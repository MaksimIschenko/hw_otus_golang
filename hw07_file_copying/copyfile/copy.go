package copyfile

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileNotFound          = errors.New("file not found")
	ErrUnknowFileSize        = errors.New("unknown file size")
	ErrReadFile              = errors.New("error reading file")
	ErrWriteFile             = errors.New("error writing file")
	ErrSeekFile              = errors.New("error seeking file")
	ErrCopyFile              = errors.New("error copying file")
)

var logger = log.New(os.Stdout, "copyfile: ", log.LstdFlags)

func CheckArgs(from, to string, offset, limit int64) error {
	// Check if from is empty
	if from == "" {
		return errors.New("from is empty")
	}

	// Check if to is empty
	if to == "" {
		return errors.New("to is empty")
	}

	// Check if offset is negative
	if offset < 0 {
		return errors.New("offset is negative")
	}

	// Check if limit is negative
	if limit < 0 {
		return errors.New("limit is negative")
	}

	return nil
}

func CheckFile(fromPath string) error {
	// Check if file exists at fromPath
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return ErrFileNotFound
	}

	// Check if file is not a directory
	if fileInfo.IsDir() {
		return ErrUnsupportedFile
	}

	// Check if file size is unknown
	if fileInfo.Size() == 0 {
		return ErrUnknowFileSize
	}

	return nil
}

func CheckOffset(fromPath string, offset int64) error {
	// Check if file exists at fromPath
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return ErrFileNotFound
	}

	// Check if offset exceeds file size
	if offset >= fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Check if arguments are valid
	err := CheckArgs(fromPath, toPath, offset, limit)
	if err != nil {
		return err
	}
	logger.Printf("Arguments are valid: from=%s, to=%s, offset=%d, limit=%d\n", fromPath, toPath, offset, limit)

	// Check if file exists at fromPath
	err = CheckFile(fromPath)
	if err != nil {
		return err
	}
	logger.Printf("File exists at fromPath: %s\n", fromPath)

	// Open file for reading
	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0o644)
	if err != nil {
		return ErrReadFile
	}
	defer fromFile.Close()
	logger.Printf("File opened for reading: %s\n", fromPath)

	// Create (if not exists) file for writing and open it for writing
	toFile, err := os.OpenFile(toPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return ErrWriteFile
	}
	defer toFile.Close()
	logger.Printf("File created for writing: %s\n", toPath)

	// If offset > 0, check offset and seek to offset
	if offset > 0 {
		err = CheckOffset(fromPath, offset)
		if err != nil {
			return err
		}
		_, err = fromFile.Seek(offset, io.SeekStart)
		if err != nil {
			return ErrSeekFile
		}
		logger.Printf("Seeked to offset: %d\n", offset)
	}

	// Copy file
	if limit == 0 {
		_, err = io.Copy(toFile, fromFile)
	} else {
		_, err = io.CopyN(toFile, fromFile, limit)
	}
	if err != nil {
		if errors.Is(err, io.EOF) {
			fmt.Println("EOF")
			return nil
		}
		fmt.Println(err)
		return ErrCopyFile
	}
	logger.Printf("File copied: %s -> %s\n", fromPath, toPath)

	return nil
}
