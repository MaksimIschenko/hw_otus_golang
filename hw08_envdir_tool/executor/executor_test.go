package executor

import (
	"bytes"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/MaksimIschenko/hw_otus_golang/hw08_envdir_tool/envreader"
	"github.com/stretchr/testify/require"
)

func TestRunCmd_Success(t *testing.T) {
	env := envreader.Environment{
		"TEST_ENV": {Value: "Hello", NeedRemove: false},
	}

	var stdout, stderr bytes.Buffer

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	rStdout, wStdout, _ := os.Pipe()
	rStderr, wStderr, _ := os.Pipe()
	os.Stdout = wStdout
	os.Stderr = wStderr
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		wStdout.Close()
		wStderr.Close()
	}()

	go func() {
		io.Copy(&stdout, rStdout)
	}()
	go func() {
		io.Copy(&stderr, rStderr)
	}()

	cmd := []string{"echo", "Hello World"}
	exitCode := RunCmd(cmd, env)

	require.Equal(t, OK, exitCode)
	require.Contains(t, stdout.String(), "Hello World")
}

func TestRunCmd_Failure(t *testing.T) {
	env := envreader.Environment{
		"TEST_ENV": {Value: "Hello", NeedRemove: false},
	}

	var stdout, stderr bytes.Buffer

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	rStdout, wStdout, _ := os.Pipe()
	rStderr, wStderr, _ := os.Pipe()
	os.Stdout = wStdout
	os.Stderr = wStderr
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		wStdout.Close()
		wStderr.Close()
	}()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(&stdout, rStdout)
	}()
	go func() {
		defer wg.Done()
		io.Copy(&stderr, rStderr)
	}()

	cmd := []string{"nonexistent_command"}
	exitCode := RunCmd(cmd, env)

	wStdout.Close()
	wStderr.Close()

	wg.Wait()

	require.Equal(t, ERR, exitCode)
	require.Empty(t, stdout.String())
}
