package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/turutcrane/mytask"
)

var verbose = flag.Bool("v", false, "verbose flag")

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	flag.Parse()
	args := flag.Args()

	if err := doMytask(ctx, args); err != nil {
		log.Fatalln(err)
	}
}

func doMytask(ctx context.Context, args []string) error {
	curDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("T32: Error: %w", err)
	}

	root, err := filepath.Abs(curDir)
	if err != nil {
		return fmt.Errorf("T37: Error: %w", err)
	}

	for {
		// check existence of the file mytask.go
		// mytaskGo := filepath.Join(abs, "mytask.go")
		// if _, err0 := os.Stat(mytaskGo); err0 == nil {
		// 	if *verbose {
		// 		log.Println("mytask.go Path", mytaskGo)
		// 	}
		// 	cmdLine := append([]string{"go", "run", "-tags", "mytask", "./mytask.go", "-root", abs, "-current", pwd}, args...)
		// 	return mytask.Exec(ctx, abs, cmdLine...)
		// }

		tomlFile := filepath.Join(root, "mytask.toml")
		if _, err0 := os.Stat(tomlFile); err0 == nil {
			if *verbose {
				slog.Info("mytask:", slog.String("mytask.toml", tomlFile))
			}
			var c mytask.Config
			var err error
			if c, err = mytask.ParseConfig(curDir, tomlFile); err != nil {
				return fmt.Errorf("T48: Error: %w", err)
			}
			if *verbose {
				slog.Info("mytask:", slog.Any("config", c))
			}
			if d, err := os.Stat(c.GetTaskDir()); err == nil && d.IsDir() {
				return mytaskDo(ctx, c, args)
			} else {
				if err != nil {
					return fmt.Errorf("T52: Error: %w", err)
				}
				return fmt.Errorf("T53: Error: %s is not directory", c.GetTaskDir())
			}
		}

		// check existence of the mytask directory
		// mytaskPath := filepath.Join(root, "mytask")
		// if d, err := os.Stat(mytaskPath); err == nil && d.IsDir() {
		// 	return mytaskDo(ctx, "", mytaskPath, root, curDir, args)
		// }
		if root == "/" {
			break
		}
		root = filepath.Clean(root + "/..")
	}
	return fmt.Errorf("T67: Error: mytask.go or mytask directory not found")
}

func mytaskDo(ctx context.Context, c mytask.Config, args []string) error {
	if *verbose {
		slog.Info("mytask:", slog.String("task dir", c.GetTaskDir()))
	}
	cmdLine := append([]string{"go", "run", ".", "-toml", c.GetTomlPath(), "-current", c.GetCurDir()}, args...)
	return mytask.Exec(ctx, c.GetTaskDir(), cmdLine...)
}
