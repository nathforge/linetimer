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
	if err := run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}

		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf(usage, filepath.Base(os.Args[0]))
	}

	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Stdin = os.Stdin

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	defer stderr.Close()

	if err := cmd.Start(); err != nil {
		return err
	}
	startedAt := time.Now().Unix()

	wg := sync.WaitGroup{}

	var stdoutErr error
	go func() {
		_, stdoutErr = io.Copy(os.Stdout, prependDuration(startedAt, stdout))
		wg.Done()
	}()
	wg.Add(1)

	var stderrErr error
	go func() {
		_, stderrErr = io.Copy(os.Stderr, prependDuration(startedAt, stderr))
		wg.Done()
	}()
	wg.Add(1)

	wg.Wait()
	cmdWaitErr := cmd.Wait()

	if cmdWaitErr != nil {
		return fmt.Errorf("error running command: %s", cmdWaitErr)
	}

	if stdoutErr != nil {
		return fmt.Errorf("error copying stdout: %s", stdoutErr)
	}

	if stderrErr != nil {
		return fmt.Errorf("error copying stderr: %s", stderrErr)
	}

	return nil
}

func prependDuration(startedAt int64, r io.Reader) io.Reader {
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

		elapsed := time.Now().Unix() - startedAt
		minutes := elapsed / 60
		seconds := elapsed % 60
		linePrefix := fmt.Sprintf("[%0d:%02d] ", minutes, seconds)

		var writeBuffer bytes.Buffer

		if startOfLine {
			writeBuffer.Write([]byte(linePrefix))
		}

		writeBuffer.Write(bytes.Replace(
			readBuffer[:n-1],
			[]byte{'\n'},
			[]byte("\n"+linePrefix),
			-1,
		))
		writeBuffer.WriteByte(readBuffer[n-1])

		startOfLine = readBuffer[n-1] == '\n'

		return writeBuffer.Bytes(), nil
	})
}
