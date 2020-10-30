package config

import (
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	Database struct {
		Name     string
		User     string
		Password string
		Host     string
	}

	Server struct {
		Port string
	}
}

var appConfig *AppConfig

func Get() *AppConfig {
	return appConfig
}

func Load() (*AppConfig, error) {

	if appConfig == nil {
		//TODO: Will review!!
		s, _ := os.Getwd()
		data, err := ioutil.ReadFile(path.Join(s, "pkg/api-server", "config.yaml"))
		if err != nil {
			return nil, err
		}

		appConfig = &AppConfig{}
		err = yaml.Unmarshal(data, appConfig)
		if err != nil {
			return nil, err
		}
	}

	return appConfig, nil
}
