package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/turutcrane/mytask"
)

var verbose = flag.Bool("v", false, "verbose flag")

type mytaskToml struct {
	MytaskDir string `toml:"mytask_dir"`
}

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
	cur, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("T32: Error: %w", err)
	}

	root, err := filepath.Abs(cur)
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
				log.Println("mytask.toml Path", tomlFile)
			}
			// toml ファイルを parse する。
			var mt mytaskToml
			if _, err := toml.DecodeFile(tomlFile, &mt); err != nil {
				return fmt.Errorf("T49: Error: %w", err)
			}
			mytaskDir := mt.MytaskDir
			if !filepath.IsAbs(mytaskDir) {
				mytaskDir = filepath.Join(root, mytaskDir)
			}
			mytaskDir = filepath.Clean(mytaskDir)
			if d, err := os.Stat(mytaskDir); err == nil && d.IsDir() {
				return mytaskDo(ctx, tomlFile, mytaskDir, root, cur, args)
			} else {
				if err != nil {
					return fmt.Errorf("T52: Error: %w", err)
				}
				return fmt.Errorf("T53: Error: %s is not directory", mt.MytaskDir)
			}
		}

		// check existence of the mytask directory
		mytaskPath := filepath.Join(root, "mytask")
		if d, err := os.Stat(mytaskPath); err == nil && d.IsDir() {
			return mytaskDo(ctx, "", mytaskPath, root, cur, args)
		}
		if root == "/" {
			break
		}
		root = filepath.Clean(root + "/..")
	}
	return fmt.Errorf("T67: Error: mytask.go or mytask directory not found")
}

func mytaskDo(ctx context.Context, tomlfile, mytaskDir, root, pwd string, args []string) error {
	if *verbose {
		log.Println("Mytask Path", mytaskDir)
	}
	cmdLine := append([]string{"go", "run", ".", "-toml", tomlfile, "-root", root, "-current", pwd, "-task", mytaskDir}, args...)
	return mytask.Exec(ctx, mytaskDir, cmdLine...)
}
