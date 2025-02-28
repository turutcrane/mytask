package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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
	// check existence of the file mytask.go
	if _, err := os.Stat("mytask.go"); err == nil {
		cmdLine := append([]string{"go", "run", "-tags", "mytask", "./mytask.go"}, args[1:]...)
		return mytask.Exec(ctx, "", cmdLine...)
	}

	// check existence of the mytask directory
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("T36: Error: %w", err)
	}
	abs, err := filepath.Abs(pwd)
	if err != nil {
		return fmt.Errorf("T41: Error: %w", err)
	}

	for ; abs != "/"; abs = filepath.Clean(abs + "/..") {
		mytaskPath := filepath.Join(abs, "mytask")
		if d, err := os.Stat(mytaskPath); err == nil {
			if d.IsDir() {
				if *verbose {
					log.Println("Mytask Path", mytaskPath)
				}
				cmdLine := append([]string{"go", "run", ".", "-root", abs, "-current", pwd}, args...)
				return mytask.Exec(ctx, mytaskPath, cmdLine...)
			}
		}
	}
	return nil
}
