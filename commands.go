package mytask

import (
	"flag"
	"fmt"
	"maps"
	"os"
	"slices"
)

type Config struct {
	RootDir string
	CurDir  string
}

var MytaskConfig Config

func Setup() {
	flag.StringVar(&MytaskConfig.RootDir, "root", ".", "root directory")
	flag.StringVar(&MytaskConfig.CurDir, "current", ".", "current directory")
}

// Command represents a Command that can be executed.
type Command struct {
	key    string
	action func(args []string) ([]string, error)
}

func (cmd Command) Do(args []string) ([]string, error) {
	return cmd.action(args)
}

var (
	cmdList = map[string]Command{}
)

func AddCmd(key string, cmd func([]string) ([]string, error)) {
	cmdList[key] = Command{
		key:    key,
		action: cmd,
	}
}

func GetTask(key string) (Command, bool) {
	if key == "help" || key == "" {
		return Command{
			key: "help",
			action: func(args []string) ([]string, error) {
				verbs := []string{"help"}
				verbs = slices.AppendSeq(verbs, maps.Keys(cmdList))
				fmt.Fprintln(os.Stderr, "Help:", verbs)
				return args, nil
			},
		}, true
	}

	cmd, ok := cmdList[key]
	if ok {
		return cmd, ok
	}
	return Command{}, false
}

func RunTasks(args []string) error {
	if len(args) == 0 {
		args = append(args, "help")
	}

	for len(args) > 0 {
		cmd, ok := GetTask(args[0])
		if ok {
			var err error
			args, err = cmd.Do(args[1:])
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("T110: Error: %s is not a valid command", args[0])
		}
	}

	return nil
}
