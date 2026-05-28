package mytask

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	rootDir        string
	tomlPath       string
	curDir         string
	taskDir        string
	completionMode string
}

type mytaskToml struct {
	MytaskDir  string `toml:"mytask_dir"`
	RootDir    string `toml:"root_dir"`
	Completion string `toml:"completion"`
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

func (c Config) GetCompletion() string {
	return c.completionMode
}

// ParseConfig tomlPath は絶対パス
func ParseConfig(curDir, tomlFile string) (Config, error) {
	var c Config
	c.curDir = curDir
	c.tomlPath = tomlFile
	tomlDir := filepath.Dir(tomlFile)

	// toml ファイルを parse する。
	var mt mytaskToml
	if _, err := toml.DecodeFile(tomlFile, &mt); err != nil {
		return Config{}, fmt.Errorf("T49: Error: %w", err)
	}

	c.taskDir = dirAbsPath(tomlDir, mt.MytaskDir)
	c.rootDir = dirAbsPath(tomlDir, mt.RootDir)
	c.completionMode = mt.Completion
	return c, nil
}

func dirAbsPath(tomlDir, dir string) string {
	if dir == "" {
		dir = tomlDir
	}
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(tomlDir, dir)
	}
	return filepath.Clean(dir)
}
