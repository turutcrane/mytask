package mytask

import (
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
	if env != nil {
		e := os.Environ()
		cmd.Env = append(e, env...)
	}

	return cmd.Run()
}

func ExecPipe(ctx context.Context, dir string, in io.Reader, cmdLine ...string) (*exec.Cmd, io.ReadCloser, error) {
	return ExecPipeEnv(ctx, dir, in, nil, cmdLine...)
}

// Exec executes a command with the given name, arguments and environment values. return pipe
func ExecPipeEnv(ctx context.Context, dir string, in io.Reader, env []string, cmdLine ...string) (*exec.Cmd, io.ReadCloser, error) {
	cmd := exec.CommandContext(ctx, cmdLine[0], cmdLine[1:]...)
	cmd.Stdin = in
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	if env != nil {
		e := os.Environ()
		cmd.Env = append(e, env...)
	}
	if p, err := cmd.StdoutPipe(); err == nil {
		if err0 := cmd.Start(); err0 != nil {
			return nil, nil, err0
		}
		return cmd, p, nil
	} else {
		return nil, nil, err
	}
}
