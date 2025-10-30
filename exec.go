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

	return ExecEnv(ctx, dir, nil, cmdLine...)
}

// Exec executes a command with the given name, arguments and environment values.
func ExecEnv(ctx context.Context, dir string, env []string, cmdLine ...string) error {
	cmd := exec.CommandContext(ctx, cmdLine[0], cmdLine[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	if env !=nil {
		e := os.Environ()
		cmd.Env = append(e, env...)
	}

	return cmd.Run()
}

// ExecInPipe executes cmdLine with in as Stdin
func ExecInPipe(ctx context.Context, dir string, in io.Reader, cmdLine ...string) (io.Reader, error) {

	return ExecInPipeEnv(ctx, dir, in, nil, cmdLine ...)
}

// ExecInPipeEnv executes cmdLine with in as Stdin and environment values
func ExecInPipeEnv(ctx context.Context, dir string, in io.Reader, env []string,  cmdLine ...string) (io.Reader, error) {
	out := &bytes.Buffer{}
	cmd := exec.CommandContext(ctx, cmdLine[0], cmdLine[1:]...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	if env !=nil {
		e := os.Environ()
		cmd.Env = append(e, env...)
	}
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return out, nil
}
