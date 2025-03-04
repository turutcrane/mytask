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
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("T32: Error: %w", err)
	}

	abs, err := filepath.Abs(pwd)
	if err != nil {
		return fmt.Errorf("T37: Error: %w", err)
	}

	for {
		// check existence of the file mytask.go
		mytaskGo := filepath.Join(abs, "mytask.go")
		if _, err0 := os.Stat(mytaskGo); err0 == nil {
			if *verbose {
				log.Println("mytask.go Path", mytaskGo)
			}
			cmdLine := append([]string{"go", "run", "-tags", "mytask", "./mytask.go", "-root", abs, "-current", pwd}, args...)
			return mytask.Exec(ctx, abs, cmdLine...)
		}

		// check existence of the mytask directory
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
		if abs == "/" {
			break
		}
		abs = filepath.Clean(abs + "/..")
	}
	return fmt.Errorf("T67: Error: mytask.go or mytask directory not found")
}
