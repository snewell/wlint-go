package wlint

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

var (
	homedir string
)

type StandardConfig interface {
	GetPurifier() string
}

type Config struct {
	Purifier string `yaml:"purifier"`
}

func (c Config) GetPurifier() string {
	return c.Purifier
}

func GetGlobalConfigPath() (string, string) {
	wlintConfigDir := path.Join(homedir, ".wlint")
	wlintConfigFile := path.Join(wlintConfigDir, "wlintrc")
	return wlintConfigFile, wlintConfigDir
}

func findParentConfigs(startDir string, foundConfigFn func(string, string) error) error {
	for {
		possiblePath := path.Join(startDir, ".wlintrc")
		if _, err := os.Stat(possiblePath); err == nil {
			err := foundConfigFn(possiblePath, startDir)
			if err != nil {
				return err
			}
		}
		nextDir := path.Dir(startDir)
		if nextDir == startDir {
			return nil
		}
		startDir = nextDir
	}
}

func loadConfig[T any](configPath string) (T, error) {
	var ret T

	cf, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// we're done
			return ret, nil
		}
		return ret, err
	}
	defer cf.Close()

	data, err := io.ReadAll(cf)
	if err == nil {
		err = yaml.Unmarshal(data, &ret)
	}
	return ret, err
}

type ConfigInfo[T any] struct {
	Config T
	Dir    string
}

func GetAllConfigs[T any]() ([]ConfigInfo[T], error) {
	ret := []ConfigInfo[T]{}
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error determining current working directory: %v\n", err)
	} else {
		err := findParentConfigs(currentDir, func(configPath string, configDir string) error {
			config, err := loadConfig[T](configPath)
			if err != nil {
				return err
			}
			ret = append(ret, ConfigInfo[T]{
				Config: config,
				Dir:    configDir,
			})
			return nil
		})
		if err != nil {
			return ret, err
		}
	}

	// try reading global config
	globalConfigFile, globalConfigDir := GetGlobalConfigPath()
	globalConfig, err := loadConfig[T](globalConfigFile)
	if err != nil {
		return ret, err
	}
	ret = append(ret, ConfigInfo[T]{
		Config: globalConfig,
		Dir:    globalConfigDir,
	})
	return ret, nil
}

func FindPurifier[T StandardConfig](cliPurifier string, configs []ConfigInfo[T]) (MakePurifierFn, error) {
	if len(cliPurifier) != 0 {
		return GetPurifier(cliPurifier)
	}
	for index := range configs {
		purifier := configs[index].Config.GetPurifier()
		if len(purifier) != 0 {
			return GetPurifier(purifier)
		}
	}
	return nil, nil
}

func init() {
	var err error
	homedir, err = os.UserHomeDir()
	if err != nil {
		log.Printf("Error determining user home directory: %v", err)
	}
}
