package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/tidwall/transform"
)

const usage = "Usage: %s COMMAND [ARGS...]"
const readBufferSize = 32 * 1024

func main() {
	err := run()

	if err == nil {
		return
	}

	// If COMMAND returned a non-zero exit code, exit with that same code.
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		os.Exit(exitErr.ExitCode())
	}

	// Log unexpected errors to stderr.
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func run() error {
	// Check args.
	if len(os.Args) < 2 {
		return fmt.Errorf(usage, filepath.Base(os.Args[0]))
	}

	// Create command object.
	cmd := exec.Command(os.Args[1], os.Args[2:]...)

	// Pass our stdin to the command.
	cmd.Stdin = os.Stdin

	// Hook into stdout.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer stdout.Close()

	// Hook into stderr.
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	defer stderr.Close()

	// Run command.
	if err := cmd.Start(); err != nil {
		return err
	}

	// Record start time.
	startedAt := time.Now().Unix()

	wg := sync.WaitGroup{}

	getPrefix := func() func() string {
		return func() string {
			elapsed := time.Now().Unix() - startedAt
			minutes := elapsed / 60
			seconds := elapsed % 60
			return fmt.Sprintf("[%0d:%02d] ", minutes, seconds)
		}
	}()

	// Transform stdout output.
	var stdoutErr error
	go func() {
		_, stdoutErr = io.Copy(os.Stdout, prependDuration(getPrefix, stdout))
		wg.Done()
	}()
	wg.Add(1)

	// Transform stderr output.
	var stderrErr error
	go func() {
		_, stderrErr = io.Copy(os.Stderr, prependDuration(getPrefix, stderr))
		wg.Done()
	}()
	wg.Add(1)

	// Wait for transforms to complete and the command to terminate.
	wg.Wait()
	cmdWaitErr := cmd.Wait()

	// Print a completion timer.
	fmt.Println(getPrefix())

	// Return an error if available.
	if cmdWaitErr != nil {
		return fmt.Errorf("error running command: %w", cmdWaitErr)
	}
	if stdoutErr != nil {
		return fmt.Errorf("error copying stdout: %w", stdoutErr)
	}
	if stderrErr != nil {
		return fmt.Errorf("error copying stderr: %w", stderrErr)
	}

	return nil
}

func prependDuration(getPrefix func() string, r io.Reader) io.Reader {
	// Prepend command duration (e.g `[1:23]`) to the start of each line.

	// Start of stream == start of line.
	startOfLine := true

	readBuffer := make([]byte, readBufferSize)

	return transform.NewTransformer(func() ([]byte, error) {
		n, err := r.Read(readBuffer)
		if err != nil {
			return nil, err
		}

		if n <= 0 {
			return []byte{}, nil
		}

		prefix := getPrefix()

		var writeBuffer bytes.Buffer

		// Write `prefix` if we're at the start of a line - either at the
		// start of the stream, or after the previous call which ended a line.
		if startOfLine {
			writeBuffer.Write([]byte(prefix))
		}

		// Add `prefix` to the start of each line.
		if n > 1 {
			writeBuffer.Write(bytes.Replace(
				readBuffer[:n-1],
				[]byte{'\n'},
				[]byte("\n"+prefix),
				-1,
			))
		}
		// The last byte in a buffer can't start a line, though might be an
		// "\n" character which *ends* a line.
		writeBuffer.WriteByte(readBuffer[n-1])

		// Check if we've ended a line. Set `startOfLine` for the next call to
		// this function.
		startOfLine = readBuffer[n-1] == '\n'

		return writeBuffer.Bytes(), nil
	})
}
