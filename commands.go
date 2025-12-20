package mytask

import (
	"context"
	"flag"
	"fmt"
	"os"
)

func Setup() (Config, error) {
	var tomlPath, curDir string
	flag.StringVar(&tomlPath, "toml", "", "toml filepath")
	flag.StringVar(&curDir, "current", ".", "current directory")
	flag.Parse()
	return SetupConfig(curDir, tomlPath)
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
	if key == "help" || key == "" {
		return Command{
			key: "help",
			action: func(_ context.Context, args []string) ([]string, error) {
				verbs := []string{"help"}
				// verbs = slices.AppendSeq(verbs, maps.Keys(cmdList))
				for k := range cmdList {
					verbs = append(verbs, k)
				}
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

func RunTasks(ctx context.Context, args []string) error {
	if len(args) == 0 {
		args = append(args, "help")
	}

	for len(args) > 0 {
		cmd, ok := GetCommand(args[0])
		if ok {
			var err error
			args, err = cmd.Do(ctx, args[1:])
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("T110: Error: %s is not a valid command", args[0])
		}
	}

	return nil
}

