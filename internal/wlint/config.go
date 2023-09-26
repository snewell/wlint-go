package wlint

import (
	"io"
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

var (
	homedir string
)

type Config struct {
}

func GetGlobalConfigPath() (string, string) {
	wlintConfigDir := path.Join(homedir, ".wlint")
	wlintConfigFile := path.Join(wlintConfigDir, "wlintrc")
	return wlintConfigFile, wlintConfigDir
}

func FindProjectConfig() (string, string, error) {
	searchDir, err := os.Getwd()
	if err != nil {
		return "", "", err
	}
	for {
		possiblePath := path.Join(searchDir, ".wlintrc")
		if _, err := os.Stat(possiblePath); err == nil {
			return possiblePath, searchDir, nil
		}
		nextDir := path.Dir(searchDir)
		if nextDir == searchDir {
			return "", "", os.ErrNotExist
		}
		searchDir = nextDir
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

func GetAllConfigs[T any]() (ConfigInfo[T], ConfigInfo[T], error) {
	// This function is ugly right now :(
	var globalConfig ConfigInfo[T]
	var localConfig ConfigInfo[T]

	var configFile string
	var err error
	// get global
	configFile, globalConfig.Dir = GetGlobalConfigPath()
	globalConfig.Config, err = loadConfig[T](configFile)
	if err != nil {
		return globalConfig, localConfig, err
	}

	// get local
	configFile, localConfig.Dir, err = FindProjectConfig()
	if err != nil {
		if os.IsNotExist(err) {
			// no local config, we're done
			err = nil
		}
		// some other error; err hasn't been cleared
		return globalConfig, localConfig, err
	}
	localConfig.Config, err = loadConfig[T](configFile)
	return globalConfig, localConfig, err
}

func init() {
	var err error
	homedir, err = os.UserHomeDir()
	if err != nil {
		log.Printf("Error determining user home directory: %v", err)
	}
}
