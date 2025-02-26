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

// func Setup(ctx context.Context) (context.Context, func()) {
// 	signalCtx, clear := signal.NotifyContext(ctx, os.Interrupt)

// 	// go func() {
// 	// 	<-signalCtx.Done()
// 	// 	if childPid.Load() != 0 {
// 	// 		syscall.Kill(int(childPid.Load()), syscall.SIGINT)
// 	// 	}
// 	// }()

// 	return signalCtx, clear
// }

// func pipeSample(ctx context.Context) {
// 	awsProfile := "myprofile"
// 	awsRegion := "us-west-2"
// 	ecrDomain := "123456789012.dkr.ecr.us-west-2.amazonaws.com"
// 	AddCmd("login", func(args []string) ([]string, error) {
// 		pipe, err := ExecInPipe(ctx, "", os.Stdin, "aws", "ecr", "--profile", awsProfile, "get-login-password", "--region", awsRegion)
// 		if err != nil {
// 			return nil, err
// 		}
// 		out, err := ExecInPipe(ctx, "", pipe, "docker", "login", "--username", "AWS", "--password-stdin", ecrDomain)
// 		if err != nil {
// 			return nil, err
// 		}
// 		fmt.Println(out)
// 		return nil
// 	})
// }
