package executor

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/MaksimIschenko/hw_otus_golang/hw08_envdir_tool/envreader"
)

const (
	OK  int = 0
	ERR int = 1
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env envreader.Environment) (returnCode int) {
	for k, v := range env {
		if v.NeedRemove {
			os.Unsetenv(k)
		} else {
			os.Setenv(k, v.Value)
		}
	}

	cleanCmd := filepath.Clean(cmd[0])
	execCmd := exec.Command(cleanCmd, cmd[1:]...)

	// Forwarding input/output
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	// Execute
	err := execCmd.Start()
	if err != nil {
		return ERR
	}

	// Wait for command to finish
	err = execCmd.Wait()
	if err != nil {
		var exitErr *exec.ExitError
		if ok := errors.As(err, &exitErr); ok {
			return exitErr.ExitCode()
		}
		return ERR
	}

	return OK
}
