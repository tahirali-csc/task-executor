package config

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"

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

func Load(config string) (*AppConfig, error) {

	if appConfig == nil {
		data, err := ioutil.ReadFile(config)
		if err != nil {
			return nil, err
		}

		log.Debug("Getting application configuration from::", config)
		appConfig = &AppConfig{}
		err = yaml.Unmarshal(data, appConfig)
		if err != nil {
			return nil, err
		}
	}

	return appConfig, nil
}
