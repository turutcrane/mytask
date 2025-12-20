package mytask

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	rootDir  string
	tomlPath string
	curDir   string
	taskDir  string
}

type mytaskToml struct {
	MytaskDir string `toml:"mytask_dir"`
	RootDir   string `toml:"root_dir"`
}

func (c Config) GetRootDir() string {
	return c.rootDir
}

func (c Config) GetCurDir() string {
	return c.curDir
}

func (c Config) GetTaskDir() string {
	return c.taskDir
}

func (c Config) GetTomlPath() string {
	return c.tomlPath
}

// SetupConfig tomlPath は絶対パス
func SetupConfig(curDir, tomlFile string) (Config, error) {
	var c Config
	c.curDir = curDir
	c.tomlPath = tomlFile

	// toml ファイルを parse する。
	var mt mytaskToml
	if _, err := toml.DecodeFile(tomlFile, &mt); err != nil {
		return Config{}, fmt.Errorf("T49: Error: %w", err)
	}
	root := mt.RootDir
	if root == "" {
		root = filepath.Dir(tomlFile)
	}
	c.rootDir = root

	mytaskDir := mt.MytaskDir
	if mytaskDir == "" {
		if !filepath.IsAbs(mytaskDir) {
			mytaskDir = filepath.Join(c.rootDir, mytaskDir)
		}
		mytaskDir = filepath.Clean(mytaskDir)
	}
	c.taskDir = mytaskDir

	return c, nil
}
