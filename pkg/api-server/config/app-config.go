package config

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

type ServerInfo struct {
	Port     string
	KeyFile  string `yaml:"keyFile"`
	CertFile string `yaml:"certFile"`
}

type AppConfig struct {
	Database struct {
		Name     string
		User     string
		Password string
		Host     string
	}

	Server struct {
		Http  ServerInfo
		Https ServerInfo
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
