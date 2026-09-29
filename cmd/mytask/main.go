package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/turutcrane/mytask"
)

var verbose = flag.Bool("v", false, "verbose flag")
var completion = flag.Bool("completion", false, "completion flag for complete -C 'mytask -complete' mytask")
var bootstrap = flag.Bool("bootstrap", false, "generate bootstrap code & mytask.toml")

const mytaskToml = "mytask.toml"
const mytaskDir = "mytask"
const tomlContent = `#
mytask_dir = "./mytask"
completion = "bash"
`
const mytaskGo = "mytask.go"
const goContent = `package main

import (
        "context"
        "flag"
        "log"
        "os"
        "os/signal"

        "github.com/turutcrane/mytask"
)

func main() {
        completion := flag.Bool("completion", false, "bash -completion command argToBeCompleted prevArg")
        _, err := mytask.GetConfig()
        if err != nil {
                log.Panicln("T19:", err)
        }

        ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
        defer stop()

        mytask.AddCommand("go-version", func(ctx context.Context, args []string) ([]string, error) {
                return args, mytask.Exec(ctx, "", "go", "version")
        })

        if *completion {
                mytask.Completion(flag.Args())
                return
        }

        if err := mytask.RunTasks(ctx, flag.Args()); err != nil {
                log.Fatalf("mytask Runtask: %v\n", err)
        }
}

`

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	flag.Parse()
	args := flag.Args()
	if *bootstrap {
		// Do nothing if file 'mytask.tom' or directory 'mytask' exists.
		_, err := os.Stat(mytaskToml)
		if !errors.Is(err, fs.ErrNotExist) {
			log.Fatalf("%s is already exist\n", mytaskToml)
		}
		_, err = os.Stat(mytaskDir)
		if !errors.Is(err, fs.ErrNotExist) {
			log.Fatalf("Dir: %s is already exist\n", mytaskDir)
		}
		err = os.WriteFile(mytaskToml, []byte(tomlContent), 0644)
		if err != nil {
			log.Fatalf("create %s: %v", mytaskToml, err)
		}
		err = os.Mkdir(mytaskDir, 0755)
		if err != nil {
			log.Fatalf("mkdir %s: %v", mytaskDir, err)
		}
		err = os.WriteFile(filepath.Join(mytaskDir, mytaskGo), []byte(goContent), 0644)
		if err != nil {
			log.Fatalf("create %s: %v", mytaskGo, err)
		}
		return
	} else {
		if err := doMytask(ctx, args); err != nil {
			log.Fatalln(err)
		}
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

		tomlFile := filepath.Join(root, mytaskToml)
		if _, err0 := os.Stat(tomlFile); err0 == nil {
			if *verbose {
				slog.Info("mytask:", slog.String(mytaskToml, tomlFile))
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
				return mytaskExec(ctx, c, args)
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

func mytaskExec(ctx context.Context, c mytask.Config, args []string) error {
	if *verbose {
		slog.Info("mytask:", slog.String("task dir", c.GetTaskDir()))
	}
	cmdLine := []string{"go", "run", ".", "-toml", c.GetTomlPath(), "-current", c.GetCurDir()}
	if *completion {
		if c.GetCompletion() == "bash" {
			cmdLine = append(cmdLine, "-completion")
		} else {
			return nil
		}
	}
	cmdLine = append(cmdLine, args...)
	return mytask.Exec(ctx, c.GetTaskDir(), cmdLine...)
}
