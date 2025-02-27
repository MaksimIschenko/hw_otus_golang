package envreader

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ProcessBytes removes trailing spaces and replaces null bytes with newlines.
func ProcessBytes(data []byte) []byte {
	data = bytes.TrimRight(data, " \t")
	data = bytes.ReplaceAll(data, []byte{0x00}, []byte("\n"))
	return data
}

// readEnvFile reads the first line from the file and returns it.
func readEnvFile(f *os.File) (string, error) {
	r := bufio.NewReader(f)
	line, _, err := r.ReadLine()
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("error reading file: %w", err)
	}
	line = ProcessBytes(line)
	return string(line), nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envMap := Environment{}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			_envMap, err := ReadDir(fmt.Sprintf("%s/%s", dir, file.Name()))
			if err != nil {
				return nil, err
			}
			for k, v := range _envMap {
				envMap[k] = v
			}
		}
		f, err := os.Open(fmt.Sprintf("%s/%s", dir, file.Name()))
		if err != nil {
			return nil, err
		}
		defer f.Close()
		value, err := readEnvFile(f)
		if err != nil {
			return nil, err
		}
		envMap[file.Name()] = EnvValue{Value: value, NeedRemove: value == ""}
	}
	return envMap, nil
}
