package mytask

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
)


// Exec executes a command with the given name and arguments.
func Exec(ctx context.Context, dir string, cmdLine ...string) error {
	cmd := exec.CommandContext(ctx, cmdLine[0], cmdLine[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir

	return cmd.Run()
}

// ExecInPipe executes a command with in io.Reader return out io.Reader.
func ExecInPipe(ctx context.Context, dir string, in io.Reader, cmdLine ...string) (io.Reader, error) {
	out := &bytes.Buffer{}
	cmd := exec.CommandContext(ctx, cmdLine[0], cmdLine[1:]...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return out, nil
}
