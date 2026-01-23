package mytask

import (
	"context"
	"flag"
	"fmt"
	"os"
)

func GetConfig() (Config, error) {
	var tomlPath, curDir string
	flag.StringVar(&tomlPath, "toml", "", "toml filepath")
	flag.StringVar(&curDir, "current", ".", "current directory")
	flag.Parse()
	return ParseConfig(curDir, tomlPath)
}

// Command represents a Command that can be executed.
type Command struct {
	key    string
	action func(ctx context.Context, args []string) ([]string, error)
}

func (cmd Command) Do(ctx context.Context, args []string) ([]string, error) {
	return cmd.action(ctx, args)
}

var (
	cmdList = map[string]Command{}
)

func AddCommand(key string, cmd func(context.Context, []string) ([]string, error)) {
	cmdList[key] = Command{
		key:    key,
		action: cmd,
	}
}

func GetCommand(key string) (Command, bool) {
	cmd, ok := cmdList[key]
	return cmd, ok
}

func HelpVerbList(_ context.Context, args []string) ([]string, error) {
	verbs := []string{}
	for k := range cmdList {
		verbs = append(verbs, k)
	}
	fmt.Fprintln(os.Stderr, "Help:", verbs)
	return args, nil
}

func RunTasks(ctx context.Context, args []string) error {
	if len(args) == 0 {
		args = []string{"help"}
	}

	for len(args) > 0 {
		cmd, ok := GetCommand(args[0])
		if ok {
			var err error
			args, err = cmd.Do(ctx, args[1:])
			if err != nil {
				return err
			}
		} else if args[0] == "help" {
			args, _ = HelpVerbList(ctx, args[1:])
		} else {
			return fmt.Errorf("mytask: Error: %s is not a valid command", args[0])
		}
	}

	return nil
}
